package webmail

// Filtering rule that has one or more initial conditions and
// one or more actions that are performed only if the
// FilterRule - filter's initial conditions are meet.
type FilterRule struct {
	Id             KId                 `json:"id"`             // [READ-ONLY] global identification
	IsEnabled      bool                `json:"isEnabled"`      // says whether rule is enabled
	Description    string              `json:"description"`    // contains rules description
	IsIncomplete   bool                `json:"isIncomplete"`   // if rule is not completed (it does not contain any definition of conditions and actions)
	Conditions     FilterConditionList `json:"conditions"`     // list of rule's initial conditions
	Actions        FilterActionList    `json:"actions"`        // list of rule's actions (performed if initial conditions are meet)
	EvaluationMode EvaluationModeType  `json:"evaluationMode"` // determines evaluation mod of initial conditions
}

type FilterRuleList []FilterRule

type FilterRawRule struct {
	Id          KId    `json:"id"`          // [READ-ONLY] global identification
	IsEnabled   bool   `json:"isEnabled"`   // says whether rule is enabled
	Description string `json:"description"` // contains rules description
	Script      string `json:"script"`
}

// FiltersGet - by user
// Return
//	dataStamp - server as concurrent modification protection
//	filters - list of all messages filtering rules defined
func (c *ClientConnection) FiltersGet() (*unsignedlong, FilterRuleList, error) {
	data, err := c.CallRaw("Filters.get", nil)
	if err != nil {
		return nil, nil, err
	}
	dataStamp := struct {
		Result struct {
			DataStamp unsignedlong   `json:"dataStamp"`
			Filters   FilterRuleList `json:"filters"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &dataStamp)
	return &dataStamp.Result.DataStamp, dataStamp.Result.Filters, err
}

// FiltersGetById - Obtain particular rule in a script form.
// Parameters
//	currentDataStamp - the stamp obtained via function get
//	id - ID of rule
// Return
//	rule - the script
func (c *ClientConnection) FiltersGetById(currentDataStamp unsignedlong, id KId) (*FilterRawRule, error) {
	params := struct {
		CurrentDataStamp unsignedlong `json:"currentDataStamp"`
		Id               KId          `json:"id"`
	}{currentDataStamp, id}
	data, err := c.CallRaw("Filters.getById", params)
	if err != nil {
		return nil, err
	}
	rule := struct {
		Result struct {
			Rule FilterRawRule `json:"rule"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &rule)
	return &rule.Result.Rule, err
}

// FiltersGenerateRule - Obtain rule in a script form generated from pattern.
// Parameters
//	pattern - structured rule
// Return
//	rule - the script
func (c *ClientConnection) FiltersGenerateRule(pattern FilterRule) (*FilterRawRule, error) {
	params := struct {
		Pattern FilterRule `json:"pattern"`
	}{pattern}
	data, err := c.CallRaw("Filters.generateRule", params)
	if err != nil {
		return nil, err
	}
	rule := struct {
		Result struct {
			Rule FilterRawRule `json:"rule"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &rule)
	return &rule.Result.Rule, err
}

// FiltersSet - by user
// Parameters
//	currentDataStamp - the stamp which was assigned to you by
//	filters - list of all messages filtering rules defined
// Return
//	newDataStamp - a new stamp that replaces your current
func (c *ClientConnection) FiltersSet(currentDataStamp unsignedlong, filters FilterRuleList) (*unsignedlong, error) {
	params := struct {
		CurrentDataStamp unsignedlong   `json:"currentDataStamp"`
		Filters          FilterRuleList `json:"filters"`
	}{currentDataStamp, filters}
	data, err := c.CallRaw("Filters.set", params)
	if err != nil {
		return nil, err
	}
	newDataStamp := struct {
		Result struct {
			NewDataStamp unsignedlong `json:"newDataStamp"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &newDataStamp)
	return &newDataStamp.Result.NewDataStamp, err
}

// FiltersSetById - Set particular rule.
// Parameters
//	currentDataStamp - the stamp obtained via function get
//	rule - the new script
// Return
//	newDataStamp - a new stamp
func (c *ClientConnection) FiltersSetById(currentDataStamp unsignedlong, rule FilterRawRule) (*unsignedlong, error) {
	params := struct {
		CurrentDataStamp unsignedlong  `json:"currentDataStamp"`
		Rule             FilterRawRule `json:"rule"`
	}{currentDataStamp, rule}
	data, err := c.CallRaw("Filters.setById", params)
	if err != nil {
		return nil, err
	}
	newDataStamp := struct {
		Result struct {
			NewDataStamp unsignedlong `json:"newDataStamp"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &newDataStamp)
	return &newDataStamp.Result.NewDataStamp, err
}
