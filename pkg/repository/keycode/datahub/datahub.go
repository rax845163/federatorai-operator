package datahub

import (
	"time"

	"github.com/containers-ai/federatorai-operator/pkg/client/datahub"
	repository_keycode "github.com/containers-ai/federatorai-operator/pkg/repository/keycode"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
)

// KeycodeRepository wraps datahub client
type KeycodeRepository struct {
	client *datahub.Client
}

var _ repository_keycode.Interface = &KeycodeRepository{}

// NewKeycodeRepository creates keycode repository that storing keycode to Alameda-Datahub
func NewKeycodeRepository(client *datahub.Client) repository_keycode.Interface {
	return &KeycodeRepository{
		client: client,
	}
}

// SendKeycode sends keycode to Alameda-Datahub
func (d *KeycodeRepository) SendKeycode(keycode string) error {
	if err := d.client.AddKeycode(keycode); err != nil {
		return err
	}
	return nil
}

// SendSignatureData sends signature data to Alameda-Datahub
func (d *KeycodeRepository) SendSignatureData(keycode string) error {
	if err := d.client.ActivateRegistrationData(keycode); err != nil {
		return err
	}
	return nil
}

// GetRegistrationData gets registration data from Alameda-Datahub
func (d *KeycodeRepository) GetRegistrationData() (string, error) {
	resp, err := d.client.GenerateRegistrationData()
	if err != nil {
		return "", err
	}
	return resp, nil
}

// GetKeycodeDetail gets Keycode details from Alameda-Datahub
func (d *KeycodeRepository) GetKeycodeDetail(code string) (repository_keycode.Detail, error) {

	keycode, err := d.client.GetKeycodeDetail(code)
	if err != nil {
		return repository_keycode.Detail{}, err
	}

	var applyTime *time.Time
	if keycode.ApplyTime != nil {
		t, err := ptypes.Timestamp(keycode.ApplyTime)
		if err != nil {
			return repository_keycode.Detail{}, errors.Errorf("convert timestamp failed: %s", err.Error())
		}
		applyTime = &t

	}
	var expireTime *time.Time
	if keycode.ExpireTime != nil {
		t, err := ptypes.Timestamp(keycode.ExpireTime)
		if err != nil {
			return repository_keycode.Detail{}, errors.Errorf("convert timestamp failed: %s", err.Error())
		}
		expireTime = &t
	}
	summary := repository_keycode.Detail{
		Keycode:        keycode.Keycode,
		KeycodeType:    keycode.KeycodeType,
		KeycodeVersion: keycode.KeycodeVersion,
		ApplyTime:      applyTime,
		ExpireTime:     expireTime,
		LicenseState:   keycode.LicenseState,
		Registered:     keycode.Registered,
	}

	return summary, nil
}

// ListKeycodes lists keycode from Alameda-Datahub
func (d *KeycodeRepository) ListKeycodes() ([]repository_keycode.Detail, error) {

	keycodes, err := d.client.ListKeycodes()
	if err != nil {
		return []repository_keycode.Detail{}, err
	}

	details := make([]repository_keycode.Detail, len(keycodes))
	for i, keycode := range keycodes {
		var applyTime *time.Time
		if keycode.ApplyTime != nil {
			t, err := ptypes.Timestamp(keycode.ApplyTime)
			if err != nil {
				return []repository_keycode.Detail{}, errors.Errorf("convert timestamp failed: %s", err.Error())
			}
			applyTime = &t
		}
		var expireTime *time.Time
		if keycode.ExpireTime != nil {
			t, err := ptypes.Timestamp(keycode.ExpireTime)
			if err != nil {
				return []repository_keycode.Detail{}, errors.Errorf("convert timestamp failed: %s", err.Error())
			}
			expireTime = &t
		}
		details[i] = repository_keycode.Detail{
			Keycode:        keycode.Keycode,
			KeycodeType:    keycode.KeycodeType,
			KeycodeVersion: keycode.KeycodeVersion,
			ApplyTime:      applyTime,
			ExpireTime:     expireTime,
			LicenseState:   keycode.LicenseState,
			Registered:     keycode.Registered,
		}
	}

	return details, nil
}

// DeleteKeycode deletes keycode from Alameda-Datahub
func (d *KeycodeRepository) DeleteKeycode(keycode string) error {
	if err := d.client.DeleteKeycode(keycode); err != nil {
		return err
	}
	return nil
}
