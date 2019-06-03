package docktor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

const timeout = time.Duration(15 * time.Second)

// ExternalAPI exposes methods to query the Docktor API
type ExternalAPI struct {
	address  string
	username string
	password string
	ldap     bool
}

// NewExternalAPI create a new ExternalAPI from authentication information
func NewExternalAPI(address, username, password string, ldap bool) (ExternalAPI, error) {
	if address == "" {
		return ExternalAPI{}, fmt.Errorf("Docktor address is empty")
	} else if username == "" {
		return ExternalAPI{}, fmt.Errorf("Docktor username is empty")
	} else if password == "" {
		return ExternalAPI{}, fmt.Errorf("Docktor password is empty")
	}

	return ExternalAPI{
		address:  address,
		username: username,
		password: password,
		ldap:     ldap,
	}, nil
}

// authenticate authenticates user to Docktor through its login API
// It updates self object with the generated JWT Token, used to authenticate through all protected API routes
func (api *ExternalAPI) authenticate() (string, error) {

	values := map[string]string{"username": api.username, "password": api.password}
	js, err := json.Marshal(values)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal username and password: %v", err.Error())
	}

	u, err := url.ParseRequestURI(api.address)
	if err != nil {
		return "", fmt.Errorf("Failed to parse Docktor URL: %v", err.Error())
	}
	u.Path = fmt.Sprintf("/api/auth/login?ldap=%v", api.ldap)

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(js))
	if err != nil {
		return "", fmt.Errorf("Failed to initiate HTTP request: %v", err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(len(js)))

	log.WithFields(log.Fields{
		"address":  api.address,
		"username": api.username,
		"ldap":     api.ldap,
	}).Debug("Authenticating to Docktor")

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"address":  api.address,
			"username": api.username,
			"ldap":     api.ldap,
		}).WithError(err).Error("Failed to authenticate to Docktor")
		return "", fmt.Errorf("Failed to authenticate to Docktor: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"address":    api.address,
			"username":   api.username,
			"ldap":       api.ldap,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
		}).Error("Failed to authenticate to Docktor because server did not return OK")
		return "", fmt.Errorf("Failed to authenticate to Docktor because server did not return OK: %v", resp.Status)
	}

	log.WithFields(log.Fields{
		"address":  api.address,
		"username": api.username,
		"ldap":     api.ldap,
	}).Debug("Docktor authentication successful")

	var token string

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		log.WithFields(log.Fields{
			"address":  api.address,
			"username": api.username,
			"ldap":     api.ldap,
		}).WithError(err).Error("Failed to get token from Docktor while decoding JSON result")
		return "", fmt.Errorf("Failed to get token from Docktor while decoding JSON result: %v", err.Error())
	}

	return token, nil
}

// GroupDocktor is a group fetched from Docktor API
type GroupDocktor struct {
	ID         string `json:"_id,omitempty"`
	Name       string `json:"name,omitempty"`
	Containers []struct {
		Image string `json:"Image"`
	} `json:"containers"`
}

// GetGroup gets a Docktor group name from its ID
func (api *ExternalAPI) GetGroup(groupID string) (GroupDocktor, error) {
	token, err := api.authenticate()
	if err != nil {
		return GroupDocktor{}, err
	}

	u, _ := url.ParseRequestURI(api.address)
	u.Path = fmt.Sprintf("/api/groups/%v", groupID)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", u.String(), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	log.WithFields(log.Fields{
		"url":     u.String(),
		"groupID": groupID,
	}).Debugf("Getting group from Docktor")

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"docktorURL": u.String(),
			"user":       api.username,
		}).Error("Failed to get group from Docktor")
		return GroupDocktor{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"address":    api.address,
			"username":   api.username,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
		}).Error("Failed to get group from Docktor because server did not return OK")
		return GroupDocktor{}, fmt.Errorf("Failed to get group from Docktor because server did not return OK: %v", resp.Status)
	}

	var docktorGroup GroupDocktor

	err = json.NewDecoder(resp.Body).Decode(&docktorGroup)
	if err != nil {
		log.WithFields(log.Fields{
			"address":  api.address,
			"username": api.username,
		}).WithError(err).Error("Failed to get group from Docktor while decoding JSON result")
		return GroupDocktor{}, fmt.Errorf("Failed to get group from Docktor while decoding JSON result: %v", err.Error())
	}

	log.WithFields(log.Fields{
		"address":   api.address,
		"groupID":   docktorGroup.ID,
		"groupName": docktorGroup.Name,
	}).Debug("Fetch group from Docktor")

	return docktorGroup, nil

}

// GetGroupIDFromURL returns the Docktor group ID from its URL
// URL is expected to be format : https://<docktor-host>/groups/<id>
func (api *ExternalAPI) GetGroupIDFromURL(docktorURL string) (string, error) {
	u, err := url.ParseRequestURI(docktorURL)
	if err != nil {
		return "", fmt.Errorf("docktorGroupURL is not a valid URL. Expected 'https://<docktor>/groups/<id>', Got '%v'", docktorURL)
	}
	path := strings.Split(u.Path, "/")
	if len(path) == 0 {
		return "", fmt.Errorf("Unable to get project id from URL. Expected 'https://<docktor>/groups/<id>', Got '%v'", u.Path)
	}
	id := path[len(path)-1]
	if id == "" {
		return "", fmt.Errorf("Unable to get project id from URL parsed path : %v. URL=%v", path, u.Path)
	}
	return id, nil
}
