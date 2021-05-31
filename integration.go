package webmail

import "encoding/json"

// SyncFolder - Class with methods for integration
type SyncFolder struct {
	Id           KId             `json:"id"`           // [READ-ONLY] global identification
	ParentId     KId             `json:"parentId"`     // [READ-ONLY] global identification
	Name         string          `json:"name"`         // [READ-ONLY] folder name displayed in folder tree
	Type         FolderType      `json:"type"`         // [READ-ONLY] type of the folder
	SubType      FolderSubType   `json:"subType"`      // [READ-ONLY] type of the folder
	PlaceType    FolderPlaceType `json:"placeType"`    // [READ-ONLY] type of place where is folder placed
	NestingLevel int             `json:"nestingLevel"` // [READ-ONLY] number 0 = root folder, 1 = subfolders of root folder, 2 = subfolder of subfolder, ...
	OwnerName    string          `json:"ownerName"`    // [READ-ONLY] name of owner of folder (available only for 'FPlacePeople', 'FPlaceResources' and 'FPlaceLocations')
	EmailAddress string          `json:"emailAddress"` // [READ-ONLY] email of owner of folder (available only for 'FPlacePeople', 'FPlaceResources' and 'FPlaceLocations')
	IsSelectable bool            `json:"isSelectable"` // [READ-ONLY] false if a setting of this folder cannot be modified
	Synchronize  bool            `json:"synchronize"`  // true if should be this folder synchronize
}

type SyncFolderList []SyncFolder

// Class with methods for integration

// IntegrationGetASyncFolderList - Obtain list of folders of currently logged user
// Return
//	list - list of folders
func (c *ClientConnection) IntegrationGetASyncFolderList() (SyncFolderList, error) {
	data, err := c.CallRaw("Integration.getASyncFolderList", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List SyncFolderList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// IntegrationSetASyncFolderList - Set folder properties
// Parameters
//	folders - properties to save
// Return
//	errors - error message list
func (c *ClientConnection) IntegrationSetASyncFolderList(folders SyncFolderList) (ErrorList, error) {
	params := struct {
		Folders SyncFolderList `json:"folders"`
	}{folders}
	data, err := c.CallRaw("Integration.setASyncFolderList", params)
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

// IntegrationGetIPhoneSyncFolderList - Obtain list of folders of currently logged user (task and calendars only)
// Return
//	list - list of folders
func (c *ClientConnection) IntegrationGetIPhoneSyncFolderList() (SyncFolderList, error) {
	data, err := c.CallRaw("Integration.getIPhoneSyncFolderList", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List SyncFolderList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// IntegrationSetIPhoneSyncFolderList - Set folder properties (task and calendars only)
// Parameters
//	folders - properties to save
// Return
//	errors - error message list
func (c *ClientConnection) IntegrationSetIPhoneSyncFolderList(folders SyncFolderList) (ErrorList, error) {
	params := struct {
		Folders SyncFolderList `json:"folders"`
	}{folders}
	data, err := c.CallRaw("Integration.setIPhoneSyncFolderList", params)
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
