package webmail

// FolderSubType - Folder sub-type enumeration.
type FolderSubType string
const (
	FSubNone FolderSubType = "FSubNone" // Ordinary folder.
	FSubInbox FolderSubType = "FSubInbox" // Special sub-type of FMail. This folder cannot be removed.
	FSubDrafts FolderSubType = "FSubDrafts" // Special sub-type of FMail. This folder cannot be removed.
	FSubSentItems FolderSubType = "FSubSentItems" // Special sub-type of FMail. This folder cannot be removed.
	FSubDeletedItems FolderSubType = "FSubDeletedItems" // Special sub-type of FMail. This folder cannot be removed.
	FSubJunkEmail FolderSubType = "FSubJunkEmail" // Special sub-type of FMail. This folder cannot be removed.
	FSubDefault FolderSubType = "FSubDefault" // These folders cannot be removed.
	FSubGalResources FolderSubType = "FSubGalResources" // This folder is created via GAL for storage contacts of resources.
	FSubGalContacts FolderSubType = "FSubGalContacts" // This folder is created via GAL for storage contacts of people.
)

// FolderType - Folder type enumeration
type FolderType string
const (
	FRoot FolderType = "FRoot" 
	FMail FolderType = "FMail" 
	FContact FolderType = "FContact" 
	FCalendar FolderType = "FCalendar" 
	FTask FolderType = "FTask" 
	FNote FolderType = "FNote" 
)

// FolderPlaceType - Type of place where is folder placed
type FolderPlaceType string
const (
	FPlaceMailbox FolderPlaceType = "FPlaceMailbox" // the mailbox of currently loged user
	FPlaceResources FolderPlaceType = "FPlaceResources" // the resource type of Equipment
	FPlaceLocations FolderPlaceType = "FPlaceLocations" // the resource type of Room
	FPlacePeople FolderPlaceType = "FPlacePeople" // the shared folder of another user
	FPlacePublic FolderPlaceType = "FPlacePublic" // the public folder
	FPlaceArchive FolderPlaceType = "FPlaceArchive" // the archive folder
)

// FolderAccess - Access to folder
type FolderAccess string
const (
	FAccessListingOnly FolderAccess = "FAccessListingOnly" 
	FAccessReadOnly FolderAccess = "FAccessReadOnly" 
	FAccessReadWrite FolderAccess = "FAccessReadWrite" 
	FAccessAdmin FolderAccess = "FAccessAdmin" // full access; E.g user can add folder
)

type FolderPermission struct {
	Access FolderAccess `json:"access"` 
	Principal kerio::jsonapi::webmail::principals::Principal `json:"principal"` 
	Inherited bool `json:"inherited"` // [READ-ONLY] permission are placed in a public root folder and there are read-only here
	IsDelegatee bool `json:"isDelegatee"` // [READ-ONLY] principal is delegatee (this flag is filled only for default calendar and INBOX otherwise is false)
}

type FolderPermissionList []FolderPermission

type Folder struct {
	Id KId `json:"id"` // [READ-ONLY] global identification
	ParentId KId `json:"parentId"` // global identification
	Name string `json:"name"` // folder name displayed in folder tree
	OwnerName string `json:"ownerName"` // [READ-ONLY] name of owner of folder (available only for 'FPlacePeople', 'FPlaceResources' and 'FPlaceLocations')
	EmailAddress string `json:"emailAddress"` // [READ-ONLY] email of owner of folder (available only for 'FPlacePeople', 'FPlaceResources' and 'FPlaceLocations')
	Type FolderType `json:"type"` // type of the folder
	SubType FolderSubType `json:"subType"` // [READ-ONLY] type of the folder
	PlaceType FolderPlaceType `json:"placeType"` // [READ-ONLY] type of place where is folder placed
	Access FolderAccess `json:"access"` // [READ-ONLY] type of access of currently loged user
	IsShared bool `json:"isShared"` // [READ-ONLY] true if a folder is shared to another user (permissions are not empty)
	IsDelegated bool `json:"isDelegated"` // [READ-ONLY] true if a folder access is R/W and user is a delegate
	NestingLevel int `json:"nestingLevel"` // [READ-ONLY] number 0 = root folder, 1 = subfolders of root folder, 2 = subfolder of subfolder, ...
	MessageCount int `json:"messageCount"` // [READ-ONLY] count of items in the folder (not set for root folder)
	MessageUnread int `json:"messageUnread"` // [READ-ONLY] count of unread items (mails or deleted items) in folder, not set for non-mail folder types
	MessageSize longlong `json:"messageSize"` // [READ-ONLY] size of all messages in the folder (without subdirectories)
	Checked bool `json:"checked"` // true if a folder is chosen to view
	Color string `json:"color"` // a color of folder, if string is empty no color is set
	Published bool `json:"published"` // [READ-ONLY] true, if folder was published to server
}

