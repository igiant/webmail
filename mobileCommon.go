package webmail

// DeviceStatus - Mobile device status.
type DeviceStatus string

const (
	OK                   DeviceStatus = "OK"                   // no wipe
	DeviceNotProvisioned DeviceStatus = "DeviceNotProvisioned" // not (fully) provisioned (yet)
	DeviceWipeInitiated  DeviceStatus = "DeviceWipeInitiated"  // wipe process submitted
	DeviceWipeInProgress DeviceStatus = "DeviceWipeInProgress" // wipe process in progress
	DeviceWipeFinished   DeviceStatus = "DeviceWipeFinished"   // wipe process finished
	DeviceConnected      DeviceStatus = "DeviceConnected"
	DeviceDisconnected   DeviceStatus = "DeviceDisconnected"
)

// FolderIcon - Folder icon enumeration.
type FolderIcon string

const (
	FIMail     FolderIcon = "FIMail"
	FIContact  FolderIcon = "FIContact"
	FICalendar FolderIcon = "FICalendar"
	FITodo     FolderIcon = "FITodo"
	FIJournal  FolderIcon = "FIJournal"
	FINote     FolderIcon = "FINote"
	FIInbox    FolderIcon = "FIInbox"
	FIDeleted  FolderIcon = "FIDeleted"
)

// MobileSyncFolder - Synchronized folder.
type MobileSyncFolder struct {
	FolderName      string        `json:"folderName"`      // folder name
	FolderTypeIcon  FolderIcon    `json:"folderTypeIcon"`  // mail,contact...
	LastSyncDate    DateTimeStamp `json:"lastSyncDate"`    // date of last synchronization
	LastSyncDateIso UtcDateTime   `json:"lastSyncDateIso"` // date of last synchronization
}

type MobileSyncFolderList []MobileSyncFolder

// SyncMethod - Used synchronization method.
type SyncMethod string

const (
	ServerWins SyncMethod = "ServerWins"
	ClientWins SyncMethod = "ClientWins"
)

type ProtocolType string

const (
	protocolASync ProtocolType = "protocolASync"
	protocolKBC   ProtocolType = "protocolKBC"
)

// MobileDevice - Mobile device properties.
type MobileDevice struct {
	ProtocolType        ProtocolType         `json:"protocolType"`
	DeviceId            string               `json:"deviceId"`
	ProtocolVersion     string               `json:"protocolVersion"`     // used ActiveSync protocol version
	RegistrationDate    DateTimeStamp        `json:"registrationDate"`    // date of registration
	RegistrationDateIso UtcDateTime          `json:"registrationDateIso"` // date of registration
	LastSyncDate        DateTimeStamp        `json:"lastSyncDate"`        // date of last synchronization
	LastSyncDateIso     UtcDateTime          `json:"lastSyncDateIso"`     // date of last synchronization
	FolderList          MobileSyncFolderList `json:"folderList"`          // list of synchronized folders
	Status              DeviceStatus         `json:"status"`              // wipe status
	Method              SyncMethod           `json:"method"`              // synchronization method
	RemoteHost          string               `json:"remoteHost"`          // typically IP address of device
	Os                  string               `json:"os"`                  // operating system - eg. Windows Mobile(R) 2003
	Platform            string               `json:"platform"`            // PocketPC
	DeviceIcon          string               `json:"deviceIcon"`          // Device icon Eg. 'pocketpc' or 'unknown'
}

// MobileDeviceList - List of mobile devices.
type MobileDeviceList []MobileDevice
