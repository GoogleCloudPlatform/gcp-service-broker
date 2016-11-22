package name_generator

import (
	"fmt"
	"time"
)

type SqlInstance interface {
	BasicInstance
	DatabaseName() string
}

type BasicInstance interface {
	InstanceName() string
}

type Generators struct {
	Sql   SqlInstance
	Basic BasicInstance
}

func New() *Generators {
	return &Generators{
		Basic: &BasicNameGenerator{},
		Sql:   &SqlNameGenerator{},
	}
}

type BasicNameGenerator struct {
	count int
}
type SqlNameGenerator struct {
	BasicNameGenerator
}

func (bng *BasicNameGenerator) InstanceName() string {
	bng.count += 1
	return fmt.Sprintf("pcf_sb_%d_%d", bng.count, time.Now().UnixNano())
}

func (sng *SqlNameGenerator) DatabaseName() string {
	return sng.InstanceName()
}
