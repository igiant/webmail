package webmail

import "encoding/json"

func marshal(id int, method string, token *string, params interface{}) ([]byte, error) {
	p := parameters{
		JsonRpc: "2.0",
		Method:  method,
		ID:      id,
		Token:   token,
		Params:  params,
	}
	return json.Marshal(&p)
}
