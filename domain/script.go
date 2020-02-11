package domain

type TestCase struct ***REMOVED***
	Responses      []string `json:"Responses"`
	ExpectedOutput []string `json:"ExpectedOutput"`
***REMOVED***

type TestScript struct ***REMOVED***
	cases  []TestCase `json:"testCases"`
	Result bool       `json:"-"`
***REMOVED***
