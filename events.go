package webmail

import "encoding/json"

type Event struct {
	Id            KId            `json:"id"`       // [READ-ONLY] global identification
	FolderId      KId            `json:"folderId"` // [REQUIRED FOR CREATE] [WRITE-ONCE] global identification of folder in which is the event defined
	Watermark     Watermark      `json:"watermark"`
	Access        EventAccess    `json:"access"` // [READ-ONLY] scope of access of user to this event
	Summary       string         `json:"summary"`
	Location      string         `json:"location"`
	Description   string         `json:"description"`
	Label         EventLabel     `json:"label"`
	Categories    StringList     `json:"categories"`
	Start         UtcDateTime    `json:"start"`
	End           UtcDateTime    `json:"end"`
	TravelMinutes int            `json:"travelMinutes"` // // X-APPLE-TRAVEL-DURATION;VALUE=DURATION:PT15M
	FreeBusy      FreeBusyStatus `json:"freeBusy"`      // also known as TimeTransparency
	IsPrivate     bool           `json:"isPrivate"`     // also known as Class
	IsAllDay      bool           `json:"isAllDay"`
	Priority      PriorityType   `json:"priority"`
	Rule          RecurrenceRule `json:"rule"`
	Attendees     AttendeeList   `json:"attendees"`
	Reminder      Reminder       `json:"reminder"`
	IsCancelled   bool           `json:"isCancelled"` // [READ-ONLY] is cancelled by organiser
}

type EventList []Event

type EventUpdateType string

const (
	EUpdateRequest EventUpdateType = "EUpdateRequest"
	EUpdateReply   EventUpdateType = "EUpdateReply"
	EUpdateCancel  EventUpdateType = "EUpdateCancel"
)

type EventActionType string

const (
	EActionCreate             EventActionType = "EActionCreate"             // new invitation
	EActionChangedTime        EventActionType = "EActionChangedTime"        // time of meating was changed
	EActionChangedSummary     EventActionType = "EActionChangedSummary"     // summary of meating was changed
	EActionChangedLocation    EventActionType = "EActionChangedLocation"    // location of meating was changed
	EActionChangedDescription EventActionType = "EActionChangedDescription" // description of meating was changed
)

type EventActionTypeList []EventActionType

type EventUpdate struct {
	Id            KId                 `json:"id"`            // [READ-ONLY] global identification (e-mail where an update or an invitation is placed)
	EventId       KId                 `json:"eventId"`       // [READ-ONLY] global identification of caused event
	EventFolderId KId                 `json:"eventFolderId"` // [READ-ONLY] global identification of caused event
	OccurrenceId  KId                 `json:"occurrenceId"`  // [READ-ONLY] global identification of caused event (if whole recurrent event is updated , there is first occurrence)
	IsException   bool                `json:"isException"`
	SeqNumber     int                 `json:"seqNumber"`
	IsObsolete    bool                `json:"isObsolete"`
	DeliveryTime  UtcDateTime         `json:"deliveryTime"`
	Type          EventUpdateType     `json:"type"`
	Summary       string              `json:"summary"`
	Location      string              `json:"location"`
	Start         UtcDateTime         `json:"start"`
	End           UtcDateTime         `json:"end"`
	TotalEnd      UtcDateTime         `json:"totalEnd"`
	Description   string              `json:"description"`
	Attendee      Attendee            `json:"attendee"`
	Actions       EventActionTypeList `json:"actions"`
}

type EventUpdateList []EventUpdate

// Constants for composing kerio::web::SearchQuery

// EventsGet - Get a list of events.
// Parameters
//	query - query attributes and limits
// Return
//	list - all found events
//  totalItems - number of events found if there is no limit
func (c *ClientConnection) EventsGet(ids KIdList, query SearchQuery) (EventList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		Ids   KIdList     `json:"ids"`
		Query SearchQuery `json:"query"`
	}{ids, query}
	data, err := c.CallRaw("Events.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       EventList `json:"list"`
			TotalItems int       `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// EventsGetById - Get an event.
// Parameters
//	id - global identifier of requested event
// Return
//	result - found event
func (c *ClientConnection) EventsGetById(id KId) (*Event, error) {
	params := struct {
		Id KId `json:"id"`
	}{id}
	data, err := c.CallRaw("Events.getById", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Result struct {
			Result Event `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	return &result.Result.Result, err
}

