package broker

// ServiceExample holds example configurations for a service that _should_
// work.
type ServiceExample struct {
	// Name is a human-readable name of the example
	Name string
	// Descrpition is a long-form description of what this example is about
	Description string
	// PlanId is the plan this example will run against.
	PlanId string

	// ProvisionParams is the JSON object that will be passed to provision
	ProvisionParams map[string]interface{}

	// BindParams is the JSON object that will be passed to bind. If nil,
	// this example DOES NOT include a bind portion.
	BindParams map[string]interface{}
}
