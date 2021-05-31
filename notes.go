package webmail

// NoteColor - Color of note
type NoteColor string

const (
	White  NoteColor = "White"
	Yellow NoteColor = "Yellow"
	Pink   NoteColor = "Pink"
	Green  NoteColor = "Green"
	Blue   NoteColor = "Blue"
)

// NotePosition - Position of note
type NotePosition struct {
	XOffset unsignedlong `json:"xOffset"`
	YOffset unsignedlong `json:"yOffset"`
	XSize   unsignedlong `json:"xSize"`
	YSize   unsignedlong `json:"ySize"`
}

// Note - Note details
type Note struct {
	Id         KId          `json:"id"`       // [READ-ONLY] global identification
	FolderId   KId          `json:"folderId"` // [REQUIRED FOR CREATE] [WRITE-ONCE] global identification of folder in which is the note defined
	Watermark  Watermark    `json:"watermark"`
	Color      NoteColor    `json:"color"`      // COLOR:BLUE
	Text       string       `json:"text"`       // TEXT:****
	Position   NotePosition `json:"position"`   // POSITION:354 206 623 326
	CreateDate UtcDateTime  `json:"createDate"` // [READ-ONLY]
	ModifyDate UtcDateTime  `json:"modifyDate"` // [READ-ONLY]
}

// NoteList - List of notes
type NoteList []Note

// Constants for composing kerio::web::SearchQuery
// Notes management.

// NotesGet - Get a list of notes.
// Parameters
//	folderIds - list of global identifiers of folders to be listed.
//	query - query attributes and limits
// Return
//	list - all found notes
//	totalItems - number of notes found if there is no limit
func (c *ClientConnection) NotesGet(folderIds KIdList, query SearchQuery) (NoteList, int, error) {
	params := struct {
		FolderIds KIdList     `json:"folderIds"`
		Query     SearchQuery `json:"query"`
	}{folderIds, query}
	data, err := c.CallRaw("Notes.get", params)
	if err != nil {
		return nil, 0, err
	}
	list := struct {
		Result struct {
			List       NoteList `json:"list"`
			TotalItems int      `json:"totalItems"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, list.Result.TotalItems, err
}

// NotesGetById - Get an note.
// Parameters
//	ids - global identifiers of requested notes
// Return
//	errors - list of errors
//	result - found notes
func (c *ClientConnection) NotesGetById(ids KIdList) (ErrorList, NoteList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Notes.getById", params)
	if err != nil {
		return nil, nil, err
	}
	errors := struct {
		Result struct {
			Errors ErrorList `json:"errors"`
			Result NoteList  `json:"result"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &errors)
	return errors.Result.Errors, errors.Result.Result, err
}

// NotesRemove - Remove a list of notes.
// Parameters
//	ids - list of global identifiers of notes to be removed
// Return
//	errors - list of notes that failed to remove
func (c *ClientConnection) NotesRemove(ids KIdList) (ErrorList, error) {
	params := struct {
		Ids KIdList `json:"ids"`
	}{ids}
	data, err := c.CallRaw("Notes.remove", params)
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

// NotesCopy - Copy existing notes to folder
// Parameters
//	ids - list of global identifiers of notes to be copied
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) NotesCopy(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Notes.copy", params)
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

// NotesCreate - Create notes.
// Parameters
//	notes - list of notes to be created
// Return
//	errors - list of notes that failed on creation
//	result - particular results for all items
func (c *ClientConnection) NotesCreate(notes NoteList) (ErrorList, CreateResultList, error) {
	params := struct {
		Notes NoteList `json:"notes"`
	}{notes}
	data, err := c.CallRaw("Notes.create", params)
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

// NotesSet - Set notes.
// Return
//	errors - error message list
func (c *ClientConnection) NotesSet(notes NoteList) (ErrorList, SetResultList, error) {
	params := struct {
		Notes NoteList `json:"notes"`
	}{notes}
	data, err := c.CallRaw("Notes.set", params)
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

// NotesMove - Move existing notes to folder
// Parameters
//	ids - list of global identifiers of notes to be moved
//	folder - target folder
// Return
//	errors - error message list
func (c *ClientConnection) NotesMove(ids KIdList, folder KId) (ErrorList, CreateResultList, error) {
	params := struct {
		Ids    KIdList `json:"ids"`
		Folder KId     `json:"folder"`
	}{ids, folder}
	data, err := c.CallRaw("Notes.move", params)
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
