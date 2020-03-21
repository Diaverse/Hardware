package domain

type TestCase struct {
	HardwareOutput []string `json:"hardwareOutput"`
	HardwareInput  []string `json:"hardwareInput"`
	Result         float64  `json:"-"`
	TotalPassed    int      `json:"totalPass, omitempty"`
	TotalFailed    int      `json:"totalFail, omitempty"`
}

type TestScript struct {
	TestCases   []TestCase `json:"testCases"`
	PassPercent float64    `json:"passPercent, omitempty"`
}
