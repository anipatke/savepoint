package codehealth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"slices"
	"time"
)

// Config is the project-owned Code Health configuration (schema version 1).
type Config struct {
	Version      int                `json:"version"`
	Capabilities []CapabilityConfig `json:"capabilities"`
}

// CapabilityConfig configures one provider instance. The instance key is the
// capability, provider, and optional name together, so a monorepo can run one
// provider over several scopes. Required marks an instance whose failure the
// project will not tolerate; it defaults to optional.
type CapabilityConfig struct {
	Capability     Capability  `json:"capability"`
	Provider       ProviderKey `json:"provider"`
	Name           string      `json:"name,omitempty"`
	Required       bool        `json:"required,omitempty"`
	Executable     string      `json:"executable,omitempty"`
	Args           []string    `json:"args,omitempty"`
	Report         string      `json:"report,omitempty"`
	TimeoutSeconds int         `json:"timeout_seconds,omitempty"`
	Scope          []string    `json:"scope,omitempty"`
	Exclusions     []string    `json:"exclusions,omitempty"`
	Thresholds     *Threshold  `json:"thresholds,omitempty"`
}

// Threshold holds current-value guidance. Direction follows the capability:
// coverage must stay at or above, every other measure at or below. A value
// within Good is good, within Watch is watch, and beyond Watch needs attention.
type Threshold struct {
	Good  float64 `json:"good"`
	Watch float64 `json:"watch"`
}

// Validate checks the whole configuration and returns the first problem.
func (c Config) Validate() error {
	if c.Version != ConfigVersion {
		return fieldError(ErrUnsupportedVersion, "version", "got %d, want %d", c.Version, ConfigVersion)
	}
	if len(c.Capabilities) > MaxResults {
		return fieldError(ErrUnboundedDetail, "capabilities", "has %d entries; at most %d", len(c.Capabilities), MaxResults)
	}
	seen := make(map[instanceKey]bool, len(c.Capabilities))
	for i, cc := range c.Capabilities {
		field := indexed("capabilities", i)
		if err := cc.validate(field); err != nil {
			return err
		}
		key := instanceKey{cc.Capability, cc.Provider, cc.Name}
		if seen[key] {
			return fieldError(ErrDuplicateInstance, field, "%s via %s named %q is configured more than once", cc.Capability, cc.Provider, cc.Name)
		}
		seen[key] = true
	}
	return c.validateSharedProviders()
}

// instanceKey identifies one configured instance, or one result of it.
type instanceKey struct {
	capability Capability
	provider   ProviderKey
	name       string
}

// validateSharedProviders requires every instance that shares a capability and
// provider to be named and to cover a different scope, so two copies of a tool
// never measure the same files twice or get confused for one another.
func (c Config) validateSharedProviders() error {
	type pair struct {
		capability Capability
		provider   ProviderKey
	}
	groups := make(map[pair][]int)
	for i, cc := range c.Capabilities {
		p := pair{cc.Capability, cc.Provider}
		groups[p] = append(groups[p], i)
	}
	for _, idx := range groups {
		if len(idx) < 2 {
			continue
		}
		for n, i := range idx {
			cc := c.Capabilities[i]
			field := indexed("capabilities", i)
			if cc.Name == "" {
				return fieldError(ErrInvalidConfig, field+".name", "%s via %s is configured %d times; each instance needs a name", cc.Capability, cc.Provider, len(idx))
			}
			for _, j := range idx[n+1:] {
				if sameScope(cc.Scope, c.Capabilities[j].Scope) {
					return fieldError(ErrInvalidConfig, indexed("capabilities", j)+".scope", "instances %q and %q of %s via %s cover the same scope", cc.Name, c.Capabilities[j].Name, cc.Capability, cc.Provider)
				}
			}
		}
	}
	return nil
}

// sameScope compares scopes as sets, since order carries no meaning.
func sameScope(a, b []string) bool {
	return slices.Equal(sortedCopy(a), sortedCopy(b))
}

// EffectiveTimeout returns the project override when set, else the provider's
// default. The second result is false when neither exists, which is the case
// for report-only providers that are never executed.
func (cc CapabilityConfig) EffectiveTimeout() (time.Duration, bool) {
	if cc.TimeoutSeconds > 0 {
		return time.Duration(cc.TimeoutSeconds) * time.Second, true
	}
	if s, ok := DefaultTimeoutSeconds(cc.Provider); ok {
		return time.Duration(s) * time.Second, true
	}
	return 0, false
}

func (cc CapabilityConfig) validate(field string) error {
	if err := validateCapabilityProvider(field, cc.Capability, cc.Provider); err != nil {
		return err
	}
	if err := validateToken(field+".name", cc.Name, false); err != nil {
		return err
	}
	if err := cc.validateExecution(field); err != nil {
		return err
	}
	if err := validatePatternList(field+".scope", cc.Scope); err != nil {
		return err
	}
	if err := validatePatternList(field+".exclusions", cc.Exclusions); err != nil {
		return err
	}
	if cc.Thresholds != nil {
		return cc.Thresholds.validate(field+".thresholds", cc.Capability)
	}
	return nil
}

