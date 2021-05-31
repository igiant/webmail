package webmail

import "encoding/json"

type Validity struct {
	IsValid bool               `json:"isValid"`
	Error   LocalizableMessage `json:"error"`
}

type NameEntry struct {
	CommonName             string     `json:"commonName"`             // CN
	OrganizationName       string     `json:"organizationName"`       // O
	OrganizationalUnitName string     `json:"organizationalUnitName"` // OU
	LocalityName           string     `json:"localityName"`           // L
	CountryName            string     `json:"countryName"`            // C
	StateOrProvinceName    string     `json:"stateOrProvinceName"`    // ST
	EmailAddresses         StringList `json:"emailAddresses"`         // 'emailAddress' or 'subjectAltName'(X509v3 Subject Alternative Name)
}

type Certificate struct {
	Id        KId         `json:"id"`
	Subject   NameEntry   `json:"subject"`
	Issuer    NameEntry   `json:"issuer"`
	ValidFrom UtcDateTime `json:"validFrom"`
	ValidTo   UtcDateTime `json:"validTo"`
	Serial    string      `json:"serial"`
	Validity  Validity    `json:"validity"`
}

type CertificateList []Certificate

type CertStoreStatus string

const (
	Uninitialized CertStoreStatus = "Uninitialized" // The user has not personal certificate store initialized yet.
	Opened        CertStoreStatus = "Opened"        // The personal certificate store is not opened. The mails are automatically decrypted. Signing and encrypting is possible.
	Closed        CertStoreStatus = "Closed"        // The personal certificate store is closed.
	FailedToOpen  CertStoreStatus = "FailedToOpen"  // Failed to open it during login in. Valid only if the user uses the login password for the personal certificate store.
)

// the personal certificate store manager class

// CertificatesInit - Initialize the personal certificate store
// Parameters
//	password - password of certificate store
//	isLoginPassword - given password is the same which user uses to log in
func (c *ClientConnection) CertificatesInit(password string, isLoginPassword bool) error {
	params := struct {
		Password        string `json:"password"`
		IsLoginPassword bool   `json:"isLoginPassword"`
	}{password, isLoginPassword}
	_, err := c.CallRaw("Certificates.init", params)
	return err
}

// CertificatesOpen - Open the personal certificate store
// Parameters
//	password - password of certificate store
func (c *ClientConnection) CertificatesOpen(password string) error {
	params := struct {
		Password string `json:"password"`
	}{password}
	_, err := c.CallRaw("Certificates.open", params)
	return err
}

// CertificatesClose - Close the personal certificate store
func (c *ClientConnection) CertificatesClose() error {
	_, err := c.CallRaw("Certificates.close", nil)
	return err
}

