package controller

// Validator is an interface defining controller validation functions
type Validator interface {
	IsControllerEnabledExecution(namespace, name, kind string) (bool, error)
}
