package carpentry

import "fmt"

type ErrorFilepathAbsolute struct {
	message string
}

func (e ErrorFilepathAbsolute) Error() string {
	return e.message
}

func NewErrorfFilepathAbsolute(e error) *ErrorFilepathAbsolute {
	return &ErrorFilepathAbsolute{
		message: fmt.Sprintf("invalid absolute path to project (%s)", e),
	}
}

type ErrorCEC struct {
	message string
}

func (e ErrorCEC) Error() string {
	return e.message
}

func NewErrorfCEC(e error) *ErrorCEC {
	return &ErrorCEC{
		message: fmt.Sprintf("invalid container engine client %s", e),
	}
}

type ErrorCarpenterConfig struct {
	message string
}

func (e ErrorCarpenterConfig) Error() string {
	return e.message
}

func NewErrorfCarpenterConfig(e error) *ErrorCarpenterConfig {
	return &ErrorCarpenterConfig{
		message: fmt.Sprintf("error in carpentry configuration file %s", e),
	}
}
