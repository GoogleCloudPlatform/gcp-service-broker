package broker

const (
	JsonTypeString  = "string"
	JsonTypeNumeric = "number"
	JsonTypeInteger = "integer"
	JsonTypeBoolean = "boolean"
)

type JsonType string

type BrokerVariable struct {
	Required  bool
	FieldName string
	Type      JsonType
	Details   string
	//Validation []Validator
	Default interface{}
}
