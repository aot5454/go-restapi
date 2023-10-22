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
	Paging
	Data any `json:"data,omitempty"`
}

type Paging struct {
	CurrentRecord int `json:"currentRecord,omitempty"`
	CurrentPage   int `json:"currentPage,omitempty"`
	TotalRecord   int `json:"totalRecord,omitempty"`
	TotalPage     int `json:"totalPage,omitempty"`
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

type TokenData struct {
	UserID    int    `json:"userId"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Role      string `json:"role"`
}

type Config struct {
	Env      string   `mapstructure:"env"`
	Server   Server   `mapstructure:"server"`
	Database Database `mapstructure:"db"`
}

type Server struct {
	Port string `mapstructure:"port"`
}

type Database struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}
