package device42

const (
	// ErrorEmptyCredentials returns credential error
	ErrorEmptyCredentials = "invalid credentials: username & password must be specified"
	// ErrorEmptyUsername returns username error
	ErrorEmptyUsername = "invalid credentials: username must be specified"
	// ErrorEmptyPassword returns password error
	ErrorEmptyPassword = "invalid credentials: password must be specified"
	// ErrorEmptyHost returns empty host error
	ErrorEmptyHost = "invalid host: you must supply the device42 host"
)
