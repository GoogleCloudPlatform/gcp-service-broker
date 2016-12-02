package fakes

type StaticNameGenerator struct {
	Val string
}

func (sg *StaticNameGenerator) InstanceName() string {
	return sg.Val
}

func (sg *StaticNameGenerator) DatabaseName() string {
	return sg.Val
}
