package webmail

import (
	"encoding/json"
	"fmt"
)

type ErrorReport struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		MessageParameters struct {
			PositionalParameters []string `json:"positionalParameters"`
			Plurality            int      `json:"plurality"`
		} `json:"messageParameters"`
	} `json:"data"`
}

type errorReport struct {
	ErrorReport `json:"error"`
}

func checkError(data []byte) error {
	errorReport := errorReport{}
	_ = json.Unmarshal(data, &errorReport)
	if errorReport.Code == 0 && errorReport.Message == "" {
		return nil
	}
	return fmt.Errorf("%d: %s", errorReport.Code, errorReport.Message)
}
