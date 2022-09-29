package model

import "time"

type EventBridge struct {
	EventVersion string `json:"eventVersion"`
	UserIdentity struct {
		Type           string `json:"type"`
		PrincipalID    string `json:"principalId"`
		Arn            string `json:"arn"`
		AccountID      string `json:"accountId"`
		AccessKeyID    string `json:"accessKeyId"`
		UserName       string `json:"userName"`
		SessionContext struct {
			SessionIssuer struct {
			} `json:"sessionIssuer"`
			WebIDFederationData struct {
			} `json:"webIdFederationData"`
			Attributes struct {
				CreationDate     time.Time `json:"creationDate"`
				MfaAuthenticated string    `json:"mfaAuthenticated"`
			} `json:"attributes"`
		} `json:"sessionContext"`
	} `json:"userIdentity"`
	EventTime         time.Time `json:"eventTime"`
	EventSource       string    `json:"eventSource"`
	EventName         string    `json:"eventName"`
	AwsRegion         string    `json:"awsRegion"`
	SourceIPAddress   string    `json:"sourceIPAddress"`
	UserAgent         string    `json:"userAgent"`
	RequestParameters struct {
		Language     string `json:"language"`
		SeverityCode string `json:"severityCode"`
		ServiceCode  string `json:"serviceCode"`
		IssueType    string `json:"issueType"`
		CategoryCode string `json:"categoryCode"`
	} `json:"requestParameters"`
	ResponseElements struct {
		CaseID string `json:"caseId"`
	} `json:"responseElements"`
	RequestID                    string `json:"requestID"`
	EventID                      string `json:"eventID"`
	ReadOnly                     bool   `json:"readOnly"`
	EventType                    string `json:"eventType"`
	ManagementEvent              bool   `json:"managementEvent"`
	RecipientAccountID           string `json:"recipientAccountId"`
	EventCategory                string `json:"eventCategory"`
	SessionCredentialFromConsole string `json:"sessionCredentialFromConsole"`
}

type Detail struct {
	CaseID          string `json:"case-id"`
	DisplayID       string `json:"display-id"`
	CommunicationID string `json:"communication-id"`
	EventName       string `json:"event-name"`
	Origin          string `json:"origin"`
}