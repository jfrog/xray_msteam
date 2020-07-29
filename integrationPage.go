package main

import (
	"net/http"
	"io/ioutil"
	"io"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
	"errors"
)

type Violation struct {
	Created            string    `json:"created"`
	TopSeverity        string    `json:"top_severity"`
	WatchName          string    `json:"watch_name"`
	PolicyName         string    `json:"policy_name"`
	Issues             Issues    `json:"issues"`
}


type Issue struct {
	Severity           string             `json:"severity"`
	Type               string             `json:"type"` // Issue type license/security
	Summary            string             `json:"summary"`
	Description        string             `json:"description"`
	ImpactedArtifacts  ImpactedArtifacts  `json:"impacted_artifacts"`
}

type Issues []Issue

type ImpactedArtifact struct {
	Name             string         `json:"name"` // Artifact name
	DisplayName      string         `json:"display_name"`
	Path             string         `json:"path"`  // Artifact path in Artifactory
	PackageType      string         `json:"pkg_type"`
	SHA256           string         `json:"sha256"` // Artifact SHA 256 checksum
	SHA1             string         `json:"sha1"`
	Depth            int            `json:"depth"`  // Artifact depth in its hierarchy
	ParentSHA        string         `json:"parent_sha"`
	InfectedFiles    InfectedFiles  `json:"infected_files"`
}

type ImpactedArtifacts []ImpactedArtifact

type InfectedFile struct {
	Name           string    `json:"name"`
	Path           string    `json:"path"`  // artifact path in Artifactory
	SHA256         string    `json:"sha256"`// artifact SHA 256 checksum
	Depth          int       `json:"depth"` // Artifact depth in its hierarchy
	ParentSHA      string    `json:"parent_sha"` // Parent artifact SHA1
	DisplayName    string    `json:"display_name"`
	PackageType    string    `json:"pkg_type"`
}

type InfectedFiles []InfectedFile

type TeamMessage struct {
	Content string
}

type TeamPayload struct {
	Body TeamMessage
}

func PingPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func SendPage(w http.ResponseWriter, r *http.Request) {
	err := SendMessage(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		w.WriteHeader(200)
	}
}


func SendMessage(r *http.Request) error {
	var violation Violation

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 5048576))
	if err != nil{
		return err
	}

	err = json.Unmarshal(body, &violation)
	if err != nil {
		return err
	}

	if len(violation.PolicyName) > 0 {
		return errors.New("Unable to read webhook payload for critical data to send to MS Teams")
	}

	// VIOLATION PAYLOAD
	violationMessage := fmt.Sprintf("🔔\nPolicy: %s \nWatch: %s \nCreated: %s \nNumber Of Issues: %d", violation.PolicyName, violation.WatchName, violation.Created, len(violation.Issues))
	count := 0
	for _, thisIssue := range violation.Issues {
		count++
		if count <= 5 {
			violationMessage = fmt.Sprintf("%s\nIssue: %s",violationMessage,thisIssue.Summary)
		}
	}
	if count > 5 {
		violationMessage = fmt.Sprintf("%s\nAdditional issues available but not displayed under watch: %s",violationMessage, violation.WatchName)
	}
	fmt.Println(violationMessage)

	client := resty.New()

	// Bearer Auth Token for all request
	client.SetAuthToken(MicrosoftAccessToken)

	// Retries are configured per client
	client.
	// Set retry count to non zero to enable retries
	SetRetryCount(3).
	// You can override initial retry wait time.
	// Default is 100 milliseconds.
	SetRetryWaitTime(5 * time.Second).
	// MaxWaitTime can be overridden as well.
	// Default is 2 seconds.
	SetRetryMaxWaitTime(20 * time.Second)

	message := TeamMessage{Content: violationMessage}
	teamPayload := TeamPayload{Body: message}

	// Marshal to JSON payload
	payload, payloadErr := json.Marshal(teamPayload)
	if payloadErr != nil {
		return payloadErr
	}

	// Send Payload toe Microsoft Graph API to post message to Teams
	url := fmt.Sprintf("https://graph.microsoft.com/v1.0/teams/%s/channels/%s/messages", MicrosoftTeamId, MicrosoftTeamChannelId)
	resp, errored := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(string(payload)).
					Post(url)

	fmt.Println(resp)

	if errored != nil {
		return errored
	}
	return nil
}