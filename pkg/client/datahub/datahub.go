package datahub

import (
	"context"
	"fmt"

	datahubv1alpha1_client "github.com/containers-ai/api/alameda_api/v1alpha1/datahub"
	datahubv1alpha1_event "github.com/containers-ai/api/alameda_api/v1alpha1/datahub/events"
	"github.com/containers-ai/api/datahub/keycodes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
)

// Client wraps datahub client
type Client struct {
	conn           *grpc.ClientConn
	v1alpha1Client datahubv1alpha1_client.DatahubServiceClient
	keycodesClient keycodes.KeycodesServiceClient
}

// NewDatahubClient creates Datahub instance base on config
func NewDatahubClient(config Config) Client {
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}
	conn, _ := grpc.Dial(config.Address, dialOptions...)
	v1alpha1Client := datahubv1alpha1_client.NewDatahubServiceClient(conn)
	keycodesClient := keycodes.NewKeycodesServiceClient(conn)
	return Client{
		conn:           conn,
		v1alpha1Client: v1alpha1Client,
		keycodesClient: keycodesClient,
	}
}

// Ping checks if Alameda-Datahub can serves request
func (d *Client) Ping() error {
	return nil
}

// Close closes grpc connection
func (d *Client) Close() error {
	return d.conn.Close()
}

// CreateEvents creates events to Alameda-Datahub
func (d *Client) CreateEvents(events []*datahubv1alpha1_event.Event) error {

	ctx := context.TODO()
	req := datahubv1alpha1_event.CreateEventsRequest{
		Events: events,
	}
	respStatus, err := d.v1alpha1Client.CreateEvents(ctx, &req)
	if err != nil {
		return errors.Errorf("create events to Alameda-Datahub failed: %s", err.Error())
	} else if !isResponseStatusOK(respStatus) {
		return errors.Errorf("create events to Alameda-Datahub failed: %s", getResponseStatusDetail(respStatus))
	}

	return nil
}

// ListKeycodes lists keycode from Alameda-Datahub
func (d *Client) ListKeycodes() ([]*keycodes.Keycode, error) {

	ctx := context.TODO()
	req := keycodes.ListKeycodesRequest{}
	resp, err := d.keycodesClient.ListKeycodes(ctx, &req)
	if err != nil {
		return []*keycodes.Keycode{}, errors.Errorf("List keycodes from Alameda-Datahub failed: %s", err.Error())
	} else if !isResponseStatusOK(resp.Status) {
		return []*keycodes.Keycode{}, errors.Errorf("List keycodes from Alameda-Datahub failed: %s", getResponseStatusDetail(resp.Status))
	}

	return resp.Keycodes, nil
}

// AddKeycode adds keycode to Alameda-Datahub
func (d *Client) AddKeycode(keycode string) error {
	ctx := context.TODO()
	req := keycodes.AddKeycodeRequest{
		Keycode: keycode,
	}
	resp, err := d.keycodesClient.AddKeycode(ctx, &req)
	if err != nil {
		return errors.Errorf("Add keycode %s to Alameda-Datahub failed: %s", keycode, err.Error())
	} else if !isResponseStatusOK(resp.Status) {
		return errors.Errorf("Add keycode %s to Alameda-Datahub failed: %s", keycode, getResponseStatusDetail(resp.Status))
	}

	return nil
}

// ActivateRegistrationData activates signature data to Alameda-Datahub
func (d *Client) ActivateRegistrationData(signatureData string) error {
	ctx := context.TODO()
	req := keycodes.ActivateRegistrationDataRequest{
		Data: signatureData,
	}
	respStatus, err := d.keycodesClient.ActivateRegistrationData(ctx, &req)
	if err != nil {
		return errors.Errorf("Activate signature data %s to Alameda-Datahub failed: %s", signatureData, err.Error())
	} else if !isResponseStatusOK(respStatus) {
		return errors.Errorf("Activate signature data %s to Alameda-Datahub failed: %s", signatureData, getResponseStatusDetail(respStatus))
	}

	return nil
}

// GenerateRegistrationData generates registration data from Alameda-Datahub
func (d *Client) GenerateRegistrationData() (string, error) {

	ctx := context.TODO()
	req := empty.Empty{}
	resp, err := d.keycodesClient.GenerateRegistrationData(ctx, &req)
	if err != nil {
		return "", errors.Errorf("Generate registration data from Alameda-Datahub failed: %s", err.Error())
	} else if !isResponseStatusOK(resp.Status) {
		return "", errors.Errorf("Generate registration data from Alameda-Datahub failed: %s", getResponseStatusDetail(resp.Status))
	}

	return resp.Data, nil
}

// DeleteKeycode delete keycode from Alameda-Datahub
func (d *Client) DeleteKeycode(keycode string) error {

	ctx := context.TODO()
	req := keycodes.DeleteKeycodeRequest{
		Keycode: keycode,
	}
	respStatus, err := d.keycodesClient.DeleteKeycode(ctx, &req)
	if err != nil {
		return errors.Errorf("Delete keycode %s from Alameda-Datahub failed: %s", keycode, err.Error())
	} else if !isResponseStatusOK(respStatus) {
		return errors.Errorf("Delete keycode %s from Alameda-Datahub failed: %s", keycode, getResponseStatusDetail(respStatus))
	}

	return nil
}

// GetKeycodeDetail get keycode detail from Alameda-Datahub
func (d *Client) GetKeycodeDetail(keycode string) (keycodes.Keycode, error) {

	ctx := context.TODO()
	req := keycodes.ListKeycodesRequest{
		Keycodes: []string{keycode},
	}
	resp, err := d.keycodesClient.ListKeycodes(ctx, &req)
	if err != nil {
		return keycodes.Keycode{}, errors.Errorf("Get keycode %s detail from Alameda-Datahub failed: %s", keycode, err.Error())
	} else if !isResponseStatusOK(resp.Status) {
		return keycodes.Keycode{}, errors.Errorf("Get keycode %s detail from Alameda-Datahub failed: %s", keycode, getResponseStatusDetail(resp.Status))
	}

	detail := resp.Summary
	if detail == nil {
		return keycodes.Keycode{}, errors.Errorf("Get keycode %s detail from Alameda-Datahub failed: detail is nil", keycode)
	}

	return *detail, nil
}

func isResponseStatusOK(s *status.Status) bool {

	if s != nil && s.Code == int32(code.Code_OK) {
		return true
	}

	return false
}

func getResponseStatusDetail(s *status.Status) string {

	code := "nil"
	msg := "nil"

	if s != nil {
		code = fmt.Sprintf("%d", s.Code)
		msg = s.Message
	}

	return fmt.Sprintf("status code: %s, errMsg: %s", code, msg)
}
