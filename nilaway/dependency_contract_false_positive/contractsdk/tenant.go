package contractsdk

type Plan struct {
	Name     string
	Features []string
}

var DefaultPlan *Plan

func init() {
	DefaultPlan = &Plan{
		Name:     "enterprise",
		Features: []string{"audit-log", "sso"},
	}
}
