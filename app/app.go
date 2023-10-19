package app

type status string

const (
	Success status = "SUCCESS"
	Fail    status = "ERROR"
)

const (
	BadRequestMsg          string = "Invalid request body, Please check your request body and try again!"
	NotFoundMsg            string = "The requested resource could not be found but may be available in the future."
	ConflictMsg            string = "The request could not be completed due to a conflict with the current state of the target resource."
	StoreErrorMsg          string = "The server encountered an unexpected condition which prevented it from fulfilling the request."
	InternalServerErrorMsg string = "The server encountered an unexpected condition which prevented it from fulfilling the request."
)

type Response struct {
	Status  status `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorField struct {
	Field string `json:"field"`
	Value any    `json:"value"`
	Tag   string `json:"tag"`
}

type Error Response

func (err *Error) Error() string {
	return err.Message
}
