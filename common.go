package webmail

type UtcDateTime string

type LongNumber uint64

// DateTimeStampList - Type for lists of date/times
type DateTimeStampList []DateTimeStamp

type SettingPath []string

type SettingQuery []SettingPath

type IdEntity struct {
	Id   KId    `json:"id"`   // global identifier of entity
	Name string `json:"name"` // [READ-ONLY] name or description of entity
}

type LangDescription struct {
	Name      string `json:"name"`      // name of language (national form)
	Code      string `json:"code"`      // code of language; E.g.: "en-gb"
	ShortCode string `json:"shortCode"` // short code of language which is used to identify language file; e.g. "en"
}

type LangDescriptionList []LangDescription
