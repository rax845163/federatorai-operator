package keycode

import (
	"fmt"
	"time"
)

// Detail includes keycode informations
type Detail struct {
	Keycode        string
	KeycodeType    string
	KeycodeVersion int32
	ApplyTime      *time.Time
	ExpireTime     *time.Time
	LicenseState   string
	Registered     bool
}

// Summary returns keycode summary
func (d Detail) Summary() string {
	return fmt.Sprintf("%+v", d)
}

// Interface defines keycode repository interface
type Interface interface {
	SendKeycode(string) error
	SendSignatureData(string) error
	GetRegistrationData() (string, error)
	GetKeycodeDetail(string) (Detail, error)
	DeleteKeycode(string) error
}
