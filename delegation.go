package webmail

type InboundDelegation struct {
	Principal kerio::jsonapi::webmail::principals::Principal `json:"principal"` // [READ-ONLY]
	MailboxId KId `json:"mailboxId"` // [READ-ONLY] root folder ID 
	Accepted bool `json:"accepted"` 
}

type InboundDelegationList []InboundDelegation

type OutboundDelagation struct {
	Principal kerio::jsonapi::webmail::principals::Principal `json:"principal"` 
	IsInboxRW bool `json:"isInboxRW"` 
}

type OutboundDelagationList []OutboundDelagation


// DelegationGet - Get list of accounts which the user set for delegation.
// Return
//	list - delegates
func (c *ClientConnection) DelegationGet() (OutboundDelagationList, error) {
	data, err := c.CallRaw("Delegation.get", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List OutboundDelagationList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// DelegationSet - Set list of accounts for delegation.
// Parameters
//	list - delegates; Only type 'User' is valid.
func (c *ClientConnection) DelegationSet(list OutboundDelagationList) error {
	params := struct {
		List OutboundDelagationList `json:"list"`
	}{list}
	_, err := c.CallRaw("Delegation.set", params)
	return err
}

// DelegationGetInbound - Get list of accounts whom is the user delegate.
// Return
//	list - delegates
func (c *ClientConnection) DelegationGetInbound() (InboundDelegationList, error) {
	data, err := c.CallRaw("Delegation.getInbound", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List InboundDelegationList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// DelegationSetInbound - Set list of accounts whom is the user delegate.
// Parameters
//	list - delegates
func (c *ClientConnection) DelegationSetInbound(list InboundDelegationList) error {
	params := struct {
		List InboundDelegationList `json:"list"`
	}{list}
	_, err := c.CallRaw("Delegation.setInbound", params)
	return err
}