func validateCapabilityProvider(field string, c Capability, p ProviderKey) error {
	if _, ok := capabilityUnits[c]; !ok {
		return fieldError(ErrUnknownCapability, field+".capability", "%q", c)
	}
	serves, ok := providerCapability[p]
	if !ok {
		return fieldError(ErrUnknownProvider, field+".provider", "%q", p)
	}
	if serves != c {
		return fieldError(ErrProviderMismatch, field+".provider", "%s serves %s, not %s", p, serves, c)
	}
	return nil
}

func (cc CapabilityConfig) validateExecution(field string) error {
	if err := validateText(field+".executable", cc.Executable, MaxPathLen); err != nil {
		return err
	}
	if len(cc.Args) > MaxArgs {
		return fieldError(ErrUnboundedDetail, field+".args", "has %d entries; at most %d", len(cc.Args), MaxArgs)
	}
	for i, a := range cc.Args {
		if err := validateText(indexed(field+".args", i), a, MaxArgLen); err != nil {
			return err
		}
	}
	if cc.Report != "" {
		if err := validateRelativePath(ErrInvalidConfig, field+".report", cc.Report); err != nil {
			return err
		}
	}
	if cc.TimeoutSeconds < 0 || cc.TimeoutSeconds > MaxTimeoutSecond {
		return fieldError(ErrInvalidConfig, field+".timeout_seconds", "%d is outside 0..%d", cc.TimeoutSeconds, MaxTimeoutSecond)
	}
	return nil
}

func (t Threshold) validate(field string, c Capability) error {
	bounds := []struct {
		name string
		v    float64
	}{{"good", t.Good}, {"watch", t.Watch}}
	for _, b := range bounds {
		if math.IsNaN(b.v) || math.IsInf(b.v, 0) || b.v < 0 {
			return fieldError(ErrInvalidThreshold, field+"."+b.name, "%v must be a finite non-negative number", b.v)
		}
		if capabilityUnits[c] == UnitPercent && b.v > 100 {
			return fieldError(ErrInvalidThreshold, field+"."+b.name, "%v exceeds 100 percent", b.v)
		}
	}
	if higherIsBetter(c) && t.Good < t.Watch {
		return fieldError(ErrInvalidThreshold, field, "good %v must not be below watch %v", t.Good, t.Watch)
	}
	if !higherIsBetter(c) && t.Good > t.Watch {
		return fieldError(ErrInvalidThreshold, field, "good %v must not exceed watch %v", t.Good, t.Watch)
	}
	return nil
}

// DecodeConfig parses and validates a JSON configuration. Unknown fields are
// rejected so a typo cannot silently become a default.
func DecodeConfig(data []byte) (Config, error) {
	var c Config
	if err := decodeStrict(data, &c); err != nil {
		return Config{}, err
	}
	return c, c.Validate()
}

// decodeRecord decodes one nested object, rejecting unknown fields like the
// root decoder does. Custom unmarshalers use it because encoding/json does not
// carry DisallowUnknownFields into them.
func decodeRecord(data []byte, into any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(into); err != nil {
		return fieldError(ErrMalformedRecord, "json", "%v", err)
	}
	return nil
}

// requiredNumber reads a number that must be present and non-null, so an
// omitted measurement is never read back as a measured zero.
func requiredNumber(field string, p *float64) (float64, error) {
	if p == nil {
		return 0, fieldError(ErrMissingValue, field, "a number is required; omitted or null is not zero")
	}
	return *p, nil
}

// UnmarshalJSON requires good and watch to be present numbers.
func (t *Threshold) UnmarshalJSON(data []byte) error {
	var raw struct {
		Good  *float64 `json:"good"`
		Watch *float64 `json:"watch"`
	}
	if err := decodeRecord(data, &raw); err != nil {
		return err
	}
	good, err := requiredNumber("thresholds.good", raw.Good)
	if err != nil {
		return err
	}
	watch, err := requiredNumber("thresholds.watch", raw.Watch)
	if err != nil {
		return err
	}
	*t = Threshold{Good: good, Watch: watch}
	return nil
}

func decodeStrict(data []byte, into any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(into); err != nil {
		if errors.Is(err, ErrMissingValue) {
			return err
		}
		return fieldError(ErrMalformedRecord, "json", "%v", err)
	}
	// More reports false on a closing delimiter, so demand end of input instead.
	if _, err := dec.Token(); !errors.Is(err, io.EOF) {
		return fieldError(ErrMalformedRecord, "json", "trailing data after record")
	}
	return nil
}
