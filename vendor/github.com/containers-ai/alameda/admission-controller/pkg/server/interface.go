package server

import (
	"net/http"
)

type AdmissionController interface {
	MutatePod(http.ResponseWriter, *http.Request)
}
