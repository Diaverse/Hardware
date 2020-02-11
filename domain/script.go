package domain

type TestCase struct {
	Responses      []string `json:"Responses"`
	ExpectedOutput []string `json:"ExpectedOutput"`
}

type TestScript struct {
	cases  []TestCase `json:"testCases"`
	Result bool       `json:"-"`
}
