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
}

// NewExternalAPI create a new ExternalAPI from authentication information
func NewExternalAPI(address, username, password string) (ExternalAPI, error) {
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
	}, nil
}

// authenticate authenticates user to Docktor through its login API
// It updates self object with the generated JWT Token, used to authenticate through all protected API routes
func (api *ExternalAPI) authenticate() ([]*http.Cookie, error) {
	data := url.Values{}
	data.Set("username", api.username)
	data.Set("password", api.password)

	u, err := url.ParseRequestURI(api.address)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse Docktor URL: %v", err.Error())
	}
	u.Path = "/auth/signin"

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Failed to initiate HTTP request: %v", err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	log.WithFields(log.Fields{
		"address":  api.address,
		"username": api.username,
	}).Debug("Authenticating to Docktor")

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"address":  api.address,
			"username": api.username,
		}).WithError(err).Error("Failed to authenticate to Docktor")
		return nil, fmt.Errorf("Failed to authenticate to Docktor: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"address":    api.address,
			"username":   api.username,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
		}).Error("Failed to authenticate to Docktor because server did not return OK")
		return nil, fmt.Errorf("Failed to authenticate to Docktor because server did not return OK: %v", resp.Status)
	}

	log.WithFields(log.Fields{
		"address":  api.address,
		"username": api.username,
	}).Debug("Docktor authentication successful")

	return resp.Cookies(), nil
}

// GroupDocktor is a group fetched from Docktor API
type GroupDocktor struct {
	ID         string `json:"_id,omitempty"`
	Title      string `json:"title,omitempty"`
	Containers []struct {
		ServiceTitle string `json:"serviceTitle"`
	} `json:"containers"`
}

// GetGroup gets a Docktor group name from its ID
func (api *ExternalAPI) GetGroup(groupID string) (GroupDocktor, error) {
	cookies, err := api.authenticate()
	if err != nil {
		return GroupDocktor{}, err
	}

	u, _ := url.ParseRequestURI(api.address)
	u.Path = fmt.Sprintf("/groups/%v", groupID)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", u.String(), nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

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
		"address":    api.address,
		"groupID":    docktorGroup.ID,
		"groupTitle": docktorGroup.Title,
	}).Debug("Fetch group from Docktor")

	return docktorGroup, nil

}

// GetGroupIDFromURL returns the Docktor group ID from its URL
// URL is expected to be format : http://<docktor-host>/groups/<id>
func (api *ExternalAPI) GetGroupIDFromURL(docktorURL string) (string, error) {
	u, err := url.ParseRequestURI(docktorURL)
	if err != nil {
		return "", fmt.Errorf("docktorGroupURL is not a valid URL. Expected 'http://<docktor>/groups/<id>', Got '%v'", docktorURL)
	}
	path := strings.Split(u.Path, "/")
	if len(path) == 0 {
		return "", fmt.Errorf("Unable to get project id from URL. Expected 'http://<docktor>/groups/<id>', Got '%v'", u.Path)
	}
	id := path[len(path)-1]
	if id == "" {
		return "", fmt.Errorf("Unable to get project id from URL parsed path : %v. URL=%v", path, u.Path)
	}
	return id, nil
}