type FolderList []Folder

type SharedMailbox struct {
	Principal kerio::jsonapi::webmail::principals::Principal `json:"principal"` // types [ptUser, ptResource]
	MailboxId KId `json:"mailboxId"` // root folder ID
	IsLoaded bool `json:"isLoaded"` // folders are loaded (are present on the same home server)
	Folders FolderList `json:"folders"` // folders with at least listing-only righs
	SubscribedFolderIds KIdList `json:"subscribedFolderIds"` 
}

type SharedMailboxList []SharedMailbox

// Folder store manager class

// FoldersClearToItemId - Remove all items in folder older then given item ID including. If folder is type of 'FMail' the items are moved to Delete Items or throw away if folder is Delete Items.
// Parameters
//	itemId - the last item ID
func (c *ClientConnection) FoldersClearToItemId(itemId KId) error {
	params := struct {
		ItemId KId `json:"itemId"`
	}{itemId}
	_, err := c.CallRaw("Folders.clearToItemId", params)
	return err
}

// FoldersCreate - Create new folders
// Parameters
//	folders - list of folders to create
// Return
//	errors - error message list
//	result - list of ID of crated folders.
func (c *ClientConnection) FoldersCreate(folders FolderList) (ErrorList, CreateResultList, error) {
	params := struct {
		Folders FolderList `json:"folders"`
	}{folders}
	data, err := c.CallRaw("Folders.create", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result CreateResultList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}


// FoldersGet - Obtain list of folders of currently logged user
// Return
//	list - list of folders
func (c *ClientConnection) FoldersGet() (FolderList, error) {
	data, err := c.CallRaw("Folders.get", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List FolderList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// FoldersGetShared - Obtain list of folders which currently logged user can access
// Return
//	list - list of folders
func (c *ClientConnection) FoldersGetShared(mailboxId KId) (FolderList, error) {
	params := struct {
		MailboxId KId `json:"mailboxId"`
	}{mailboxId}
	data, err := c.CallRaw("Folders.getShared", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List FolderList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// FoldersGetPublic - Obtain list of public folders which currently logged user can access
// Return
//	list - list of public folders
func (c *ClientConnection) FoldersGetPublic() (FolderList, error) {
	data, err := c.CallRaw("Folders.getPublic", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List FolderList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// FoldersGetSubscribed - Obtain list of folders acording SubscriptionList.
// Return
//	list - list of folders
func (c *ClientConnection) FoldersGetSubscribed() (SharedMailboxList, error) {
	data, err := c.CallRaw("Folders.getSubscribed", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List SharedMailboxList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}


// FoldersGetAutoCompleteContactsFolderId - Obtain ID of special folder for auto-complete contacts
// Return
//	folderId - ID of special folder
func (c *ClientConnection) FoldersGetAutoCompleteContactsFolderId() (*KId, error) {
	data, err := c.CallRaw("Folders.getAutoCompleteContactsFolderId", nil)
	if err != nil {
		return nil, err
	}
	folderId := struct {
		Result struct {
			FolderId KId `json:"folderId"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &folderId)
	return &folderId.Result.FolderId, err
}


// FoldersGetSharedMailboxList - Obtain list of mailboxes with their folders which currently logged user can access
// Return
//	mailboxes - list of mailboxes with their folders
func (c *ClientConnection) FoldersGetSharedMailboxList() (SharedMailboxList, error) {
	data, err := c.CallRaw("Folders.getSharedMailboxList", nil)
	if err != nil {
		return nil, err
	}
	mailboxes := struct {
		Result struct {
			Mailboxes SharedMailboxList `json:"mailboxes"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &mailboxes)
	return mailboxes.Result.Mailboxes, err
}


// FoldersMoveByType - Take a note that mail folders are moved recursively (the whole subtree)! Folders of other types (e.g. calendars) are not moved recursively.
// Parameters
//	targetId - target folder ID
//	ids - folder IDs
// Return
//	errors - error message list
func (c *ClientConnection) FoldersMoveByType(targetId KId, ids KIdList) (ErrorList, error) {
	params := struct {
		TargetId KId `json:"targetId"`
		Ids KIdList `json:"ids"`
	}{targetId, ids}
	data, err := c.CallRaw("Folders.moveByType", params)
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


// FoldersSet - Set folder properties
// Parameters
//	folders - properties to save
// Return
//	errors - error message list
func (c *ClientConnection) FoldersSet(folders FolderList) (ErrorList, error) {
	params := struct {
		Folders FolderList `json:"folders"`
	}{folders}
	data, err := c.CallRaw("Folders.set", params)
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


// FoldersRemove - Remove folder. Sub-folders are removed if recursive is true.
// Parameters
//	ids - folder IDs
//	recursive - remove sub-folders
// Return
//	errors - error message list
func (c *ClientConnection) FoldersRemove(ids KIdList, recursive bool) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
		Recursive bool `json:"recursive"`
	}{ids, recursive}
	data, err := c.CallRaw("Folders.remove", params)
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


// FoldersRemoveByType - Take a note that mail folders are removed recursively! Folders of other types (e.g. calendars) are not removed recursively.
// Parameters
//	ids - folder IDs
// Return
//	errors - error message list
func (c *ClientConnection) FoldersRemoveByType(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Folders.removeByType", params)
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


// FoldersGetPermissions - Get sharing permissions
// Parameters
//	folderId - ID of folder
// Return
//	permissions - sharing settings
func (c *ClientConnection) FoldersGetPermissions(folderId KId) (FolderPermissionList, error) {
	params := struct {
		FolderId KId `json:"folderId"`
	}{folderId}
	data, err := c.CallRaw("Folders.getPermissions", params)
	if err != nil {
		return nil, err
	}
	permissions := struct {
		Result struct {
			Permissions FolderPermissionList `json:"permissions"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &permissions)
	return permissions.Result.Permissions, err
}


// FoldersSetPermissions - Set sharing permissions
// Parameters
//	permissions - sharing settings
//	folderId - ID of folder
func (c *ClientConnection) FoldersSetPermissions(permissions FolderPermissionList, folderId KId, recursive bool) error {
	params := struct {
		Permissions FolderPermissionList `json:"permissions"`
		FolderId KId `json:"folderId"`
		Recursive bool `json:"recursive"`
	}{permissions, folderId, recursive}
	_, err := c.CallRaw("Folders.setPermissions", params)
	return err
}

// FoldersGetSubscriptionList - Get list of subscribed folders
func (c *ClientConnection) FoldersGetSubscriptionList() (KIdList, error) {
	data, err := c.CallRaw("Folders.getSubscriptionList", nil)
	if err != nil {
		return nil, err
	}
	folderIds := struct {
		Result struct {
			FolderIds KIdList `json:"folderIds"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &folderIds)
	return folderIds.Result.FolderIds, err
}


// FoldersSetSubscriptionList - Set list of subscribed folders
func (c *ClientConnection) FoldersSetSubscriptionList(folderIds KIdList) error {
	params := struct {
		FolderIds KIdList `json:"folderIds"`
	}{folderIds}
	_, err := c.CallRaw("Folders.setSubscriptionList", params)
	return err
}

// FoldersCopyAllMessages - Copies (or moves) all messages from the source folder to the destination
// Parameters
//	sourceId - ID of the source folder
//	destId - ID of the destionation folder
//	doMove - if true move the messages instead of copy the
func (c *ClientConnection) FoldersCopyAllMessages(sourceId KId, destId KId, doMove bool) error {
	params := struct {
		SourceId KId `json:"sourceId"`
		DestId KId `json:"destId"`
		DoMove bool `json:"doMove"`
	}{sourceId, destId, doMove}
	_, err := c.CallRaw("Folders.copyAllMessages", params)
	return err
}