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
	Cve                string             `json:"cve"`
}

type Issues []Issue

type TeamMessageActionTarget struct {
	Os  string `json:"os"`
	Uri string `json:"uri"`
}

type TeamMessageAction struct {
	Type string `json:"@type"`
	Name string `json:"name"`
	Targets []TeamMessageActionTarget `json:"targets"`
}

type TeamPayload struct {
	Context string `json:"@context"`
	Type string `json:"@type"`
	ThemeColor string `json:"themeColor"`
	Title string `json:"title"`
	Text string `json:"text"`
	PotentialActions []TeamMessageAction `json:"potentialAction"`
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

	if len(violation.PolicyName) == 0 {
		return errors.New("Unable to read webhook payload for critical data to send to MS Teams")
	}

	// VIOLATION PAYLOAD
	violationMessage := fmt.Sprintf("🔔 Policy: %s 🕐 Watch: %s ⌚ Created: %s 🔢 Number Of Issues: %d", violation.PolicyName, violation.WatchName, violation.Created, len(violation.Issues))
	count := 0
	issueMessage := ""
	for _, thisIssue := range violation.Issues {
		count++
		if count <= 5 {
			issueMessage = fmt.Sprintf("%s<br/>⚙️ Issue: %s",issueMessage,thisIssue.Summary)
		}
	}
	if count > 5 {
		issueMessage = fmt.Sprintf("%s\nAdditional issues available but not displayed under watch: %s",issueMessage, violation.WatchName)
	}
	client := resty.New()

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


	firstCve := ""
	if len(violation.Issues) > 0 {
		firstCve = violation.Issues[0].Cve
	}

	actionTarget := fmt.Sprintf("https://cve.mitre.org/cgi-bin/cvename.cgi?name=%s", firstCve)
	target := TeamMessageActionTarget{Os:"default",Uri:actionTarget}
	targets := []TeamMessageActionTarget{target}
	action := TeamMessageAction{Type: "OpenUri", Name: "Research more...", Targets: targets}
	actions := []TeamMessageAction{action}
	teamPayload := TeamPayload{Context:"https://schema.org/extensions",Type:"MessageCard",ThemeColor:"0ac70d",Title:violationMessage,Text:issueMessage,PotentialActions:actions}

	// Marshal to JSON payload
	payload, payloadErr := json.Marshal(teamPayload)
	if payloadErr != nil {
		return payloadErr
	}

	// Send Payload toe Microsoft Teams Channel Webhook to post message to Teams Channel setup by the user
	_, errored := client.R().
					SetHeader("Content-Type", "application/json").
					SetBody(string(payload)).
					Post(MicrosoftTeamWebhook)
	if errored != nil {
		return errored
	}
	return nil
}