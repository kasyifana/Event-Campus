package request

// SubmitWhitelistRequest represents whitelist submission request
type SubmitWhitelistRequest struct {
	OrganizationName string `form:"organization_name" binding:"required"`
}

// ReviewWhitelistRequest represents whitelist review request
type ReviewWhitelistRequest struct {
	Approved   bool    `json:"approved" binding:"required"`
	AdminNotes *string `json:"admin_notes,omitempty"`
}
