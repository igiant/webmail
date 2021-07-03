package webmail

import "encoding/json"

type ContactType string

const (
	ctContact          ContactType = "ctContact"
	ctDistributionList ContactType = "ctDistributionList"
)

// Contact - Contact detail.
type Contact struct {
	Id              KId               `json:"id"`
	FolderId        KId               `json:"folderId"`
	Watermark       Watermark         `json:"watermark"`
	Type            ContactType       `json:"type"`
	CommonName      string            `json:"commonName"`
	FirstName       string            `json:"firstName"`
	MiddleName      string            `json:"middleName"`
	SurName         string            `json:"surName"`
	TitleBefore     string            `json:"titleBefore"`
	TitleAfter      string            `json:"titleAfter"`
	NickName        string            `json:"nickName"`
	PhoneNumbers    PhoneNumberList   `json:"phoneNumbers"`
	EmailAddresses  EmailAddressList  `json:"emailAddresses"`
	PostalAddresses PostalAddressList `json:"postalAddresses"`
	Urls            UrlList           `json:"urls"`
	BirthDay        UtcDateTime       `json:"birthDay"`
	Anniversary     UtcDateTime       `json:"anniversary"`
	CompanyName     string            `json:"companyName"`
	DepartmentName  string            `json:"departmentName"`
	Profession      string            `json:"profession"`
	ManagerName     string            `json:"managerName"`
	AssistantName   string            `json:"assistantName"`
	Comment         string            `json:"comment"`
	IMAddress       string            `json:"IMAddress"`
	Photo           PhotoAttachment   `json:"photo"`
	Categories      StringList        `json:"categories"`
	CertSourceId    KId               `json:"certSourceId"` // [WRITE-ONLY]
	IsGalContact    bool              `json:"isGalContact"` // [READ-ONLY]
}

type ContactList []Contact

// ResourceType - Export format type
type ResourceType string

const (
	ResourceRoom      ResourceType = "ResourceRoom"      // resource is a room
	ResourceEquipment ResourceType = "ResourceEquipment" // resource is something else, eg: a car
)

// Resource - Resource details [READ-ONLY]
type Resource struct {
	Name        string       `json:"name"`        // resource name
	Address     string       `json:"address"`     // email of resource
	Description string       `json:"description"` // resource description
	Type        ResourceType `json:"type"`        // type of the resource
}

// ResourceList - List of resources
type ResourceList []Resource

// Constants for composing kerio::web::SearchQuery
// Contacts management.

