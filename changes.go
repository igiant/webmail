package webmail

type ChangeType string

const (
	chtCreated          ChangeType = "chtCreated"
	chtCopied           ChangeType = "chtCopied"
	chtDeleted          ChangeType = "chtDeleted"
	chtModified         ChangeType = "chtModified"
	chtMoved            ChangeType = "chtMoved"
	chtNewMail          ChangeType = "chtNewMail" // Valid only for item type 'itMail' and not for folder. Folder filter is not applied for this.
	chtStatus           ChangeType = "chtStatus"
	chtReadFlagChanged  ChangeType = "chtReadFlagChanged"  // Valid only for item type 'itMail' and not for folder.
	chtModifiedMetadata ChangeType = "chtModifiedMetadata" // Valid only for item type 'itMail' and not for folder. Content of messge is not changed, basicaly only flags.
	chtModifiedContent  ChangeType = "chtModifiedContent"  // Valid only for folder. It means there was a change in messages which are placed in the folder. E.g.: number of unread messages.
)

type ItemType string

const (
	itMail          ItemType = "itMail" // change per mailbox for type 'chtNewMail' (folder filter is not applied)
	itCalendar      ItemType = "itCalendar"
	itContact       ItemType = "itContact"
	itTask          ItemType = "itTask"
	itNote          ItemType = "itNote"
	itCalendarInbox ItemType = "itCalendarInbox" // change per mailbox (folder filter is not applied)
	itDelegation    ItemType = "itDelegation"    // Valid ChangeType is 'chtCreated' and 'chtDeleted'
)

type AccountSyncKey struct {
	Guid      string    `json:"guid"`
	Watermark Watermark `json:"watermark"`
}

type AccountSyncKeyList []AccountSyncKey

type SyncKey struct {
	Id             int                `json:"id"`
	Version        int                `json:"version"`
	Watermark      Watermark          `json:"watermark"`
	PublicFolder   Watermark          `json:"publicFolder"`
	AccountSyncKey AccountSyncKeyList `json:"accountSyncKey"`
}

type Change struct {
	IsFolder      bool       `json:"isFolder"`
	Type          ChangeType `json:"type"`
	ItemType      ItemType   `json:"itemType"`
	ItemId        KId        `json:"itemId"`
	ParentId      KId        `json:"parentId"`
	OrigId        KId        `json:"origId"`
	OrigParentId  KId        `json:"origParentId"`
	Watermark     Watermark  `json:"watermark"`
	MessageUnread int        `json:"messageUnread"` // when type is chtModifiedContent, it contains number of unreaded messages; when type is chtReadFlagChanged, it contains 1 if a message is unreaded
}

type ChangeList []Change

// Changes manager class

// ChangesGet - Is permitted only one long-poll request with the same 'lastSyncKey'.
// Parameters
//	lastSyncKey - last watermark
//	timeout - max time to wait for new changes. If value is zero response is returned imediately.
// Return
//	list - all found changes
//	syncKey - new watermark
func (c *ClientConnection) ChangesGet(lastSyncKey SyncKey, timeout int) (ChangeList, *SyncKey, error) {
	params := struct {
		LastSyncKey SyncKey `json:"lastSyncKey"`
		Timeout     int     `json:"timeout"`
	}{lastSyncKey, timeout}
	data, err := c.CallRaw("Changes.get", params)
	if err != nil {
		return nil, nil, err
	}
	list := struct {
		Result struct {
			List    ChangeList `json:"list"`
			SyncKey SyncKey    `json:"syncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, &list.Result.SyncKey, err
}

// ChangesGetAccount - Get changes for all accessible folders of particular user or resource.
// Parameters
//	lastAsyncKey - last watermark
//	folderIds - IDs of subcribed folders
// Return
//	list - all found changes
//	asyncKey - new watermark
func (c *ClientConnection) ChangesGetAccount(lastAsyncKey AccountSyncKey, folderIds KIdList) (ChangeList, *AccountSyncKey, error) {
	params := struct {
		LastAsyncKey AccountSyncKey `json:"lastAsyncKey"`
		FolderIds    KIdList        `json:"folderIds"`
	}{lastAsyncKey, folderIds}
	data, err := c.CallRaw("Changes.getAccount", params)
	if err != nil {
		return nil, nil, err
	}
	list := struct {
		Result struct {
			List     ChangeList     `json:"list"`
			AsyncKey AccountSyncKey `json:"asyncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, &list.Result.AsyncKey, err
}

// ChangesKillRequest - Kill current running Changes.get's request. It supposed that timeout was specified > 0.
// Parameters
//	lastSyncKey - last watermark
func (c *ClientConnection) ChangesKillRequest(lastSyncKey SyncKey) error {
	params := struct {
		LastSyncKey SyncKey `json:"lastSyncKey"`
	}{lastSyncKey}
	_, err := c.CallRaw("Changes.killRequest", params)
	return err
}

// ChangesGetFolder - Get changes in a folder.
// Parameters
//	folderId - folder from which we want get item changes
//	lastSyncKey - last synckey (watermark)
// Return
//	list - all found changes
//	syncKey - new last synckey (watermark)
func (c *ClientConnection) ChangesGetFolder(folderId KId, lastSyncKey Watermark) (ChangeList, *Watermark, error) {
	params := struct {
		FolderId    KId       `json:"folderId"`
		LastSyncKey Watermark `json:"lastSyncKey"`
	}{folderId, lastSyncKey}
	data, err := c.CallRaw("Changes.getFolder", params)
	if err != nil {
		return nil, nil, err
	}
	list := struct {
		Result struct {
			List    ChangeList `json:"list"`
			SyncKey Watermark  `json:"syncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, &list.Result.SyncKey, err
}

// ChangesGetSyncKey - Get actual watermark.
func (c *ClientConnection) ChangesGetSyncKey() (*SyncKey, error) {
	data, err := c.CallRaw("Changes.getSyncKey", nil)
	if err != nil {
		return nil, err
	}
	syncKey := struct {
		Result struct {
			SyncKey SyncKey `json:"syncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &syncKey)
	return &syncKey.Result.SyncKey, err
}

// ChangesGetAccountSyncKey - Get actual watermark.
func (c *ClientConnection) ChangesGetAccountSyncKey(mailboxId KId) (*AccountSyncKey, error) {
	params := struct {
		MailboxId KId `json:"mailboxId"`
	}{mailboxId}
	data, err := c.CallRaw("Changes.getAccountSyncKey", params)
	if err != nil {
		return nil, err
	}
	asyncKey := struct {
		Result struct {
			AsyncKey AccountSyncKey `json:"asyncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &asyncKey)
	return &asyncKey.Result.AsyncKey, err
}

// ChangesGetFolderSyncKey - Get actual sync key for a folder.
// Parameters
//	folderId - wanted folder
// Return
//	syncKey - actual synckey (watermark) for folder
func (c *ClientConnection) ChangesGetFolderSyncKey(folderId KId) (*Watermark, error) {
	params := struct {
		FolderId KId `json:"folderId"`
	}{folderId}
	data, err := c.CallRaw("Changes.getFolderSyncKey", params)
	if err != nil {
		return nil, err
	}
	syncKey := struct {
		Result struct {
			SyncKey Watermark `json:"syncKey"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &syncKey)
	return &syncKey.Result.SyncKey, err
}
