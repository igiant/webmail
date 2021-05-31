package webmail

const (
	ErrorCodeParseError     = -32700 // Parse error. Request wasn't valid JSON. (HTTP Status 500)
	ErrorCodeInternalError  = -32603 // Internal error. (HTTP Status 500
	ErrorCodeInvalidParams  = -32602 // Invalid params. (HTTP Status 500)
	ErrorCodeMethodNotFound = -32601 // Method not found.The requested method doesn't exist. (HTTP Status 404)
	ErrorCodeInvalidRequest = -32600 // Invalid request. (Too many files uploaded.etc.) (HTTP Status 500)

	// -32099 to -32000 (Reserved for Kerio-defined server errors.)
	ErrorCodeMultiServerBackendMaintenance = -32003 ///< Error generated on client side only
	ErrorCodeTimedout                      = -32002 ///< Error generated on client side only
	ErrorCodeSessionExpired                = -32001 ///< Session expired. (HTTP Status 500)

	ErrorCodeCommunicationFailure = -1 ///< Error generated on client side only

	// 1 to 999 (Reserved)
	ErrorCodeRequestEntityTooLarge = 413 ///< The client tried to upload a file bigger than limit. (Generally any HTTP POST request, including JSON-RPC Request.)

	// 1000 to 1999 (Reserved for Kerio-common errors)
	ErrorCodeOperationFailed = 1000 ///< The command was accepted, the operation was run and returned an error.
	ErrorCodeAlreadyExists   = 1001 ///< Can't create the item, as it already exists.
	ErrorCodeNoSuchEntity    = 1002 ///< Message / folder / etc. doesn't exist.
	ErrorCodeNotPermitted    = 1003 ///< Server refused to proceed. E.g. Can't delete default folder.
	ErrorCodeAccessDenied    = 1004 ///< Insufficient privileges for the required operation.
	// 4000 to 4999 (Reserved for Connect)
	ErrorCodeDangerousOperation = 4000 ///< Operation is dangerous. Used for reporting in validation functions.
	ErrorCodePartialSuccess     = 4001 ///< Operation ended partial successful.
	ErrorCodeChangePswFailed    = 4002 ///< Failed to change your password. The new password was not accepted.

	// 4100 to 4199 (Reserved for Connect/Webmail specific)
	ErrorCodeFolderReindexing       = 4100 ///< Folder is already being reindexed. Try later.
	ErrorCodeOperationInProgress    = 4101 ///< Long lasting operation blocks perform the request. Try later.
	ErrorCodeQuotaReached           = 4102 ///< Items quota or disk size quota was reached.
	ErrorCodeSendingFailed          = 4103 ///< Failed to send email.
	ErrorCodeNoSuchFolder           = 4104 ///< Folder of given GUID doesn't exist.
	ErrorCodeOperatorSessionExpired = 4105 ///< Session to Operator expired.
)

type UtcTime string

type Watermark int64

type PriorityType string

type jsonstring string

const (
	Normal PriorityType = "Normal" // default value
	Low    PriorityType = "Low"
	High   PriorityType = "High"
)

// CreateResult - Details about a particular item created.
type CreateResult struct {
	InputIndex int       `json:"inputIndex"` // 0-based index to input array, e.g. 3 means that the relates to the 4th element of the input parameter array
	Id         KId       `json:"id"`         // ID of created item.
	Watermark  Watermark `json:"watermark"`  // item version from journal
}

type CreateResultList []CreateResult

// SetResult - Details about a particular item updated.
type SetResult struct {
	InputIndex int       `json:"inputIndex"` // 0-based index to input array, e.g. 3 means that the relates to the 4th element of the input parameter array
	Id         KId       `json:"id"`         // if not empty the item was moved and its ID was changed to this value
	Watermark  Watermark `json:"watermark"`  // item version from journal
}

type SetResultList []SetResult

type Image struct {
	Url string `json:"url"` // [READ ONLY] URL to obtain image via HTTP GET request.
	Id  string `json:"id"`  // [READ ONLY] Id of uploaded image.
}

type ImageList []Image
