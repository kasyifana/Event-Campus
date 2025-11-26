package request

// MarkAttendanceRequest represents attendance marking request
type MarkAttendanceRequest struct {
	UserID string  `json:"user_id" binding:"required"`
	Notes  *string `json:"notes,omitempty"`
}

// BulkMarkAttendanceRequest represents bulk attendance marking request
type BulkMarkAttendanceRequest struct {
	UserIDs []string `json:"user_ids" binding:"required"`
}
