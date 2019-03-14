package controller

import (
	"github.com/containers-ai/federatorai-operator/pkg/controller/alamedaservice"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, alamedaservice.Add)
}
