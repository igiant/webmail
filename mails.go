package webmail

import "encoding/json"

type EMail struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	ContactId string `json:"contactId"` // [WRITE-ONCE]
}

type EMailList []EMail

type Attachment struct {
	Id          string `json:"id"`          // origin ID of attachment or ID from upload response
	Url         string `json:"url"`         // [READ-ONLY] Relative URL from root of web. Eg.: /webmail/api/download/attachment/ba5767a9-7a70-4c90-a6bf-dc8dd62e259c/14/0-1-0-1/picture.jpg
	Name        string `json:"name"`        // [WRITE-ONCE] Filename. Can be empty if attachment is inline.
	ContentType string `json:"contentType"` // [WRITE-ONCE]
	ContentId   string `json:"contentId"`   // [WRITE-ONCE] Inline attachment parameter. URI is used in format cid:contentId in HTML parts.
	Size        int    `json:"size"`        // [READ-ONLY]
}

// AttachmentList - array of mail attachment
type AttachmentList []Attachment

type DisplayableContentType string

const (
	ctTextPlain DisplayableContentType = "ctTextPlain"
	ctTextHtml  DisplayableContentType = "ctTextHtml"
)

type DisplayableMimePart struct {
	Id          KId                    `json:"id"` // [READ-ONLY] local identification, 0 = root mime part, 0-0 = sub-mime part, 0.1 = sub-mime part second item...
	ContentType DisplayableContentType `json:"contentType"`
	Content     string                 `json:"content"` // UTF-8 encoded
}

type DisplayableMimePartList []DisplayableMimePart

type MimeHeaderType string

const (
	mhMessageID       MimeHeaderType = "mhMessageID"       // [READ-ONLY]
	mhInReplayTo      MimeHeaderType = "mhInReplayTo"      // The contents of this field identify previous correspondence which this message answers. It contents an original 'Message-ID' value. rfc0822
	mhResentMessageID MimeHeaderType = "mhResentMessageID" // The contents of this field identify a forwarded message. It contains an original 'Message-ID' value. rfc0822
)

type MimeHeader struct {
	Type  MimeHeaderType `json:"type"`
	Value string         `json:"value"`
}

type MimeHeaderList []MimeHeader

type SignInfo struct {
	IsSigned bool               `json:"isSigned"`
	IsValid  bool               `json:"isValid"`
	Error    LocalizableMessage `json:"error"`
	Cert     Certificate        `json:"cert"`
}

type DecryptResult string

const (
	DecryptSuccesful DecryptResult = "DecryptSuccesful"
	DecryptNoKey     DecryptResult = "DecryptNoKey"
	DecryptError     DecryptResult = "DecryptError"
)

type EncryptInfo struct {
	IsEncrypted bool               `json:"isEncrypted"`
	Result      DecryptResult      `json:"result"`
	Error       LocalizableMessage `json:"error"`
}

// Mail - Constants for composing kerio::web::SearchQuery
type Mail struct {
	Id             KId          `json:"id"`       // [READ-ONLY] global identification
	FolderId       KId          `json:"folderId"` // global identification
	Watermark      Watermark    `json:"watermark"`
	From           EMail        `json:"from"`        // contents of the From header
	Sender         EMail        `json:"sender"`      // contents of the Sender header. Use it for delegation. It shoudn't be the same as From header.
	To             EMailList    `json:"to"`          // contents of To header
	Cc             EMailList    `json:"cc"`          // contents of Cc header
	Bcc            EMailList    `json:"bcc"`         // contents of BCc header
	SendDate       UtcDateTime  `json:"sendDate"`    // contents of Date header (to be displayed in Sent Items). Not set for drafts.
	ReceiveDate    UtcDateTime  `json:"receiveDate"` // mail delivery time. Not set for drafts.
	ModifiedDate   UtcDateTime  `json:"modifiedDate"`
	ReplyTo        EMailList    `json:"replyTo"`        // contents of Reply-To header
	NotificationTo EMail        `json:"notificationTo"` // contents of Disposition-Notification-To header
	Subject        string       `json:"subject"`        // contents of Subject header
	Priority       PriorityType `json:"priority"`       // mail priority (contents of X-Priority header). Defaults to normal priority.
	Size           int          `json:"size"`           // mail size in bytes
	// flags ---------------------------------
	IsSeen        bool `json:"isSeen"`
	IsAnswered    bool `json:"isAnswered"`
	IsFlagged     bool `json:"isFlagged"`
	IsForwarded   bool `json:"isForwarded"`
	IsJunk        bool `json:"isJunk"`
	IsMDNSent     bool `json:"isMDNSent"`     // rfc3503; When saving an unfinished message to any folder client MUST set $MDNSent keyword to prevent another client from sending MDN for the message.
	ShowExternal  bool `json:"showExternal"`  // It means the user confirm to show external sources (eg. images). It's false by default.
	RequestDSN    bool `json:"requestDSN"`    // It is true if should be requested DSN even for successful delivered message. DSN for failure is default.
	HasAttachment bool `json:"hasAttachment"` // [read-only]
	IsDraft       bool `json:"isDraft"`       // [read-only]
	IsReadOnly    bool `json:"isReadOnly"`    // [read-only]
	// ---------------------------------------
	// SMIME ---------------------------------
	SignInfo    SignInfo    `json:"signInfo"`    // [read-only]
	EncryptInfo EncryptInfo `json:"encryptInfo"` // [read-only]
	// ---------------------------------------
	DisplayableParts DisplayableMimePartList `json:"displayableParts"`
	Attachments      AttachmentList          `json:"attachments"`
	Headers          MimeHeaderList          `json:"headers"`
	Send             bool                    `json:"send"`    // [WRITE-ONLY] Send mail, beware it invalidate ID of this mail
	Sign             bool                    `json:"sign"`    // [WRITE-ONLY] Appends a signature as attachment
	Encrypt          bool                    `json:"encrypt"` // [WRITE-ONLY] Encypts the email. Certificates of each recipient must be known.
}

