package errs

var (
	ErrInternal = New("internal_error", 500, "An internal server error occurred")
	ErrInvalidJSON = New("invalid_json", 400, "The JSON payload is invalid")
	// ErrWorkflowNotFound = New("workflow_not_found", 404, "The specified workflow was not found")
)
