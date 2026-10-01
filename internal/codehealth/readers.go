package codehealth

// DefaultReaders registers the production Reader of every provider in the
// approved catalogue.
func DefaultReaders() Readers {
	junit := JUnitReader{}
	return Readers{
		ProviderGoTestJSON:     GoTestReader{},
		ProviderVitestJUnit:    junit,
		ProviderPytestJUnit:    junit,
		ProviderGoCoverProfile: GoCoverReader{},
		ProviderVitestV8:       VitestCoverageReader{},
		ProviderCoveragePyJSON: CoveragePyReader{},
		ProviderLizardCSV:      LizardReader{},
		ProviderJscpdJSON:      JscpdReader{},
		ProviderOSVScannerJSON: OSVScannerReader{},
	}
}
