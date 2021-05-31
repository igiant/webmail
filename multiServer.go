package webmail

import "encoding/json"

type RemoteItem struct {
	Url         string      `json:"url"`
	Flags       uint64      `json:"flags"`
	ReceiveDate UtcDateTime `json:"receiveDate"`
	IsMove      bool        `json:"isMove"`
}

type RemoteItemList []RemoteItem

type EmailCertificate struct {
	Email       EMail  `json:"email"`
	Certificate string `json:"certificate"`
}

type EmailCertificateList []EmailCertificate

// MultiServerAppendRemoteItem -
func (c *ClientConnection) MultiServerAppendRemoteItem(items RemoteItem, folderId KId) (*CreateResult, error) {
	params := struct {
		Items    RemoteItem `json:"items"`
		FolderId KId        `json:"folderId"`
	}{items, folderId}
	data, err := c.CallRaw("MultiServer.appendRemoteItem", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Result struct {
			Result CreateResult `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	return &result.Result.Result, err
}

// MultiServerAppendRemoteItems -
func (c *ClientConnection) MultiServerAppendRemoteItems(items RemoteItemList, folderId KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Items    RemoteItemList `json:"items"`
		FolderId KId            `json:"folderId"`
	}{items, folderId}
	data, err := c.CallRaw("MultiServer.appendRemoteItems", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList        `json:"errors"`
			Result CreateResultList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// MultiServerGetCertificates -
func (c *ClientConnection) MultiServerGetCertificates(emails EMailList) (ErrorList, EmailCertificateList, error) {
	params := struct {
		Emails EMailList `json:"emails"`
	}{emails}
	data, err := c.CallRaw("MultiServer.getCertificates", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList            `json:"errors"`
			Result EmailCertificateList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}
