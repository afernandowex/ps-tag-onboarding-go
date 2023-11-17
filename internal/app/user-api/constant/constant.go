package constant

const (
	NoUser                 = -1
	HttpPort               = 8080
	ErrorEmailRequired     = "Email is required"
	ErrorEmailFormat       = "Email is not valid"
	ErrorAgeMinimum        = "Age must be at least 18"
	ErrorNameRequired      = "Name is required"
	ErrorNameAlreadyExists = "Name already exists"
	ErrorInvalidUserID     = "Invalid user ID. Should be a number"
	ErrorUserNotFound      = "User not found"
	ErrorInvalidUserObject = "User object cannot be unmarshalled"
)
