package webmail

import "encoding/json"

type OperatorExtension struct {
	ExtensionId  KId    `json:"extensionId"`
	TelNum       string `json:"telNum"`
	Description  string `json:"description"`
	IsRegistered bool   `json:"isRegistered"`
}

type OperatorExtensionList []OperatorExtension

type OperatorCallStatus string

const (
	OcsUnknown   OperatorCallStatus = "OcsUnknown" // Disconnected, invalid callId, or other error
	OcsPickUp    OperatorCallStatus = "OcsPickUp"
	OcsRinging   OperatorCallStatus = "OcsRinging"
	OcsConnected OperatorCallStatus = "OcsConnected"
)

// Manager to handle operator requests

// CallManagerGetExtensions - empty extensions = no extension available
// Return
//	extensions
func (c *ClientConnection) CallManagerGetExtensions() (OperatorExtensionList, error) {
	data, err := c.CallRaw("CallManager.getExtensions", nil)
	if err != nil {
		return nil, err
	}
	extensions := struct {
		Result struct {
			Extensions OperatorExtensionList `json:"extensions"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &extensions)
	return extensions.Result.Extensions, err
}

// CallManagerDial - Dials requested phone number
// Parameters
//	extensionId - is phone number from which will be call initiated
//	phoneNumber - is phone number to call
// Return
//	callId - returns id of phone call
func (c *ClientConnection) CallManagerDial(extensionId KId, phoneNumber string) (*KId, error) {
	params := struct {
		ExtensionId KId    `json:"extensionId"`
		PhoneNumber string `json:"phoneNumber"`
	}{extensionId, phoneNumber}
	data, err := c.CallRaw("CallManager.dial", params)
	if err != nil {
		return nil, err
	}
	callId := struct {
		Result struct {
			CallId KId `json:"callId"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &callId)
	return &callId.Result.CallId, err
}

// CallManagerLogin - Dials requested phone number
// Parameters
//	userName - user login name
//	password - user password
func (c *ClientConnection) CallManagerLogin(userName string, password string) error {
	params := struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}{userName, password}
	_, err := c.CallRaw("CallManager.login", params)
	return err
}

// CallManagerHangup - Dials requested phone number
func (c *ClientConnection) CallManagerHangup(callId KId) error {
	params := struct {
		CallId KId `json:"callId"`
	}{callId}
	_, err := c.CallRaw("CallManager.hangup", params)
	return err
}

// CallManagerGetCallStatus - Dials requested phone number
func (c *ClientConnection) CallManagerGetCallStatus(lastStatus OperatorCallStatus, callId KId) (*OperatorCallStatus, error) {
	params := struct {
		LastStatus OperatorCallStatus `json:"lastStatus"`
		CallId     KId                `json:"callId"`
	}{lastStatus, callId}
	data, err := c.CallRaw("CallManager.getCallStatus", params)
	if err != nil {
		return nil, err
	}
	status := struct {
		Result struct {
			Status OperatorCallStatus `json:"status"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &status)
	return &status.Result.Status, err
}
