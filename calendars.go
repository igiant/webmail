package webmail

type LongList []int

type FreeBusyStatus string

const (
	Busy         FreeBusyStatus = "Busy" // opaque
	Tentative    FreeBusyStatus = "Tentative"
	Free         FreeBusyStatus = "Free" // transparent
	OutOfOffice  FreeBusyStatus = "OutOfOffice"
	NotAvailable FreeBusyStatus = "NotAvailable"
)

// EventLabel - Any event can be marked by a label
type EventLabel string

const (
	None             EventLabel = "None"
	Important        EventLabel = "Important"
	Business         EventLabel = "Business"
	Personal         EventLabel = "Personal"
	Vacation         EventLabel = "Vacation"
	MustAttend       EventLabel = "MustAttend"
	TravelRequired   EventLabel = "TravelRequired"
	NeedsPreparation EventLabel = "NeedsPreparation"
	BirthDay         EventLabel = "BirthDay"
	Anniversary      EventLabel = "Anniversary"
	PhoneCall        EventLabel = "PhoneCall"
)

// AttendeeRole - Purposly merged role and type to avoid useless combinations like optional room
type AttendeeRole string

const (
	RoleOrganizer        AttendeeRole = "RoleOrganizer"
	RoleRequiredAttendee AttendeeRole = "RoleRequiredAttendee"
	RoleOptionalAttendee AttendeeRole = "RoleOptionalAttendee"
	RoleRoom             AttendeeRole = "RoleRoom"
	RoleEquipment        AttendeeRole = "RoleEquipment"
)

// Enumerate of attendee participation status. It specifies the participation status for the calendar user
// PartStatus - specified by the property. If not specified, the default value is PartNotResponded.
type PartStatus string

const (
	PartNotResponded PartStatus = "PartNotResponded" // Event needs action
	PartAccepted     PartStatus = "PartAccepted"     // Event accepted
	PartDeclined     PartStatus = "PartDeclined"     // Event declined
	PartDelegated    PartStatus = "PartDelegated"    // Event delegated
	PartTentative    PartStatus = "PartTentative"    // Event tentatively accepted
)

// PartStatusResponse - Response to particular occurrence or event.
type PartStatusResponse struct {
	Status  PartStatus `json:"status"`
	Message string     `json:"message"`
}

type Attendee struct {
	DisplayName  string       `json:"displayName"`
	EmailAddress string       `json:"emailAddress"` // [REQUIRED]
	Role         AttendeeRole `json:"role"`         // [REQUIRED]
	IsNotified   bool         `json:"isNotified"`   // [READ-ONLY] is Attendee notified by email on event update? also known as RSVP
	PartStatus   PartStatus   `json:"partStatus"`   // [READ-ONLY] A participation status for the event*/
}

type AttendeeList []Attendee

type ReminderType string

const (
	ReminderRelative ReminderType = "ReminderRelative"
	ReminderAbsolute ReminderType = "ReminderAbsolute"
)

type Reminder struct {
	IsSet              bool         `json:"isSet"`
	Type               ReminderType `json:"type"` // if it is not send the default value is 'ReminderRelative'
	MinutesBeforeStart int          `json:"minutesBeforeStart"`
	Date               UtcDateTime  `json:"date"`
}

type FrequencyType string

const (
	Daily   FrequencyType = "Daily"
	Weekly  FrequencyType = "Weekly"
	Monthly FrequencyType = "Monthly"
	Yearly  FrequencyType = "Yearly"
)

type EndByType string

const (
	ByRecurrenceNever EndByType = "ByRecurrenceNever"
	ByRecurrenceDate  EndByType = "ByRecurrenceDate"
)

type EndBy struct {
	Type EndByType   `json:"type"`
	Date UtcDateTime `json:"date"` // also known as until, used for ByRecurrenceDate
}

type PreciseBy struct {
	ByDay      LongList `json:"byDay"`
	ByMonthDay LongList `json:"byMonthDay"`
	ByMonth    LongList `json:"byMonth"`
	ByPosition LongList `json:"byPosition"` // 2 = 2nd day, 3 = 3rd day
	ByInterval int      `json:"byInterval"` // 2 = every 2nd day, 3 = every 3rd day
}

type RecurrenceRule struct {
	IsSet     bool          `json:"isSet"`
	Frequency FrequencyType `json:"frequency"` // period in which event occurs
	EndBy     EndBy         `json:"endBy"`     // end limitation
	PreciseBy PreciseBy     `json:"preciseBy"` // further specification
}

// EventAccess - Purposly merged role and type to avoid useless combinations like optional room
type EventAccess string

const (
	EAccessCreator  EventAccess = "EAccessCreator"
	EAccessInvitee  EventAccess = "EAccessInvitee"
	EAccessReadOnly EventAccess = "EAccessReadOnly"
)
