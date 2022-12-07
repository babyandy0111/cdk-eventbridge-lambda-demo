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

type FSTickets struct {
	Tickets []struct {
		Subject         string        `json:"subject"`
		GroupID         int           `json:"group_id"`
		DepartmentID    int           `json:"department_id"`
		Category        string        `json:"category"`
		SubCategory     string        `json:"sub_category"`
		ItemCategory    string        `json:"item_category"`
		RequesterID     int           `json:"requester_id"`
		ResponderID     int           `json:"responder_id"`
		DueBy           time.Time     `json:"due_by"`
		FrEscalated     bool          `json:"fr_escalated"`
		Deleted         bool          `json:"deleted"`
		Spam            bool          `json:"spam"`
		EmailConfigID   interface{}   `json:"email_config_id"`
		FwdEmails       []interface{} `json:"fwd_emails"`
		ReplyCcEmails   []interface{} `json:"reply_cc_emails"`
		CcEmails        []interface{} `json:"cc_emails"`
		IsEscalated     bool          `json:"is_escalated"`
		FrDueBy         time.Time     `json:"fr_due_by"`
		Priority        int           `json:"priority"`
		Source          int           `json:"source"`
		Status          int           `json:"status"`
		CreatedAt       time.Time     `json:"created_at"`
		UpdatedAt       time.Time     `json:"updated_at"`
		ToEmails        interface{}   `json:"to_emails"`
		ID              int           `json:"id"`
		Type            string        `json:"type"`
		Description     string        `json:"description"`
		DescriptionText string        `json:"description_text"`
		CustomFields    struct {
			Level3          interface{} `json:"level_3"`
			ReleaseCheckbox bool        `json:"release_checkbox"`
			CustomDropdown  string      `json:"custom_dropdown"`
			Level2          interface{} `json:"level_2"`
			CustomDate      string      `json:"custom_date"`
			CustomDateTime  time.Time   `json:"custom_date_time"`
			CustomCheckbox  bool        `json:"custom_checkbox"`
			SingleLineText  string      `json:"single_line_text"`
			ChangeCheckbox  bool        `json:"change_checkbox"`
			CustomParagraph string      `json:"custom_paragraph"`
			TfWithSection   string      `json:"tf_with_section"`
			CustomNumber    interface{} `json:"custom_number"`
			LastField       interface{} `json:"last_field"`
			SbDropdown      string      `json:"sb_dropdown"`
			SbSectionField  string      `json:"sb_section_field"`
			CustomDependent string      `json:"custom_dependent"`
		} `json:"custom_fields"`
	} `json:"tickets"`
}
