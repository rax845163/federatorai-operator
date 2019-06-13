package datahub

import (
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
)

// IsResponseStatusOK check if response is ok
func IsResponseStatusOK(s *status.Status) (bool, error) {

	if s == nil {
		return false, errors.New("receive nil status from datahub")
	} else if s.Code != int32(code.Code_OK) {
		return false, errors.Errorf("status code not 0: receive status code: %d,message: %s", s.Code, s.Message)
	}

	return true, nil
}
