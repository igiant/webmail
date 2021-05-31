package webmail

// FreeBusyInterval - FreeBusy status for particular interval.
type FreeBusyInterval struct {
	Status kerio::jsonapi::webmail::calendars::FreeBusyStatus `json:"status"` 
	Start UtcDateTime `json:"start"` 
	End UtcDateTime `json:"end"` 
}

// FreeBusySequence - FreeBusy sequence for particular interval.
type FreeBusySequence []FreeBusyInterval

// FreeBusyList - List of free busy sequences.
type FreeBusyList []FreeBusySequence

// FreeBusy management.

// FreeBusyGet - Free status is not being inserted into the result lists. Empty FreeBusySequence means the user is free for whole the interval.
func (c *ClientConnection) FreeBusyGet(userAddresses StringList, start UtcDateTime, end UtcDateTime) (FreeBusyList, error) {
	params := struct {
		UserAddresses StringList `json:"userAddresses"`
		Start UtcDateTime `json:"start"`
		End UtcDateTime `json:"end"`
	}{userAddresses, start, end}
	data, err := c.CallRaw("FreeBusy.get", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List FreeBusyList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}
