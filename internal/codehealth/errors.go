package codehealth

import (
	"errors"
	"fmt"
)

// Named validation errors. Callers match them with errors.Is; the wrapped
// message names the offending field.
var (
	ErrUnsupportedVersion = errors.New("unsupported schema version")
	ErrUnknownCapability  = errors.New("unknown capability")
	ErrUnknownProvider    = errors.New("unknown provider")
	ErrProviderMismatch   = errors.New("provider does not serve capability")
	ErrDuplicateInstance  = errors.New("duplicate configured instance")
	ErrInvalidConfig      = errors.New("invalid configuration")
	ErrInvalidThreshold   = errors.New("invalid threshold")
	ErrInvalidOutcome     = errors.New("invalid outcome")
	ErrInvalidFreshness   = errors.New("invalid freshness")
	ErrValueNotAllowed    = errors.New("value not allowed for outcome")
	ErrMissingValue       = errors.New("measured outcome requires a value")
	ErrIncompatibleUnit   = errors.New("incompatible value or unit")
	ErrUnsafeEvidence     = errors.New("unsafe evidence reference")
	ErrUnboundedDetail    = errors.New("detail exceeds bounds")
	ErrMalformedTimestamp = errors.New("malformed timestamp")
	ErrMalformedIdentity  = errors.New("malformed identity")
	ErrIdentityMismatch   = errors.New("identity does not match content")
	ErrInvalidOrigin      = errors.New("invalid origin or retention")
	ErrInvalidRepository  = errors.New("invalid repository identity")
	ErrInvalidSummary     = errors.New("invalid summary")
	ErrBlockingNotAllowed = errors.New("blocking flag not allowed for capability")
	ErrManualSnapshot     = errors.New("manual snapshot cannot be evaluated")
	ErrMalformedRecord    = errors.New("malformed record")
)

// fieldError wraps a sentinel with the field it was raised for.
func fieldError(sentinel error, field, format string, args ...any) error {
	return fmt.Errorf("%w: %s: %s", sentinel, field, fmt.Sprintf(format, args...))
}
