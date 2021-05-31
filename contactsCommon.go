package webmail

// ABExtension - Extension of Apple Address Book
type ABExtension struct {
	GroupId string `json:"groupId"`
	Label   string `json:"label"`
}

// PostalAddressType - Type of a postal address
type PostalAddressType string

const (
	AddressHome   PostalAddressType = "AddressHome"
	AddressWork   PostalAddressType = "AddressWork"
	AddressOther  PostalAddressType = "AddressOther"
	AddressCustom PostalAddressType = "AddressCustom" // no type defined
)

// PostalAddress - Structure describing a postal address in contact.
type PostalAddress struct {
	Preferred       bool              `json:"preferred"`
	Pobox           string            `json:"pobox"`           // the post office box
	ExtendedAddress string            `json:"extendedAddress"` // e.g., apartment or suite number
	Street          string            `json:"street"`          // the street address
	Locality        string            `json:"locality"`        // the locality (e.g., city)
	State           string            `json:"state"`           // the region (e.g., state or province);
	Zip             string            `json:"zip"`             // the postal code
	Country         string            `json:"country"`         // the country name (full name)
	Label           string            `json:"label"`
	Type            PostalAddressType `json:"type"`
	Extension       ABExtension       `json:"extension"`
}

// PostalAddressList - A sequence of postal addresses.
type PostalAddressList []PostalAddress

// EmailAddressType - Email address type.
type EmailAddressType string

const (
	EmailWork   EmailAddressType = "EmailWork"
	EmailHome   EmailAddressType = "EmailHome"
	EmailOther  EmailAddressType = "EmailOther"
	EmailCustom EmailAddressType = "EmailCustom" // no type defined
	// valid for distribution lists only
	RefContact          EmailAddressType = "RefContact"          // Reference to existing conatact
	RefDistributionList EmailAddressType = "RefDistributionList" // Reference to existing distribution list
)

// EmailAddress - Structure describing an email address in contact.
type EmailAddress struct {
	Address            string           `json:"address"`
	Name               string           `json:"name"`
	Preferred          bool             `json:"preferred"`
	IsValidCertificate bool             `json:"isValidCertificate"`
	Type               EmailAddressType `json:"type"`
	RefId              KId              `json:"refId"` // Global identification of reference. Valid for types 'RefContact' and 'RefDistributionList'.
	Extension          ABExtension      `json:"extension"`
}

// EmailAddressList - Sequence of email addresses.
type EmailAddressList []EmailAddress

// PhoneNumberType - Type of a contact phone number
type PhoneNumberType string

const (
	TypeAssistant  PhoneNumberType = "TypeAssistant"
	TypeWorkVoice  PhoneNumberType = "TypeWorkVoice"
	TypeWorkFax    PhoneNumberType = "TypeWorkFax"
	TypeCallback   PhoneNumberType = "TypeCallback"
	TypeCar        PhoneNumberType = "TypeCar"
	TypeCompany    PhoneNumberType = "TypeCompany"
	TypeHomeVoice  PhoneNumberType = "TypeHomeVoice"
	TypeHomeFax    PhoneNumberType = "TypeHomeFax"
	TypeIsdn       PhoneNumberType = "TypeIsdn"
	TypeMobile     PhoneNumberType = "TypeMobile"
	TypeOtherVoice PhoneNumberType = "TypeOtherVoice"
	TypeOtherFax   PhoneNumberType = "TypeOtherFax"
	TypePager      PhoneNumberType = "TypePager"
	TypePrimary    PhoneNumberType = "TypePrimary"
	TypeRadio      PhoneNumberType = "TypeRadio"
	TypeTelex      PhoneNumberType = "TypeTelex"
	TypeTtyTdd     PhoneNumberType = "TypeTtyTdd"
	TypeCustom     PhoneNumberType = "TypeCustom" // no type defined
)

// PhoneNumber - Structure desribing a contact phone number
type PhoneNumber struct {
	Type      PhoneNumberType `json:"type"`
	Number    string          `json:"number"` // A number - based on the X.500 Telephone Number attribute
	Extension ABExtension     `json:"extension"`
}

type PhoneNumberList []PhoneNumber

// UrlType - Type of URL
type UrlType string

const (
	UrlHome   UrlType = "UrlHome"
	UrlWork   UrlType = "UrlWork"
	UrlOther  UrlType = "UrlOther"
	UrlCustom UrlType = "UrlCustom" // no type defined
)

// Url - Structure desribing URL
type Url struct {
	Type      UrlType     `json:"type"`
	Url       string      `json:"url"`
	Extension ABExtension `json:"extension"`
}

type UrlList []Url

// PhotoAttachment - A contact photo. Only JPEG format is supported. Maximum size is 256 kB.
type PhotoAttachment struct {
	Id  string `json:"id"`  // origin ID of attachment or ID from upload response
	Url string `json:"url"` // [READ-ONLY] Relative URL from root of web. Eg.: /webmail/api/download/attachment/ba5767a9-7a70-4c90-a6bf-dc8dd62e259c/14/0-1-0-1/picture.jpg
}

// PersonalContact - Personal Contact detail.
type PersonalContact struct {
	CommonName           string           `json:"commonName"`
	FirstName            string           `json:"firstName"`
	MiddleName           string           `json:"middleName"`
	SurName              string           `json:"surName"`
	TitleBefore          string           `json:"titleBefore"`
	TitleAfter           string           `json:"titleAfter"`
	NickName             string           `json:"nickName"`
	PhoneNumberWorkVoice string           `json:"phoneNumberWorkVoice"`
	PhoneNumberMobile    string           `json:"phoneNumberMobile"`
	PostalAddressWork    PostalAddress    `json:"postalAddressWork"`
	UrlWork              string           `json:"urlWork"`
	BirthDay             UtcDateTime      `json:"birthDay"`
	Anniversary          UtcDateTime      `json:"anniversary"`
	CompanyName          string           `json:"companyName"`
	DepartmentName       string           `json:"departmentName"`
	Profession           string           `json:"profession"`
	ManagerName          string           `json:"managerName"`
	AssistantName        string           `json:"assistantName"`
	Comment              string           `json:"comment"`
	IMAddress            string           `json:"IMAddress"`
	Photo                PhotoAttachment  `json:"photo"`
	IsReadOnly           bool             `json:"isReadOnly"`
	EmailAddresses       EmailAddressList `json:"emailAddresses"`
}

type PersonalContactList []PersonalContact
