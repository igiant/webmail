package webmail

type ModificationType string
const (
	modifyAll ModificationType = "modifyAll" 
	modifyThis ModificationType = "modifyThis" 
	modifyAllFollowing ModificationType = "modifyAllFollowing" 
	modifyMasterEvent ModificationType = "modifyMasterEvent" 
)

type Occurrence struct {
	Id KId `json:"id"` // [READ-ONLY] global identification
	EventId KId `json:"eventId"` // [READ-ONLY] global identification of appropriate event
	FolderId KId `json:"folderId"` // [REQUIRED FOR CREATE] [WRITE-ONCE] global identification of folder in which is the event defined
	Watermark Watermark `json:"watermark"` 
	Access kerio::jsonapi::webmail::calendars::EventAccess `json:"access"` // [READ-ONLY] scope of access of user to this occurrence
	Summary string `json:"summary"` 
	Location string `json:"location"` 
	Description string `json:"description"` 
	Label kerio::jsonapi::webmail::calendars::EventLabel `json:"label"` 
	Categories StringList `json:"categories"` 
	Start UtcDateTime `json:"start"` 
	End UtcDateTime `json:"end"` 
	TravelMinutes int `json:"travelMinutes"` // // X-APPLE-TRAVEL-DURATION;VALUE=DURATION:PT15M
	FreeBusy kerio::jsonapi::webmail::calendars::FreeBusyStatus `json:"freeBusy"` // also known as TimeTransparency
	IsPrivate bool `json:"isPrivate"` // also known as Class
	IsAllDay bool `json:"isAllDay"` 
	Priority PriorityType `json:"priority"` 
	Rule kerio::jsonapi::webmail::calendars::RecurrenceRule `json:"rule"` // not filled in listing method
	Attendees kerio::jsonapi::webmail::calendars::AttendeeList `json:"attendees"` 
	Reminder kerio::jsonapi::webmail::calendars::Reminder `json:"reminder"` // not filled in listing method
	IsException bool `json:"isException"` // [READ-ONLY] it does not make sense to write it
	HasReminder bool `json:"hasReminder"` // [READ-ONLY]
	IsRecurrent bool `json:"isRecurrent"` // [READ-ONLY]
	IsCancelled bool `json:"isCancelled"` // [READ-ONLY] is cancelled by organiser
	SeqNumber int `json:"seqNumber"` // [READ-ONLY]
	Modification ModificationType `json:"modification"` // [WRITE-ONLY]
}

type OccurrenceList []Occurrence

// Constants for composing kerio::web::SearchQuery 

// OccurrencesGet - Items rule and reminder in the occurrence aren't filled. If necessary use getOccurrence method.
// Parameters
//	folderIds - list of global identifiers of folders to be listed
//	query - query attributes and limits
// Return
//	list - all found events
//	totalItems - number of events found if there is no limit
func (c *ClientConnection) OccurrencesGet(folderIds KIdList, query SearchQuery) (OccurrenceList, int, error) {
	params := struct {
		FolderIds KIdList `json:"folderIds"`
		Query SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Occurrences.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List OccurrenceList `json:"list"`
			TotalItems int `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}


// OccurrencesGetById - Get an occurrence.
// Return
//	result - found occurrence
func (c *ClientConnection) OccurrencesGetById(ids KIdList) (ErrorList, OccurrenceList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Occurrences.getById", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result OccurrenceList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}


// OccurrencesGetFromAttachment - Get an occurrence.
// Parameters
//	attachmentId - global identifier of attachment
// Return
//	result - found occurrence
func (c *ClientConnection) OccurrencesGetFromAttachment(attachmentId KId) (*Occurrence, error) {
	params := struct {
		AttachmentId KId `json:"attachmentId"`
	}{attachmentId}
	data, err := c.CallRaw("Occurrences.getFromAttachment", params)
	if err != nil {
		return nil, err
	}
	result := struct {
		Result struct {
			Result Occurrence `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	return &result.Result.Result, err
}


// OccurrencesRemove - Remove a list of occurrences.
// Parameters
//	occurrences - occurrences to be removed. Only fields 'id' and 'modification' are required.
// Return
//	errors - list of occurrences that failed to remove
func (c *ClientConnection) OccurrencesRemove(occurrences OccurrenceList) (ErrorList, error) {
	params := struct {
		Occurrences OccurrenceList `json:"occurrences"`
	}{occurrences}
	data, err := c.CallRaw("Occurrences.remove", params)
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


// OccurrencesSet - Set occurrences.
// Parameters
//	occurrences - modifications of occurrences.
// Return
//	errors - error message list
func (c *ClientConnection) OccurrencesSet(occurrences OccurrenceList) (ErrorList, SetResultList, error) {
	params := struct {
		Occurrences OccurrenceList `json:"occurrences"`
	}{occurrences}
	data, err := c.CallRaw("Occurrences.set", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result SetResultList `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}


// OccurrencesSetPartStatus - Set part status to occurrence or event and send response to organizer.
// Parameters
//	id - identifiers of events or occurrence
//	response - response and status
func (c *ClientConnection) OccurrencesSetPartStatus(id KId, response kerio::jsonapi::webmail::calendars::PartStatusResponse) error {
	params := struct {
		Id KId `json:"id"`
		Response kerio::jsonapi::webmail::calendars::PartStatusResponse `json:"response"`
	}{id, response}
	_, err := c.CallRaw("Occurrences.setPartStatus", params)
	return err
}