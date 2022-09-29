package model

import "time"

type FS struct {
	Description string   `json:"description"`
	Subject     string   `json:"subject"`
	Email       string   `json:"email"`
	Priority    int      `json:"priority"`
	Status      int      `json:"status"`
	CcEmail     []string `json:"cc_emails"`
	Source      int      `json:"source"`
	//DueBy        string       `json:"due_by"`
	//FrDueBy      string       `json:"fr_due_by"`
	CustomFields CustomFields `json:"custom_fields"`
}

type FSResponse struct {
	Ticket struct {
		CcEmails        []string      `json:"cc_emails"`
		FwdEmails       []interface{} `json:"fwd_emails"`
		ReplyCcEmails   []string      `json:"reply_cc_emails"`
		FrEscalated     bool          `json:"fr_escalated"`
		Spam            bool          `json:"spam"`
		EmailConfigID   interface{}   `json:"email_config_id"`
		GroupID         interface{}   `json:"group_id"`
		Priority        int           `json:"priority"`
		RequesterID     int           `json:"requester_id"`
		RequestedForID  int           `json:"requested_for_id"`
		ResponderID     interface{}   `json:"responder_id"`
		Source          int           `json:"source"`
		Status          int           `json:"status"`
		Subject         string        `json:"subject"`
		ToEmails        interface{}   `json:"to_emails"`
		DepartmentID    interface{}   `json:"department_id"`
		ID              int           `json:"id"`
		Type            string        `json:"type"`
		DueBy           time.Time     `json:"due_by"`
		FrDueBy         time.Time     `json:"fr_due_by"`
		IsEscalated     bool          `json:"is_escalated"`
		Description     string        `json:"description"`
		DescriptionText string        `json:"description_text"`
		Category        interface{}   `json:"category"`
		SubCategory     interface{}   `json:"sub_category"`
		ItemCategory    interface{}   `json:"item_category"`
		CustomFields    struct {
			AutoCheckbox interface{} `json:"auto_checkbox"`
		} `json:"custom_fields"`
		CreatedAt   time.Time     `json:"created_at"`
		UpdatedAt   time.Time     `json:"updated_at"`
		Tags        []interface{} `json:"tags"`
		Attachments []interface{} `json:"attachments"`
	} `json:"ticket"`
}

type CustomFields struct {
	CaseId    string `json:"case_id"`
	DisplayId string `json:"display_id"`
}

type ReplyFS struct {
	Body string `json:"body"`
}