// EventsGetEventUpdates - Get updates or invitations from Calendar INBOX by global identifiers.
// Parameters
//	ids - list of global identifiers of EventUpdates
// Return
//	errors - list of updates that failed to optain
//	eventUpdates - list of updates or invitattions
func (c *ClientConnection) EventsGetEventUpdates(ids KIdList) (ErrorList, EventUpdateList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Events.getEventUpdates", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors       ErrorList       `json:"errors"`
			EventUpdates EventUpdateList `json:"eventUpdates"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.EventUpdates, err
}

// EventsGetEventUpdateList - Get all updates or invitations from Calendar INBOX.
// Return
//	eventUpdates - list of updates or invitattions
func (c *ClientConnection) EventsGetEventUpdateList() (EventUpdateList, error) {
	data, err := c.CallRaw("Events.getEventUpdateList", nil)
	if err != nil {
		return nil, err
	}
	eventUpdates := struct {
		Result struct {
			EventUpdates EventUpdateList `json:"eventUpdates"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &eventUpdates)
	return eventUpdates.Result.EventUpdates, err
}

// EventsGetSharedEventUpdateList - Get all updates or invitations from Calendar INBOX.
// Parameters
//	mailboxIds - list of global identifiers of mailboxes
// Return
//	errors - list of mailboxes that failed to search
//	eventUpdates - list of updates or invitattions
func (c *ClientConnection) EventsGetSharedEventUpdateList(mailboxIds KIdList) (ErrorList, EventUpdateList, error) {
	params := struct {
		MailboxIds KIdList `json:"mailboxIds"`
	}{mailboxIds}
	data, err := c.CallRaw("Events.getSharedEventUpdateList", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors       ErrorList       `json:"errors"`
			EventUpdates EventUpdateList `json:"eventUpdates"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.EventUpdates, err
}

// EventsRemove - Remove a list of events.
// Parameters
//	ids - list of global identifiers of events to be removed
// Return
//	errors - list of events that failed to remove
func (c *ClientConnection) EventsRemove(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Events.remove", params)
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

// EventsRemoveEventUpdates - Remove a list of EventUpdates.
// Parameters
//	ids - list of global identifiers of EventUpdates to be removed
// Return
//	errors - list of updates that failed to remove
func (c *ClientConnection) EventsRemoveEventUpdates(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Events.removeEventUpdates", params)
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

// EventsCopy - Copy existing events to folder
// Parameters
//	ids - list of global identifiers of events to be copied
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) EventsCopy(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Events.copy", params)
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

// EventsCreate - Create events.
// Parameters
//	events - list of events to be created
// Return
//	errors - list of events that failed on creation
//	result - particular results for all items
func (c *ClientConnection) EventsCreate(events EventList) (ErrorList, CreateResultList, error) {
	params := struct {
		Events EventList `json:"events"`
	}{events}
	data, err := c.CallRaw("Events.create", params)
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

// EventsCreateFromAttachment - Get an occurrence.
// Parameters
//	attachmentId - global identifier of attachment
// Return
//	result - result
func (c *ClientConnection) EventsCreateFromAttachment(attachmentId KId) (*CreateResult, error) {
	params := struct {
		AttachmentId KId `json:"attachmentId"`
	}{attachmentId}
	data, err := c.CallRaw("Events.createFromAttachment", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Result struct {
			Result CreateResult `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	return &result.Result.Result, err
}

// EventsSet - Set events.
// Parameters
//	events - modifications of events.
// Return
//	errors - error message list
func (c *ClientConnection) EventsSet(events EventList) (ErrorList, SetResultList, error) {
	params := struct {
		Events EventList `json:"events"`
	}{events}
	data, err := c.CallRaw("Events.set", params)
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

// EventsMove - Move existing events to folder
// Parameters
//	ids - list of global identifiers of events to be moved
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) EventsMove(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Events.move", params)
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
