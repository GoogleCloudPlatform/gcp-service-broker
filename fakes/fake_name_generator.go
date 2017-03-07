package fakes

type StaticNameGenerator struct {
	Val string
}

func (sg *StaticNameGenerator) InstanceName() string {
	return sg.Val
}

func (sg *StaticNameGenerator) InstanceNameWithSeparator(sep string) string {
	return sg.Val
}

func (sg *StaticNameGenerator) DatabaseName() string {
	return sg.Val
}

type StaticSQLNameGenerator struct {
	StaticNameGenerator
}

func (sng *StaticSQLNameGenerator) InstanceName() string {
	return sng.Val
}

func (sng *StaticSQLNameGenerator) DatabaseName() string {
	return sng.Val
}

func (sng *StaticSQLNameGenerator) GenerateUsername(instanceID, bindingID string) (string, error) {
	return sng.Val[:16], nil
}

func (sng *StaticSQLNameGenerator) GeneratePassword() (string, error) {
	return sng.Val, nil
}
