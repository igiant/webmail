package webmail

// StringList - Type for lists of strings.
type StringList []string

// IntegerList - Type for lists of integers.
type IntegerList []int

// NamedValue - Note: all fields must be assigned if used in set methods
type NamedValue struct {
	Name  string `json:"name"` // Name-value pair
	Value string `json:"value"`
}

// NamedValueList - List of name-value pairs
type NamedValueList []NamedValue

// NamedMultiValue - Note: all fields must be assigned if used in set methods
type NamedMultiValue struct {
	Name  string     `json:"name"` // Name-multivalue pair
	Value StringList `json:"value"`
}

// NamedMultiValueList - List of name-multivalue pairs
type NamedMultiValueList []NamedMultiValue

// KId - global object identification
type KId string

// KIdList - list of global object identifiers
type KIdList []KId

// SortDirection - Sorting Direction
type SortDirection string

const (
	Asc  SortDirection = "Asc"  // ascending order
	Desc SortDirection = "Desc" // descending order
)

// CompareOperator - Simple Query Operator
type CompareOperator string

const (
	Eq          CompareOperator = "Eq"          // '=' - equal
	NotEq       CompareOperator = "NotEq"       // '<>' - not equal
	LessThan    CompareOperator = "LessThan"    // '<' - lower that
	GreaterThan CompareOperator = "GreaterThan" // '>' - greater that
	LessEq      CompareOperator = "LessEq"      // '<=' - lower or equal
	GreaterEq   CompareOperator = "GreaterEq"   // '>=' - greater or equal
	Like        CompareOperator = "Like"        // contains substring, % is wild character
)

// LogicalOperator - Compound Operator
type LogicalOperator string

const (
	Or  LogicalOperator = "Or"
	And LogicalOperator = "And"
)

// SubCondition - A Part of a Condition
type SubCondition struct {
	FieldName  string          `json:"fieldName"`  // left side of condition
	Comparator CompareOperator `json:"comparator"` // middle of condition
	Value      string          `json:"value"`      // right side of condition
}

// SubConditionList - A Complete Condition
type SubConditionList []SubCondition

// SortOrder - Sorting Order
type SortOrder struct {
	ColumnName    string        `json:"columnName"`
	Direction     SortDirection `json:"direction"`
	CaseSensitive bool          `json:"caseSensitive"`
}

// SortOrderList - List of Sorting Orders
type SortOrderList []SortOrder

// SearchQuery - General Query for Searching
// Query substitution (quicksearch):
// SearchQuery doesn't support complex queries, only queries
// with all AND operators (or all OR operators) are supported.
// Combination of AND and OR is not allowed. This limitation is for special cases solved by using
// substitution of complicated query-part by simple condition.
// Only the quicksearch is currently implemented and only in "Users::get()" method.
// Behavior of quicksearch in Users::get():
// QUICKSEACH  = "x"  is equal to:  (loginName  = "x") OR (fullName  = "x")
// QUICKSEACH LIKE "x*" is equal to:  (loginName LIKE "x*") OR (fullName LIKE "x*")
// SearchQuery - QUICKSEACH  <> "x"  is equal to:  (loginName  <> "x") AND (fullName  <> "x")
type SearchQuery struct {
	Fields     StringList       `json:"fields"`     // empty = give me all fields, applicable constants: ADD_USERS, LIST_USERS
	Conditions SubConditionList `json:"conditions"` // empty = without condition
	Combining  LogicalOperator  `json:"combining"`  // the list of conditions can be either combined by 'ORs' or 'ANDs'
	Start      int              `json:"start"`      // how many items to skip before filling a result list (0 means skip none)
	Limit      int              `json:"limit"`      // how many items to put to a result list (if there are enough items); applicable constant: Unlimited
	OrderBy    SortOrderList    `json:"orderBy"`
}

// LocalizableMessage - Message can contain replacement marks: { "User %1 cannot be deleted.", ["jsmith"], 1 }
type LocalizableMessage struct {
	Message              string     `json:"message"`              // text with placeholders %1, %2, etc., e.g. "User %1 cannot be deleted."
	PositionalParameters StringList `json:"positionalParameters"` // additional strings to replace the placeholders in message (first string replaces %1 etc.)
	Plurality            int        `json:"plurality"`            // count of items, used to distinguish among singular/paucal/plural; 1 for messages with no counted items
}

type LocalizableMessageList []LocalizableMessage

// ManipulationError - error structure to be used when manipulating with globally addressable list items
type ManipulationError struct {
	Id           KId                `json:"id"` // entity KId, can be user, group, alias, ML...
	ErrorMessage LocalizableMessage `json:"errorMessage"`
}

type ManipulationErrorList []ManipulationError

// RestrictionKind - A kind of restriction
type RestrictionKind string

const (
	Regex                  RestrictionKind = "Regex"                  // regular expression
	ByteLength             RestrictionKind = "ByteLength"             // maximal length in Bytes
	ForbiddenNameList      RestrictionKind = "ForbiddenNameList"      // list of denied exact names due to filesystem or KMS store
	ForbiddenPrefixList    RestrictionKind = "ForbiddenPrefixList"    // list of denied preffixes due to filesystem or KMS store
	ForbiddenSuffixList    RestrictionKind = "ForbiddenSuffixList"    // list of denied suffixes due to filesystem or KMS store
	ForbiddenCharacterList RestrictionKind = "ForbiddenCharacterList" // list of denied characters
)

// ItemName - Item of the Entity; used in restrictions
type ItemName string

