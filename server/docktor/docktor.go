package docktor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

const timeout = time.Duration(5 * time.Second)

// ExternalAPI expose methods of Docktor API
// For instance, it exposes method to get the group name from an ID of a Docktor group
type ExternalAPI struct {
	docktorAddr     string
	docktorUser     string
	docktorPassword string
}

// NewExternalAPI create a new ExternalAPI from authentication information
// addr is the address of the Docktor instance
// user is the username used to authenticate to Docktor instance
// password is the password used in association to user
func NewExternalAPI(addr, user, password string) (ExternalAPI, error) {
	if addr == "" || user == "" || password == "" {
		return ExternalAPI{}, fmt.Errorf("Docktor address, user and password are mandatory. Addr=%v User=%v", addr, user)
	}
	return ExternalAPI{
		docktorAddr:     addr,
		docktorUser:     user,
		docktorPassword: password,
	}, nil
}

// authenticate authenticates user to Docktor through its login API
// It updates self object with the generated JWT Token, used to authenticate through all protected API routes
func (api *ExternalAPI) authenticate() ([]*http.Cookie, error) {

	data := url.Values{}
	data.Set("username", api.docktorUser)
	data.Set("password", api.docktorPassword)

	u, _ := url.ParseRequestURI(api.docktorAddr)
	u.Path = "/auth/signin"

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("POST", u.String(), bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	log.WithFields(log.Fields{
		"DocktorURL": api.docktorAddr,
		"user":       api.docktorUser,
	}).Debug("Authenticating to Docktor...")

	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"DocktorURL": api.docktorAddr,
			"user":       api.docktorUser,
		}).WithError(err).Error("Failed to authenticate to Docktor")
		return nil, fmt.Errorf("Failed to authenticate to Docktor: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"DocktorURL":  api.docktorAddr,
			"user":        api.docktorUser,
			"HTTP_STATUS": resp.StatusCode,
			"Status":      resp.Status,
		}).Error("Failed to authenticate to Docktor because server did not return OK")
		return nil, fmt.Errorf("Failed to authenticate to Docktor because server did not return OK: %v", resp.Status)
	}

	log.WithFields(log.Fields{
		"DocktorURL": api.docktorAddr,
		"user":       api.docktorUser,
	}).Debug("Authenticating to Docktor [OK]")

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

	u, _ := url.ParseRequestURI(api.docktorAddr)
	u.Path = fmt.Sprintf("/groups/%v", groupID)

	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", u.String(), nil)
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	log.WithField("url", u.String()).Debugf("Getting group from Docktor with id=%v ...", groupID)

	resp, err := client.Do(req)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"docktorURL": u.String(),
			"user":       api.docktorUser,
		}).Error("Failed to get group from Docktor")
		return GroupDocktor{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.WithFields(log.Fields{
			"docktorURL":  api.docktorAddr,
			"user":        api.docktorUser,
			"HTTP_STATUS": resp.StatusCode,
			"Status":      resp.Status,
		}).Error("Failed to get group from Docktor because server did not return OK")
		return GroupDocktor{}, fmt.Errorf("Failed to get group from Docktor because server did not return OK: %v", resp.Status)
	}

	var docktorGroup GroupDocktor

	err = json.NewDecoder(resp.Body).Decode(&docktorGroup)
	if err != nil {
		log.WithFields(log.Fields{
			"docktorURL": api.docktorAddr,
			"user":       api.docktorUser,
		}).WithError(err).Error("Failed to get group from Docktor while decoding JSON result")
		return GroupDocktor{}, fmt.Errorf("Failed to get group from Docktor while decoding JSON result: %v", err.Error())
	}

	log.WithFields(log.Fields{
		"docktorURL":  api.docktorAddr,
		"group.id":    docktorGroup.ID,
		"group.title": docktorGroup.Title,
	}).Debug("Fetch group from Docktor")

	return docktorGroup, nil

}