// ContactsCopy - Copy existing contacts to folder
// Parameters
//	ids - list of global identifiers of contacts to be copied
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) ContactsCopy(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Contacts.copy", params)
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

// ContactsCreate - Create contact in particular folder
// Parameters
//	contacts - new contacts; Field 'folderId' must be set.
// Return
//	errors - error message list
//	result - list of ID of crated contacts
func (c *ClientConnection) ContactsCreate(contacts ContactList) (ErrorList, CreateResultList, error) {
	params := struct {
		Contacts ContactList `json:"contacts"`
	}{contacts}
	data, err := c.CallRaw("Contacts.create", params)
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

// ContactsGet - Get a list of contacts.
// Parameters
//	folderIds - list of global identifiers of folders to be listed.
//	query - query attributes and limits
// Return
//	list - all found contacts
//  totalItems - number of contacts found if there is no limit
func (c *ClientConnection) ContactsGet(folderIds KIdList, query SearchQuery) (ContactList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Contacts.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       ContactList `json:"list"`
			TotalItems int         `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// ContactsGetFromCache - id, folderId, watermark, type, commonName, titleAfter, titleBefore, firstName, middleName, surName, nickName, emailAddresses, phoneNumbers, photo
// Parameters
//	folderIds - list of global identifiers of folders to be listed.
//	query - query attributes and limits
// Return
//	list - all found contacts
//  totalItems - number of contacts found if there is no limit
func (c *ClientConnection) ContactsGetFromCache(folderIds KIdList, query SearchQuery) (ContactList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Contacts.getFromCache", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       ContactList `json:"list"`
			TotalItems int         `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// ContactsGetById - Get particular contacts. All members of struct Contact are filed in response.
// Parameters
//	ids - global identifiers of contact.
// Return
//	errors - list of errors which happened
//	result - contacts of given IDs. All members of struct are returned.
func (c *ClientConnection) ContactsGetById(ids KIdList) (ErrorList, ContactList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Contacts.getById", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList   `json:"errors"`
			Result ContactList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// ContactsGetByIdFromCache - id, folderId, watermark, type, commonName, titleAfter, titleBefore, firstName, middleName, surName, nickName, emailAddresses, phoneNumbers, photo
// Parameters
//	ids - global identifiers of contact.
// Return
//	errors - list of errors which happened
//	result - contacts of given IDs.
func (c *ClientConnection) ContactsGetByIdFromCache(ids KIdList) (ErrorList, ContactList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Contacts.getByIdFromCache", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList   `json:"errors"`
			Result ContactList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// ContactsGetFromAttachment - Get contact from attachment.
// Parameters
//	id - global identifiers of mail attachment.
// Return
//	result - contact of given IDs. All members of struct are returned.
func (c *ClientConnection) ContactsGetFromAttachment(id KId) (*Contact, error) {
	params := struct {
		Id KId `json:"id"`
	}{id}
	data, err := c.CallRaw("Contacts.getFromAttachment", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Result struct {
			Result Contact `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	return &result.Result.Result, err
}

// ContactsGetResources - Get a list of resources that an user can schedule.
// Parameters
//	query - query attributes and limits (empty query obtain all resources)
// Return
//	list - all found resources
//  totalItems - number of resources found if there is no limit
func (c *ClientConnection) ContactsGetResources(query SearchQuery) (ResourceList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		Query SearchQuery `json:"query"`
	}{query}
	data, err := c.CallRaw("Contacts.getResources", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       ResourceList `json:"list"`
			TotalItems int          `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// ContactsGetCertificate - Get a certificate for given email address.
// Parameters
//	email - email address of requested certificate
//	id - global identifier of contacts to be searched
// Return
//	cert - found certificate
func (c *ClientConnection) ContactsGetCertificate(email string, id KId) (*Certificate, error) {
	params := struct {
		Email string `json:"email"`
		Id    KId    `json:"id"`
	}{email, id}
	data, err := c.CallRaw("Contacts.getCertificate", params)
	if err != nil {
		return nil, err
	}
	cert := struct {
		Result struct {
			Cert Certificate `json:"cert"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &cert)
	return &cert.Result.Cert, err
}

// ContactsRemove - Remove a list of contacts.
// Parameters
//	ids - list of global identifiers of contacts to be removed
// Return
//	errors - list of contacts that failed to remove
func (c *ClientConnection) ContactsRemove(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Contacts.remove", params)
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

// ContactsSet - Set existing contacts.
// Parameters
//	contacts - modifications of contacts. Field 'folderId' must be set.
// Return
//	errors - error message list
func (c *ClientConnection) ContactsSet(contacts ContactList) (ErrorList, SetResultList, error) {
	params := struct {
		Contacts ContactList `json:"contacts"`
	}{contacts}
	data, err := c.CallRaw("Contacts.set", params)
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

// ContactsMove - Move existing contacts to folder
// Parameters
//	ids - list of global identifiers of contacts to be moved
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) ContactsMove(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Contacts.move", params)
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

// ContactsGetPersonal - Get personal user contact
func (c *ClientConnection) ContactsGetPersonal() (*PersonalContact, error) {
	data, err := c.CallRaw("Contacts.getPersonal", nil)
	if err != nil {
		return nil, err
	}
	contact := struct {
		Result struct {
			Contact PersonalContact `json:"contact"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &contact)
	return &contact.Result.Contact, err
}

// ContactsSetPersonal - Set personal user contact
func (c *ClientConnection) ContactsSetPersonal(contact PersonalContact) error {
	params := struct {
		Contact PersonalContact `json:"contact"`
	}{contact}
	_, err := c.CallRaw("Contacts.setPersonal", params)
	return err
}
