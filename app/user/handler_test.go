package user

import (
	"bytes"
	"errors"
	"go-restapi/app"
	"go-restapi/logger"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestCases struct {
	name           string
	url            string
	method         string
	reqBody        string
	expectedStatus int
	expectedBody   string
}

var CreateUserSuccessCases = []TestCases{
	{
		name:           "CreateUser: Should return success message",
		url:            "/users",
		method:         "POST",
		reqBody:        `{"username":"test","password":"password","firstname":"test","lastname":"test"}`,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":""}`,
	},
}

var CreateUserFailCases = []TestCases{
	{
		name:           "CreateUser: Should return error (Bind error)",
		url:            "/users",
		method:         "POST",
		reqBody:        `{"username":"test","password":"test","firstname":"test","lastname":"test"`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
	{
		name:           "CreateUser: Should return error (Validate error)",
		url:            "/users",
		method:         "POST",
		reqBody:        `{"username":"test","password":"test","firstname":"test","lastname":"test"}`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
	{
		name:           "CreateUser: Should return error (Service error)",
		url:            "/users",
		method:         "POST",
		reqBody:        `{"username":"test","password":"password","firstname":"test","lastname":"test"}`,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"The server encountered an unexpected condition which prevented it from fulfilling the request."}`,
	},
}

var CreateUserBadReqFailCases = []TestCases{
	{
		name:           "CreateUser: Should return error (BadRequest)",
		url:            "/users",
		method:         "POST",
		reqBody:        `{"username":"test","password":"password","firstname":"test","lastname":"test"}`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
}

func TestCreateUserHandler(t *testing.T) {
	serviceSuccess := &mockUserService{}
	serviceSuccess.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	serviceFailStore := &mockUserService{}
	serviceFailStore.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("error"))

	serviceFailBadRequest := &mockUserService{}
	serviceFailBadRequest.On("CreateUser", mock.Anything, mock.Anything).Return(ErrUsernameAlreadyExists)

	t.Run("Success Case", RunTest(serviceSuccess, CreateUserSuccessCases))
	t.Run("Fail Case Store", RunTest(serviceFailStore, CreateUserFailCases))
	t.Run("Fail Case BadRequest", RunTest(serviceFailBadRequest, CreateUserBadReqFailCases))

}

// ----------------------------

var GetListUserSuccessCasesEmptyData = []TestCases{
	{
		name:           "GetListUser: Should return success message (Empty data)",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","data":[]}`,
	},
}

var GetListUserSuccessCases = []TestCases{
	{
		name:           "GetListUser: Should return success message",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","data":[{"id":0,"username":"test","firstname":"test","lastname":"test","status":"Active"}]}`,
	},
}

var GetListUserFailCases = []TestCases{
	{
		name:           "GetListUser: Should return error (Service error)",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"The server encountered an unexpected condition which prevented it from fulfilling the request."}`,
	},
}

func TestGetListUserHandler(t *testing.T) {
	mockData := []GetListUserResponse{
		{
			ID:        0,
			Username:  "test",
			FirstName: "test",
			LastName:  "test",
			Status:    "Active",
		},
	}
	serviceSuccess := &mockUserService{}
	serviceSuccess.On("GetListUser", mock.Anything).Return(mockData, nil)
	t.Run("Success Case", RunTest(serviceSuccess, GetListUserSuccessCases))

	serviceSuccessEmptyData := &mockUserService{}
	serviceSuccessEmptyData.On("GetListUser", mock.Anything).Return([]GetListUserResponse{}, nil)
	t.Run("Success Case Empty data", RunTest(serviceSuccessEmptyData, GetListUserSuccessCasesEmptyData))

	serviceFail := &mockUserService{}
	serviceFail.On("GetListUser", mock.Anything).Return([]GetListUserResponse{}, errors.New("error"))
	t.Run("Fail Case", RunTest(serviceFail, GetListUserFailCases))
}

func RunTest(service UserService, testCases []TestCases) func(t *testing.T) {
	return func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		h := NewUserHandler(service)
		r.POST("/users", toGinHandlerFunc(h.CreateUser))
		r.GET("/users", toGinHandlerFunc(h.GetListUser))

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body := bytes.NewBufferString(tc.reqBody)

				req := httptest.NewRequest(tc.method, tc.url, body)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)

				assert.Equal(t, tc.expectedStatus, rec.Code)
				assert.Equal(t, tc.expectedBody, rec.Body.String())
			})
		}
	}
}

func toGinHandlerFunc(f func(ctx app.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.New()
		ctx := app.NewContext(c, l.Handler())
		f(ctx)
	}
}
