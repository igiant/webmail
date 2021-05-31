package webmail

import "encoding/json"

// UserInfo - Details of the logged user into the webmail.
type UserInfo struct {
	Id               KId        `json:"id"`               // [READ-ONLY] global id of user
	LoginName        string     `json:"loginName"`        // [READ-ONLY] loginName name, also primary email address
	FullName         string     `json:"fullName"`         // [READ-ONLY] full name
	Emails           StringList `json:"emails"`           // [READ-ONLY]
	PreferredAddress string     `json:"preferredAddress"` // preferred address
	ReplyToAddress   string     `json:"replyToAddress"`   // address for reply
}

type OutOfOfficeSettings struct {
	IsEnabled          bool        `json:"isEnabled"`
	Text               string      `json:"text"`
	IsTimeRangeEnabled bool        `json:"isTimeRangeEnabled"`
	TimeRangeStart     UtcDateTime `json:"timeRangeStart"`
	TimeRangeEnd       UtcDateTime `json:"timeRangeEnd"`
}

type SpamSettings struct {
	IsEnabled           bool       `json:"isEnabled"`           // If enabled then spam is moved to the Junk E-mail folder
	WhiteListContacts   bool       `json:"whiteListContacts"`   // Also trust senders from Contacts folder
	AutoupdateWhiteList bool       `json:"autoupdateWhiteList"` // If enable e-mail address of original sender will be added into white list while sending reply
	WhiteList           StringList `json:"whiteList"`           // Trust to these senders
}

// QuotaInfo - Stores user's quota info
type QuotaInfo struct {
	MessagesLimit          int    `json:"messagesLimit"`          // Maximum number of messages that current user is allowed to have, value 0 means user has no limit
	MessagesUsed           int    `json:"messagesUsed"`           // Number of messages that curent user has
	SpaceLimit             uint64 `json:"spaceLimit"`             // Maximum amount of space [Bytes] reserved for current user, value 0 means user has no limit
	SpaceUsed              uint64 `json:"spaceUsed"`              // Amount of space [Bytes] consumed by current user
	PercentLimitForWarning int    `json:"percentLimitForWarning"` // Value in percent that user has to exceed to be warned
}

// PasswordPolicy - If this policy is enabled, passwords must meet the following minimum requirements when they are changed or created:
// - Passwords must not contain the user's entire and checks are not case sensitive
// - Passwords must contain characters from three of the following five categories:
// - Uppercase characters of European languages (A through Z, with diacritic marks, Greek and Cyrillic characters)
// - Lowercase characters of European languages (a through z, sharp-s, with diacritic marks, Greek and Cyrillic characters)
// - Base 10 digits (0 through 9)
// - Nonalphanumeric characters: ~!@#$%^&*_-+=`|\(){}[]:;"'<>,.?/
// - Any Unicode character that is categorized as an alphabetic character but is not uppercase or lowercase. This includes Unicode characters from Asian languages.
type PasswordPolicy struct {
	IsEnabled bool `json:"isEnabled"`
	MinLength int  `json:"minLength"`
}

// Currently logged user manager

