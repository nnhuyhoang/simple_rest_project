package consts

const (
	// AuthInfo is a key to get authentication info from gin context
	AuthInfo = "AuthInfo"

	// Role Code
	FieldTechnician        = "field_technician"
	FieldTechnicianManager = "field_technician_manager"

	//status
	StatusInProgress    = "in_progress"
	StatusCompleted     = "completed"
	StatusMalfunctioned = "malfunctioned"

	// Response Message
	SuccessMessage = "success"

	// Target Type
	TargetTypeInspection  = "inspection"
	TargetTypeIssue       = "issue"
	TargetPurchaseRequest = "purchase_request"

	// User Action
	ActionCreateInspection      = "create_inspection"
	ActionUpdateInspection      = "update_inspection"
	ActionCreateIssue           = "create_issue"
	ActionUpdateIssue           = "update_issue"
	ActionDeleteIssue           = "delete_issue"
	ActionCreatePurchaseRequest = "create_purchase_request"

	//Email service
	ExchangeNameSignUp = "email_signup"
	WorkerRoutingKey   = "worker"
	ManagerRoutingKey  = "manager"

	//Default
	DefaultTimezone = "Asia/Singapore"
)
