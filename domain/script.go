package domain

type TestCase struct ***REMOVED***
	HardwareOutput []string `json:"hardwareOutput"`
	HardwareInput  []string `json:"hardwareInput"`
***REMOVED***

type TestScript struct ***REMOVED***
	TestCases []TestCase `json:"testCases"`
	Result    bool       `json:"-"`
***REMOVED***