// CertificatesGet - Obtain a list of certificates
// Return
//	certificates - current list of certificates
func (c *ClientConnection) CertificatesGet() (CertificateList, error) {
	data, err := c.CallRaw("Certificates.get", nil)
	if err != nil {
		return nil, err
	}
	certificates := struct {
		Result struct {
			Certificates CertificateList `json:"certificates"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &certificates)
	return certificates.Result.Certificates, err
}

// CertificatesGetById - Obtain particular certificate
// Return
//	certificate - a certificate
//	certificate - global identifier
func (c *ClientConnection) CertificatesGetById(id KId) (*Certificate, error) {
	params := struct {
		Id KId `json:"id"`
	}{id}
	data, err := c.CallRaw("Certificates.getById", params)
	if err != nil {
		return nil, err
	}
	certificate := struct {
		Result struct {
			Certificate Certificate `json:"certificate"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &certificate)
	return &certificate.Result.Certificate, err
}

// CertificatesGetStatus - Obtain a list of certificates
func (c *ClientConnection) CertificatesGetStatus() (*CertStoreStatus, error) {
	data, err := c.CallRaw("Certificates.getStatus", nil)
	if err != nil {
		return nil, err
	}
	status := struct {
		Result struct {
			Status CertStoreStatus `json:"status"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &status)
	return &status.Result.Status, err
}

// CertificatesToSource - Obtain source (plain-text representation) of the certificate
// Parameters
//	id - global identifier
// Return
//	source - certificate in plain text
func (c *ClientConnection) CertificatesToSource(id KId) (string, error) {
	params := struct {
		Id KId `json:"id"`
	}{id}
	data, err := c.CallRaw("Certificates.toSource", params)
	if err != nil {
		return "", err
	}
	source := struct {
		Result struct {
			Source string `json:"source"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &source)
	return source.Result.Source, err
}

// CertificatesOpenWithOldLoginPassword - Calling is valid only if login password is used as well for certificate store.
// Parameters
//	oldPassword - password to certificate store (old login password)
func (c *ClientConnection) CertificatesOpenWithOldLoginPassword(oldPassword string) error {
	params := struct {
		OldPassword string `json:"oldPassword"`
	}{oldPassword}
	_, err := c.CallRaw("Certificates.openWithOldLoginPassword", params)
	return err
}

// CertificatesOpenEditWithOldLoginPassword - Calling is valid only if login password is used as well for certificate store.
func (c *ClientConnection) CertificatesOpenEditWithOldLoginPassword(oldPassword string) error {
	params := struct {
		OldPassword string `json:"oldPassword"`
	}{oldPassword}
	_, err := c.CallRaw("Certificates.openEditWithOldLoginPassword", params)
	return err
}

// CertificatesReset - Reset personal certificate store to uninitialized state. All current store will be removed!
// Parameters
//	loginPassword - current login password to verify user)
func (c *ClientConnection) CertificatesReset(loginPassword string) error {
	params := struct {
		LoginPassword string `json:"loginPassword"`
	}{loginPassword}
	_, err := c.CallRaw("Certificates.reset", params)
	return err
}

// CertificatesOpenEdit - Unlock edit functions
// Parameters
//	password - password of certificate store
func (c *ClientConnection) CertificatesOpenEdit(password string) error {
	params := struct {
		Password string `json:"password"`
	}{password}
	_, err := c.CallRaw("Certificates.openEdit", params)
	return err
}

// CertificatesCloseEdit - Lock edit functions
func (c *ClientConnection) CertificatesCloseEdit() error {
	_, err := c.CallRaw("Certificates.closeEdit", nil)
	return err
}

// CertificatesSetPreferred - Preferred flag is removed from other certificates issued for the same email address.
// Parameters
//	id - ID of the certificate
func (c *ClientConnection) CertificatesSetPreferred(id KId) error {
	params := struct {
		Id KId `json:"id"`
	}{id}
	_, err := c.CallRaw("Certificates.setPreferred", params)
	return err
}

// CertificatesChangePassword - Preferred flag is removed from other certificates issued for the same email address.
func (c *ClientConnection) CertificatesChangePassword(oldPassword string, newPassword string, isLoginPassword bool) error {
	params := struct {
		OldPassword     string `json:"oldPassword"`
		NewPassword     string `json:"newPassword"`
		IsLoginPassword bool   `json:"isLoginPassword"`
	}{oldPassword, newPassword, isLoginPassword}
	_, err := c.CallRaw("Certificates.changePassword", params)
	return err
}

// CertificatesImportPKCS12 - Preferred flag is removed from other certificates issued for the same email address.
func (c *ClientConnection) CertificatesImportPKCS12(fileId KId, password string) error {
	params := struct {
		FileId   KId    `json:"fileId"`
		Password string `json:"password"`
	}{fileId, password}
	_, err := c.CallRaw("Certificates.importPKCS12", params)
	return err
}

// CertificatesExportPKCS12 - Note: "export" is a keyword in C++, so the name of the method must be changed: exportPrivateKey
// Parameters
//	id - ID of the certificate
// Return
//	fileDownload - description of the output file
func (c *ClientConnection) CertificatesExportPKCS12(newPassword string, id KId) (*Download, error) {
	params := struct {
		NewPassword string `json:"newPassword"`
		Id          KId    `json:"id"`
	}{newPassword, id}
	data, err := c.CallRaw("Certificates.exportPKCS12", params)
	if err != nil {
		return nil, err
	}
	fileDownload := struct {
		Result struct {
			FileDownload Download `json:"fileDownload"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(data, &fileDownload)
	return &fileDownload.Result.FileDownload, err
}

// CertificatesRemove - Note: "export" is a keyword in C++, so the name of the method must be changed: exportPrivateKey
func (c *ClientConnection) CertificatesRemove(id KId) error {
	params := struct {
		Id KId `json:"id"`
	}{id}
	_, err := c.CallRaw("Certificates.remove", params)
	return err
}