// SessionCanUserChangePassword - to change his/her password.
// Return
//	isEligible - is set to true as long as user is eligible
func (c *ClientConnection) SessionCanUserChangePassword() (bool, error) {
	data, err := c.CallRaw("Session.canUserChangePassword", nil)
	if err != nil {
		return false, err
	}
	isEligible := struct {
		Result struct {
			IsEligible bool `json:"isEligible"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &isEligible)
	return isEligible.Result.IsEligible, err
}

// SessionGetAvailableTimeZones - Get list of all available time zones.
// Return
//	zones - list of time zones
func (c *ClientConnection) SessionGetAvailableTimeZones() (StringList, error) {
	data, err := c.CallRaw("Session.getAvailableTimeZones", nil)
	if err != nil {
		return nil, err
	}
	zones := struct {
		Result struct {
			Zones StringList `json:"zones"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &zones)
	return zones.Result.Zones, err
}

// SessionGetAvailableLanguages - Get list of all languages supported by server.
func (c *ClientConnection) SessionGetAvailableLanguages() (LangDescriptionList, error) {
	data, err := c.CallRaw("Session.getAvailableLanguages", nil)
	if err != nil {
		return nil, err
	}
	languages := struct {
		Result struct {
			Languages LangDescriptionList `json:"languages"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &languages)
	return languages.Result.Languages, err
}

// SessionGetOutOfOffice - Obtain the Auto Reply settings
// Return
//	settings - details
func (c *ClientConnection) SessionGetOutOfOffice() (*OutOfOfficeSettings, error) {
	data, err := c.CallRaw("Session.getOutOfOffice", nil)
	if err != nil {
		return nil, err
	}
	settings := struct {
		Result struct {
			Settings OutOfOfficeSettings `json:"settings"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &settings)
	return &settings.Result.Settings, err
}

// SessionGetQuotaInformation - Obtain iformations about quota of current user.
func (c *ClientConnection) SessionGetQuotaInformation() (*QuotaInfo, error) {
	data, err := c.CallRaw("Session.getQuotaInformation", nil)
	if err != nil {
		return nil, err
	}
	quotaInfo := struct {
		Result struct {
			QuotaInfo QuotaInfo `json:"quotaInfo"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &quotaInfo)
	return &quotaInfo.Result.QuotaInfo, err
}

// SessionGetSettings - Obtain currently logged user's settings.
// Return
//	settings - WAM settings
func (c *ClientConnection) SessionGetSettings(query SettingQuery) (*jsonstring, error) {
	params := struct {
		Query SettingQuery `json:"query"`
	}{query}
	data, err := c.CallRaw("Session.getSettings", params)
	if err != nil {
		return nil, err
	}
	settings := struct {
		Result struct {
			Settings jsonstring `json:"settings"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &settings)
	return &settings.Result.Settings, err
}

// SessionGetSpamSettings - Obtain the spam settings
// Return
//	settings - details
func (c *ClientConnection) SessionGetSpamSettings() (*SpamSettings, error) {
	data, err := c.CallRaw("Session.getSpamSettings", nil)
	if err != nil {
		return nil, err
	}
	settings := struct {
		Result struct {
			Settings SpamSettings `json:"settings"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &settings)
	return &settings.Result.Settings, err
}

// SessionGetUserVoiceUrl - Obtain URL for users' access to UserVoice
// Return
//	accessUrl - URL for access to UserVoice
func (c *ClientConnection) SessionGetUserVoiceUrl() (string, error) {
	data, err := c.CallRaw("Session.getUserVoiceUrl", nil)
	if err != nil {
		return "", err
	}
	accessUrl := struct {
		Result struct {
			AccessUrl string `json:"accessUrl"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &accessUrl)
	return accessUrl.Result.AccessUrl, err
}

// Login - [KLoginMethod]
// Parameters
//	userName
//	password
//	application - application descriminator, note that with session to admin you cannot log in webmail
// Return
//	token
func (c *ClientConnection) Login(userName string, password string, app *ApiApplication) error {
	if app == nil {
		app = NewApplication("", "", "")
	}
	params := loginStruct{userName, password, *app}
	data, err := c.CallRaw("Session.login", params)
	if err != nil {
		return err
	}
	token := struct {
		Result struct {
			Token string `json:"token"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &token)
	if err != nil {
		return err
	}
	c.Token = &token.Result.Token
	return nil
}

// Logout - [KLogoutMethod]
func (c *ClientConnection) Logout() error {
	_, err := c.CallRaw("Session.logout", nil)
	return err
}

// SessionSetOutOfOffice - Set the Auto Reply settings
// Parameters
//	settings - details
func (c *ClientConnection) SessionSetOutOfOffice(settings OutOfOfficeSettings) error {
	params := struct {
		Settings OutOfOfficeSettings `json:"settings"`
	}{settings}
	_, err := c.CallRaw("Session.setOutOfOffice", params)
	return err
}

// SessionSetPassword - Change password of current user.
// Parameters
//	currentPassword - current users' password
//	newPassword - new users' password
func (c *ClientConnection) SessionSetPassword(currentPassword string, newPassword string) error {
	params := struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}{currentPassword, newPassword}
	_, err := c.CallRaw("Session.setPassword", params)
	return err
}

// SessionSetSettings - Set settings of the currently logged user.
// Parameters
//	settings - WAM settings
func (c *ClientConnection) SessionSetSettings(settings jsonstring) error {
	params := struct {
		Settings jsonstring `json:"settings"`
	}{settings}
	_, err := c.CallRaw("Session.setSettings", params)
	return err
}

// SessionSetSpamSettings - Set the spam settings
// Parameters
//	settings - details
func (c *ClientConnection) SessionSetSpamSettings(settings SpamSettings) error {
	params := struct {
		Settings SpamSettings `json:"settings"`
	}{settings}
	_, err := c.CallRaw("Session.setSpamSettings", params)
	return err
}

// SessionSetUserInfo - Set user details.
// Parameters
//	userDetails - details about the currently logged user
func (c *ClientConnection) SessionSetUserInfo(userDetails UserInfo) error {
	params := struct {
		UserDetails UserInfo `json:"userDetails"`
	}{userDetails}
	_, err := c.CallRaw("Session.setUserInfo", params)
	return err
}

// SessionWhoAmI - Determines the currently logged user (caller).
// Return
//	userDetails - details about the currently logged user
func (c *ClientConnection) SessionWhoAmI() (*UserInfo, error) {
	data, err := c.CallRaw("Session.whoAmI", nil)
	if err != nil {
		return nil, err
	}
	userDetails := struct {
		Result struct {
			UserDetails UserInfo `json:"userDetails"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &userDetails)
	return &userDetails.Result.UserDetails, err
}

// SessionGetMobileDeviceList - Obtain a list of mobile devices of given user.
// Parameters
//	query - query attributes and limits
// Return
//	list - mobile devices of given user
//	totalItems - number of mobile devices found for given user
func (c *ClientConnection) SessionGetMobileDeviceList(query SearchQuery) (MobileDeviceList, int, error) {
	params := struct {
		Query SearchQuery `json:"query"`
	}{query}
	data, err := c.CallRaw("Session.getMobileDeviceList", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       MobileDeviceList `json:"list"`
			TotalItems int              `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// SessionRemoveMobileDevice - Remove mobile device from the list of user's mobile devices.
// Parameters
//	deviceId - ID of user's mobile device to be removed
func (c *ClientConnection) SessionRemoveMobileDevice(deviceId string) error {
	params := struct {
		DeviceId string `json:"deviceId"`
	}{deviceId}
	_, err := c.CallRaw("Session.removeMobileDevice", params)
	return err
}

// SessionWipeMobileDevice - Wipe user's mobile device.
// Parameters
//	deviceId - ID of user's mobile device to be wiped
//	password - password of current user
func (c *ClientConnection) SessionWipeMobileDevice(deviceId string, password string) error {
	params := struct {
		DeviceId string `json:"deviceId"`
		Password string `json:"password"`
	}{deviceId, password}
	_, err := c.CallRaw("Session.wipeMobileDevice", params)
	return err
}

// SessionCancelWipeMobileDevice - Cancel wiping of user's mobile device.
// Parameters
//	deviceId - ID of user's mobile device to cancel wipe
func (c *ClientConnection) SessionCancelWipeMobileDevice(deviceId string) error {
	params := struct {
		DeviceId string `json:"deviceId"`
	}{deviceId}
	_, err := c.CallRaw("Session.cancelWipeMobileDevice", params)
	return err
}

// SessionGetSignatureImageList - Obtain list of images stored in user account
func (c *ClientConnection) SessionGetSignatureImageList() (ImageList, error) {
	data, err := c.CallRaw("Session.getSignatureImageList", nil)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List ImageList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// SessionAddSignatureImage - Add image into user's store
// Parameters
//	ids - Upload IDs of images to add into user's store
// Return
//	errors - list of errors
//	result - succesfuly added images
func (c *ClientConnection) SessionAddSignatureImage(ids KIdList) (ErrorList, ImageList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Session.addSignatureImage", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result ImageList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// SessionRemoveSignatureImage - Remove image from user's store
// Parameters
//	ids - Image IDs to remove
func (c *ClientConnection) SessionRemoveSignatureImage(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Session.removeSignatureImage", params)
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
