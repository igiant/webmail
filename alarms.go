package webmail

type Alarm struct {
	Type         ItemType    `json:"type"`         // [READ-ONLY] only 'Calendar' and 'Task' are valid
	ItemId       KId         `json:"itemId"`       // [READ-ONLY] global identification of occurrence
	BaseId       KId         `json:"baseId"`       // [READ-ONLY] global identification of event or task
	Summary      string      `json:"summary"`      // [READ-ONLY]
	Location     string      `json:"location"`     // [READ-ONLY]
	Start        UtcDateTime `json:"start"`        // [READ-ONLY] can be empty in case of item 'Task'
	End          UtcDateTime `json:"end"`          // [READ-ONLY] can be empty in case of item 'Task'. In case of an all day event is there begin of next day.
	Due          UtcDateTime `json:"due"`          // [READ-ONLY] can be empty and is valid only for item 'Task'.
	IsAllDay     bool        `json:"isAllDay"`     // [READ-ONLY]
	ReminderTime UtcDateTime `json:"reminderTime"` // time when remainder should appear
}

type AlarmList []Alarm

// AlarmsDismiss - the method then a reminder will be removed as well as this value.
// Parameters
//	itemIds - list of event or occurrence IDs
// Return
//	errors - list of errors
func (c *ClientConnection) AlarmsDismiss(itemIds KIdList) (ErrorList, error) {
	params := struct {
		ItemIds KIdList `json:"itemIds"`
	}{itemIds}
	data, err := c.CallRaw("Alarms.dismiss", params)
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

// AlarmsGet - Get alarms. Alarms are searched in range from now to value 'until'.
// Parameters
//	since - lower bound of time range
//	until - upper bound of time range
// Return
//	list - list of alarms
func (c *ClientConnection) AlarmsGet(since UtcTime, until UtcTime) (AlarmList, error) {
	params := struct {
		Since UtcTime `json:"since"`
		Until UtcTime `json:"until"`
	}{since, until}
	data, err := c.CallRaw("Alarms.get", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List AlarmList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// AlarmsSet - Value is placed in X-NEXT-ALARM property.
// Parameters
//	nextTime - time
//	itemIds - list of event or occurrence IDs
// Return
//	errors - list of errors
func (c *ClientConnection) AlarmsSet(nextTime UtcTime, itemIds KIdList) (ErrorList, error) {
	params := struct {
		NextTime UtcTime `json:"nextTime"`
		ItemIds  KIdList `json:"itemIds"`
	}{nextTime, itemIds}
	data, err := c.CallRaw("Alarms.set", params)
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