const (
	Name        ItemName = "Name"        // Entity Name
	Description ItemName = "Description" // Entity Description
	Email       ItemName = "Email"       // Entity Email Address
	FullName    ItemName = "FullName"    // Entity Full Name
	TimeItem    ItemName = "TimeItem"    // Entity Time - it cannot be simply Time because of C++ conflict - see bug 34684 comment #3
	DateItem    ItemName = "DateItem"    // Entity Date - I expect same problem with Date as with Time
	DomainName  ItemName = "DomainName"  // differs from name (eg. cannot contains underscore)
)

// Units used for handling large values of bytes.

// ByteUnits - See also userinfo.idl: enum UserValueUnits.
type ByteUnits string

const (
	Bytes     ByteUnits = "Bytes"
	KiloBytes ByteUnits = "KiloBytes"
	MegaBytes ByteUnits = "MegaBytes"
	GigaBytes ByteUnits = "GigaBytes"
	TeraBytes ByteUnits = "TeraBytes"
	PetaBytes ByteUnits = "PetaBytes"
)

// Stores size of very large values of bytes e.g. for user quota

// ByteValueWithUnits - Note: all fields must be assigned if used in set methods
type ByteValueWithUnits struct {
	Value int       `json:"value"`
	Units ByteUnits `json:"units"`
}

// Settings of size limit

// SizeLimit - Note: all fields must be assigned if used in set methods
type SizeLimit struct {
	IsActive bool               `json:"isActive"`
	Limit    ByteValueWithUnits `json:"limit"`
}

// AddResult - Result of the add operation
type AddResult struct {
	Id           KId                `json:"id"`           // purposely not id - loginName is shown
	Success      bool               `json:"success"`      // was operation successful? if yes so id is new id for this item else errorMessage tells why it failed
	ErrorMessage LocalizableMessage `json:"errorMessage"` // contains number of recovered user messages or error message
}

// AddResultList - list of add operation results
type AddResultList []AddResult

type IpAddress string

type IpAddressList []IpAddress

// StoreStatus - Status of entry in persistent manager
type StoreStatus string

const (
	StoreStatusClean    StoreStatus = "StoreStatusClean"    // already present in configuration store
	StoreStatusModified StoreStatus = "StoreStatusModified" // update waiting for apply()
	StoreStatusNew      StoreStatus = "StoreStatusNew"      // added to manager but not synced to configuration store
)

// Time
// When using start and limit to only get a part of all results
// (e.g. only 20 users, skipping the first 40 users),
// use this special limit value for unlimited count
// (of course the service still respects the value of start).
// Note that each service is allowed to use its safety limit
// (such as 50,000) to prevent useless overload.
// The limits are documented per-service or per-method.
// Implementation note: Some source code transformations may lead to signed long, i.e. 4294967295.
// But the correct value is -1.
// Date and Time - should be used instead of time_t, where time zones can affect time interpretation
// Time - Note: all fields must be assigned if used in set methods
type Time struct {
	Hour int `json:"hour"` // 0-23
	Min  int `json:"min"`  // 0-59
}

// Date - Note: all fields must be assigned if used in set methods
type Date struct {
	Year  int `json:"year"`
	Month int `json:"month"` // 0-11
	Day   int `json:"day"`   // 1-31 max day is limited by month
}

// A string that can be switched on/off. String is meaningful only if switched on.

// OptionalString - Note: all fields must be assigned if used in set methods
type OptionalString struct {
	Enabled bool   `json:"enabled"`
	Value   string `json:"value"`
}

// OptionalLong - Note: all fields must be assigned if used in set methods
type OptionalLong struct {
	Enabled bool `json:"enabled"`
	Value   int  `json:"value"`
}

// IP Address Group / Time Range / ... that can be switched on/off

// OptionalEntity - Note: all fields must be assigned if used in set methods
type OptionalEntity struct {
	Enabled bool   `json:"enabled"`
	Id      KId    `json:"id"` // global identifier
	Name    string `json:"name"`
}

// Message can contain replacement marks: { "User %1 cannot be deleted.", ["jsmith"], 1 }.

// LocalizableMessageParameters - This is the parameters structure.
type LocalizableMessageParameters struct {
	PositionalParameters StringList `json:"positionalParameters"` // additional strings to replace the placeholders in message (first string replaces %1 etc.)
	Plurality            int        `json:"plurality"`            // count of items, used to distinguish among singular/paucal/plural; 1 for messages with no counted items
}

// Error - Error details regarding a particular item, e.g. one of users that could not be updated or removed.
type Error struct {
	InputIndex        int                          `json:"inputIndex"`        // 0-based index to input array, e.g. 3 means that the relates to the 4th element of the input parameter array
	Code              int                          `json:"code"`              // -32767..-1 (JSON-RPC) or 1..32767 (application)
	Message           string                       `json:"message"`           // text with placeholders %1, %2, etc., e.g. "User %1 cannot be deleted."
	MessageParameters LocalizableMessageParameters `json:"messageParameters"` // strings to replace placeholders in message, and message plurality.
}

type ErrorList []Error

// Download - important information about download
type Download struct {
	Url    string `json:"url"`    // download url
	Name   string `json:"name"`   // filename
	Length int    `json:"length"` // file size in bytes
}

// DateTimeStamp - Type for date/time representation
type DateTimeStamp int

// ApiApplication - Describes client (third-party) application or script which uses the Administration API.
type ApiApplication struct {
	Name    string `json:"name"`    // E.g. "Simple server monitor"
	Vendor  string `json:"vendor"`  // E.g. "MyScript Ltd."
	Version string `json:"version"` // E.g. "1.0.0 beta 1"
}

// Credentials - Credentials contains userName and password
type Credentials struct {
	UserName string `json:"userName"` // UserName
	Password string `json:"password"` // Password
}
