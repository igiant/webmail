package webmail

import "encoding/json"

type TaskStatus string

const (
	tsNotStarted TaskStatus = "tsNotStarted"
	tsCompleted  TaskStatus = "tsCompleted"
	tsInProgress TaskStatus = "tsInProgress"
	tsWaiting    TaskStatus = "tsWaiting"
	tsDeferred   TaskStatus = "tsDeferred"
)

// Task - Task details
type Task struct {
	Id          KId            `json:"id"`       // [READ-ONLY] global identification
	FolderId    KId            `json:"folderId"` // [REQUIRED FOR CREATE] [WRITE-ONCE] global identification of folder in which is the event defined
	Watermark   Watermark      `json:"watermark"`
	Access      EventAccess    `json:"access"` // [READ-ONLY] scope of access of user to this event
	Summary     string         `json:"summary"`
	Location    string         `json:"location"`
	Description string         `json:"description"`
	Status      TaskStatus     `json:"status"`
	Start       UtcDateTime    `json:"start"`
	Due         UtcDateTime    `json:"due"`  // Deadline
	End         UtcDateTime    `json:"end"`  // [READ-ONLY] Date when task was completed. Valid only if the status is 'tsCompleted'.
	Done        int            `json:"done"` // Percent completed. If the status is set to 'tsCompleted' this value is always set to 100%.
	Priority    PriorityType   `json:"priority"`
	Rule        RecurrenceRule `json:"rule"`
	Attendees   AttendeeList   `json:"attendees"`
	Reminder    Reminder       `json:"reminder"`
	SortOrder   int            `json:"sortOrder"` // [0-7FFFFFFF] Zero means a newest tasks.
	IsPrivate   bool           `json:"isPrivate"`
	IsCancelled bool           `json:"isCancelled"` // [READ-ONLY] is canceled by organizer
}

// TaskList - List of resources
type TaskList []Task

// Constants for composing kerio::web::SearchQuery
// Tasks management.

// TasksGet - Get a list of tasks.
//	folderIds - list of global identifiers of folders to be listed.
//	query - query attributes and limits
// Return
//	list - all found tasks
//  totalItems - number of tasks found if there is no limit
func (c *ClientConnection) TasksGet(folderIds KIdList, query SearchQuery) (TaskList, int, error) {
	query = addMissedParametersToSearchQuery(query)
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Tasks.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       TaskList `json:"list"`
			TotalItems int      `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// TasksGetById - Get an tasks.
//	ids - global identifiers of requested tasks
// Return
//	errors - list of tasks that failed to obtain
//	result - found tasks
func (c *ClientConnection) TasksGetById(ids KIdList) (ErrorList, TaskList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Tasks.getById", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result TaskList  `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// TasksRemove - Remove a list of tasks.
//	ids - list of global identifiers of tasks to be removed
// Return
//	errors - list of tasks that failed to remove
func (c *ClientConnection) TasksRemove(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Tasks.remove", params)
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

// TasksCopy - Copy existing tasks to folder
//	ids - list of global identifiers of tasks to be copied
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) TasksCopy(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Tasks.copy", params)
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

// TasksCreate - Create tasks.
// Return
//	errors - list of tasks that failed on creation
//	result - particular results for all items
func (c *ClientConnection) TasksCreate(tasks TaskList) (ErrorList, CreateResultList, error) {
	params := struct {
		Tasks TaskList `json:"tasks"`
	}{tasks}
	data, err := c.CallRaw("Tasks.create", params)
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

// TasksSet - Set tasks.
// Return
//	errors - error message list
func (c *ClientConnection) TasksSet(tasks TaskList) (ErrorList, SetResultList, error) {
	params := struct {
		Tasks TaskList `json:"tasks"`
	}{tasks}
	data, err := c.CallRaw("Tasks.set", params)
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

// TasksMove - Move existing tasks to folder
//	ids - list of global identifiers of tasks to be moved
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) TasksMove(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Tasks.move", params)
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
