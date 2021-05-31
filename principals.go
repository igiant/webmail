package webmail

import "encoding/json"

type PrincipalType string

const (
	ptUser      PrincipalType = "ptUser"
	ptResource  PrincipalType = "ptResource"
	ptGroup     PrincipalType = "ptGroup"
	ptDomain    PrincipalType = "ptDomain"
	ptAnonymous PrincipalType = "ptAnonymous" // Special type without ID: anyone
	ptAuthUser  PrincipalType = "ptAuthUser"  // Special type without ID: every authenticated user
)

type Principal struct {
	Id          KId           `json:"id"` // global identification
	Type        PrincipalType `json:"type"`
	DisplayName string        `json:"displayName"`
	MailAddress string        `json:"mailAddress"`
}

type PrincipalList []Principal

// PrincipalsGet - Get list of principals from server.
// Return
//	list - principals
func (c *ClientConnection) PrincipalsGet(users bool, groups bool, domains bool) (PrincipalList, error) {
	params := struct {
		Users   bool `json:"users"`
		Groups  bool `json:"groups"`
		Domains bool `json:"domains"`
	}{users, groups, domains}
	data, err := c.CallRaw("Principals.get", params)
	if err != nil {
		return nil, err
	}
	list := struct {
		Result struct {
			List PrincipalList `json:"list"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &list)
	return list.Result.List, err
}

// PrincipalsGetByEmail - Find principal according his primary email (login name) on server.
// Parameters
//	email - email/login name
// Return
//	principal - principal
func (c *ClientConnection) PrincipalsGetByEmail(email string) (*Principal, error) {
	params := struct {
		Email string `json:"email"`
	}{email}
	data, err := c.CallRaw("Principals.getByEmail", params)
	if err != nil {
		return nil, err
	}
	principal := struct {
		Result struct {
			Principal Principal `json:"principal"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &principal)
	return &principal.Result.Principal, err
}
