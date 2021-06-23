package device42

const (
	ErrorEmptyCredentials = "invalid credentials: username & password must be specified"
	ErrorEmptyUsername    = "invalid credentials: username must be specified"
	ErrorEmptyPassword    = "invalid credentials: password must be specified"
	ErrorEmptyBaseUrl     = "invalide url: you must supply a base URL"

	// HTTP errors
	ErrorHttpBadRequest          = "bad request (a validation exception has occurred)"
	ErrorHttpUnauthorized        = "unauthorized (your credentials suck)"
	ErrorHttpForbidden           = "forbidden (the resource requested is hidden)"
	ErrorHttpNotFound            = "not found (the specified resource could not be found)"
	ErrorHttpMethodNotAllowed    = "method not allowed (you tried to access a resource with an invalid method)"
	ErrorHttpGone                = "gone (the resource requested has been removed from our servers)"
	ErrorHttpInternalServerError = "internal server error (some parameter missing or issue with the server. check with returned “msg” from the call)"
	ErrorHttpServiceUnavaliable  = "service unavailable (please check if your device42 instance is working normally)"
)