type MailList []Mail

// Mail store manager class

// MailsGet - Get a list of e-mails.
//	folderIds - list of global identifiers of folders to be listed.
//	query - query attributes and limits
// Return
//	list - all found e-mails
//  totalItems - number of mails found if there is no limit
func (c *ClientConnection) MailsGet(folderIds KIdList, query SearchQuery) (MailList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Mails.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       MailList `json:"list"`
			TotalItems int      `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// MailsGetPageWithId - Get a list of e-mails.
//	folderIds - list of global identifiers of folders to be listed
//	query - query attributes and limits. Mind that offset is not used
//	id - global identifier of requested email
// Return
//	list - all found e-mails
//  totalItems - number of mails found if there is no limit
func (c *ClientConnection) MailsGetPageWithId(folderIds KIdList, query SearchQuery, id KId) (MailList, int, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
		Id        KId         `json:"id"`
	}{folderIds, query, id}
	data, err := c.CallRaw("Mails.getPageWithId", params)
	if err != nil {
		return nil, 0, 0, err
	}
	list := struct {
		Result struct {
			List       MailList `json:"list"`
			Start      int      `json:"start"`
			TotalItems int      `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.Start, list.Result.TotalItems, err
}

// MailsGetById - Get one particular email. All members of struct Mail are filed in response.
//	ids - global identifiers of requested emails
// Return
//	errors - list of email that failed to obtain
//	result - found emails
func (c *ClientConnection) MailsGetById(ids KIdList) (ErrorList, MailList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Mails.getById", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result MailList  `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// MailsCreate - ErrorCodeSendingFailed - Failed to send email and failed to create mail.
//	mails - new mails.
// Return
//	errors - error message list
//	result - list of ID of crated mails.
func (c *ClientConnection) MailsCreate(mails MailList) (ErrorList, CreateResultList, error) {
	params := struct {
		Mails MailList `json:"mails"`
	}{mails}
	data, err := c.CallRaw("Mails.create", params)
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

// MailsRemove - Remove a list of mails.
//	ids - list of global identifiers of mails to be removed
// Return
//	errors - list of mails that failed to remove
func (c *ClientConnection) MailsRemove(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Mails.remove", params)
	if err != nil {
		return nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, err
}

// MailsSet - ErrorCodeSendingFailed - Failed to send email and failed to update mail.
//	mails - modifications of mails.
// Return
//	errors - error message list
func (c *ClientConnection) MailsSet(mails MailList) (ErrorList, SetResultList, error) {
	params := struct {
		Mails MailList `json:"mails"`
	}{mails}
	data, err := c.CallRaw("Mails.set", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList     `json:"errors"`
			Result SetResultList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// MailsSetAllSeen - Set all e-mail in folder as seen.
//	folderId - target folder
func (c *ClientConnection) MailsSetAllSeen(folderId KId) error {
	params := struct {
		FolderId KId `json:"folderId"`
	}{folderId}
	_, err := c.CallRaw("Mails.setAllSeen", params)
	return err
}

// MailsCopy - Copy existing e-mails to folder
//	ids - list of global identifiers of mails to be copied
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) MailsCopy(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Mails.copy", params)
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

// MailsMove - Move existing e-mails to folder
//	ids - list of global identifiers of e-mails to be moved
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) MailsMove(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Mails.move", params)
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

// MailsExportAttachments - Export attachments from mail and pack them into zip.
//	attachmentIds - list of global identifiers of attachments. All attachments must be from the same e-mail.
// Return
//	fileDownload - description of output file
func (c *ClientConnection) MailsExportAttachments(attachmentIds KIdList) (*Download, error) {
	params := struct {
		AttachmentIds KIdList `json:"attachmentIds"`
	}{attachmentIds}
	data, err := c.CallRaw("Mails.exportAttachments", params)
	if err != nil {
		return nil, err
	}
	fileDownload := struct {
		Result struct {
			FileDownload Download `json:"fileDownload"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &fileDownload)
	return &fileDownload.Result.FileDownload, err
}
