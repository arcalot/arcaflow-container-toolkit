package docker

type ErrorLine struct {
	Error       string      `json:"error"`
	ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

type StreamLine struct {
	Stream string `json:"stream"`
}

type MalformedErrorDetails struct {
	messge string
}

func (err MalformedErrorDetails) Error() string {
	return err.messge
}

func NewMalformedErrorDetails(msg string) *MalformedErrorDetails {
	return &MalformedErrorDetails{
		messge: msg,
	}
}

type ErrorDetails struct {
	message string
}

func (err ErrorDetails) Error() string {
	return err.message
}

func NewErrorDetails(msg string) *ErrorDetails {
	return &ErrorDetails{
		message: msg,
	}
}
