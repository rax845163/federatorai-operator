package mock

import (
	repository_keycode "github.com/containers-ai/federatorai-operator/pkg/repository/keycode"
)

// KeycodeRepository wraps datahub client
type KeycodeRepository struct {
}

var _ repository_keycode.Interface = &KeycodeRepository{}

// NewKeycodeRepository creates keycode repository that storing keycode to Alameda-Datahub
func NewKeycodeRepository() repository_keycode.Interface {
	return &KeycodeRepository{}
}

// SendKeycode sends keycode to Alameda-Datahub
func (d *KeycodeRepository) SendKeycode(keycode string) error {
	return nil
}

// SendSignatureData sends signature data to Alameda-Datahub
func (d *KeycodeRepository) SendSignatureData(keycode string) error {
	return nil
}

// GetRegistrationData gets registration data from Alameda-Datahub
func (d *KeycodeRepository) GetRegistrationData() (string, error) {
	return "mock-registration-date", nil
}

// GetKeycodeDetail gets Keycode details from Alameda-Datahub
func (d *KeycodeRepository) GetKeycodeDetail(code string) (repository_keycode.Detail, error) {
	return repository_keycode.Detail{}, nil
}

// DeleteKeycode deletes keycode from Alameda-Datahub
func (d *KeycodeRepository) DeleteKeycode(keycode string) error {
	return nil
}

func (d *KeycodeRepository) ListKeycodes() ([]repository_keycode.Detail, error) {

	details := []repository_keycode.Detail{
		repository_keycode.Detail{
			Keycode:        "test-keycode",
			KeycodeType:    "test-keycode-type",
			KeycodeVersion: 0,
			ApplyTime:      nil,
			ExpireTime:     nil,
			LicenseState:   "test-keycode-state",
			Registered:     true,
		},
	}

	return details, nil
}
