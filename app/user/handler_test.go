package user

import (
	"bytes"
	"errors"
	"go-restapi/app"
	"go-restapi/logger"
	"go-restapi/utils"
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
	mockUtils := &mockUtils{}
	serviceSuccess := &mockUserService{}
	serviceSuccess.On("CreateUser", mock.Anything, mock.Anything).Return(nil)

	serviceFailStore := &mockUserService{}
	serviceFailStore.On("CreateUser", mock.Anything, mock.Anything).Return(errors.New("error"))

	serviceFailBadRequest := &mockUserService{}
	serviceFailBadRequest.On("CreateUser", mock.Anything, mock.Anything).Return(ErrUsernameAlreadyExists)

	t.Run("Success Case", RunTest(serviceSuccess, mockUtils, CreateUserSuccessCases))
	t.Run("Fail Case Store", RunTest(serviceFailStore, mockUtils, CreateUserFailCases))
	t.Run("Fail Case BadRequest", RunTest(serviceFailBadRequest, mockUtils, CreateUserBadReqFailCases))

}

// ----------------------------

var GetListUserSuccessCasesEmptyData = []TestCases{
	{
		name:           "GetListUser: Should return success message (Empty data)",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","currentPage":1,"data":[]}`,
	},
}

var GetListUserSuccessCases = []TestCases{
	{
		name:           "GetListUser: Should return success message",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","currentRecord":1,"currentPage":1,"totalRecord":1,"totalPage":1,"data":[{"id":0,"username":"test","firstname":"test","lastname":"test","status":"Active"}]}`,
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

var GetCountListUserFailCases = []TestCases{
	{
		name:           "GetCountListUser: Should return error (Service error)",
		url:            "/users",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"The server encountered an unexpected condition which prevented it from fulfilling the request."}`,
	},
}

func TestGetListUserHandler(t *testing.T) {
	RunGetListUserSuccessCase(t)
	RunGetListUserSuccessEmptyDataCase(t)
	RunGetListUserFailCase(t)
	RunGetListUserFailCaseCountListUser(t)
}

func RunGetListUserSuccessCase(t *testing.T) {
	mockData := []GetListUserResponse{
		{
			ID:        0,
			Username:  "test",
			FirstName: "test",
			LastName:  "test",
			Status:    "Active",
		},
	}

	mockUtils := &mockUtils{}
	mockUtils.On("GetPage", mock.Anything).Return(1, nil)
	mockUtils.On("GetPageSize", mock.Anything).Return(10, nil)
	mockUtils.On("GetTotalPage", mock.Anything, mock.Anything).Return(1)

	serviceSuccess := &mockUserService{}
	serviceSuccess.On("GetListUser", mock.Anything, mock.Anything, mock.Anything).Return(mockData, nil)
	serviceSuccess.On("CountListUser", mock.Anything).Return(1, nil)
	t.Run("Success Case", RunTest(serviceSuccess, mockUtils, GetListUserSuccessCases))
}

func RunGetListUserSuccessEmptyDataCase(t *testing.T) {
	mockUtils := &mockUtils{}
	mockUtils.On("GetPage", mock.Anything).Return(1, nil)
	mockUtils.On("GetPageSize", mock.Anything).Return(10, nil)
	mockUtils.On("GetTotalPage", mock.Anything, mock.Anything).Return(0)

	serviceSuccessEmptyData := &mockUserService{}
	serviceSuccessEmptyData.On("GetListUser", mock.Anything, mock.Anything, mock.Anything).Return([]GetListUserResponse{}, nil)
	serviceSuccessEmptyData.On("CountListUser", mock.Anything).Return(0, nil)
	t.Run("Success Case Empty data", RunTest(serviceSuccessEmptyData, mockUtils, GetListUserSuccessCasesEmptyData))
}

func RunGetListUserFailCase(t *testing.T) {
	mockUtils := &mockUtils{}
	mockUtils.On("GetPage", mock.Anything).Return(1, nil)
	mockUtils.On("GetPageSize", mock.Anything).Return(10, nil)
	mockUtils.On("GetTotalPage", mock.Anything, mock.Anything).Return(0)

	serviceFail := &mockUserService{}
	serviceFail.On("GetListUser", mock.Anything, mock.Anything, mock.Anything).Return([]GetListUserResponse{}, errors.New("error"))
	t.Run("Fail Case GetListUser", RunTest(serviceFail, mockUtils, GetListUserFailCases))
}

func RunGetListUserFailCaseCountListUser(t *testing.T) {
	mockUtils := &mockUtils{}
	mockUtils.On("GetPage", mock.Anything).Return(1, nil)
	mockUtils.On("GetPageSize", mock.Anything).Return(10, nil)
	mockUtils.On("GetTotalPage", mock.Anything, mock.Anything).Return(0)

	serviceFail := &mockUserService{}
	serviceFail.On("GetListUser", mock.Anything, mock.Anything, mock.Anything).Return([]GetListUserResponse{}, nil)
	serviceFail.On("CountListUser", mock.Anything).Return(0, errors.New("error"))
	t.Run("Fail Case CountListUser", RunTest(serviceFail, mockUtils, GetCountListUserFailCases))
}

// ----------------------------

var GetUserByIDSuccessCases = []TestCases{
	{
		name:           "GetUserByID: Should return success message",
		url:            "/users/0",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","data":{"id":0,"username":"test","firstname":"test","lastname":"test","status":"Active"}}`,
	},
}

var GetUserByIDFailCases = []TestCases{
	{
		name:           "GetUserByID: Should return error (ID invalid)",
		url:            "/users/abc",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
	{
		name:           "GetUserByID: Should return error (Service error)",
		url:            "/users/0",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"The server encountered an unexpected condition which prevented it from fulfilling the request."}`,
	},
}

var GetUserByIDNotFoundCases = []TestCases{
	{
		name:           "GetUserByID: Should return error (Service error Not found)",
		url:            "/users/0",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 404,
		expectedBody:   `{"status":"ERROR","message":"The requested resource could not be found but may be available in the future."}`,
	},
}

func TestGetUserByIDHandler(t *testing.T) {
	RunGetUserByIDHandlerSuccessCase(t)
	RunGetUserByIDHandlerFailCase(t)
	RunGetUserByIDHandlerNotFoundCase(t)
}

func RunGetUserByIDHandlerSuccessCase(t *testing.T) {
	mockData := &GetUserResponse{
		ID:        0,
		Username:  "test",
		FirstName: "test",
		LastName:  "test",
		Status:    "Active",
	}

	mockUtils := &mockUtils{}
	serviceSuccess := &mockUserService{}
	serviceSuccess.On("GetUserByID", mock.Anything, 0).Return(mockData, nil)
	t.Run("Success Case", RunTest(serviceSuccess, mockUtils, GetUserByIDSuccessCases))
}

func RunGetUserByIDHandlerFailCase(t *testing.T) {
	mockUtils := &mockUtils{}
	serviceFail := &mockUserService{}
	serviceFail.On("GetUserByID", mock.Anything, 0).Return(nil, errors.New("error"))
	t.Run("Fail Case", RunTest(serviceFail, mockUtils, GetUserByIDFailCases))
}

func RunGetUserByIDHandlerNotFoundCase(t *testing.T) {
	mockUtils := &mockUtils{}
	serviceFail := &mockUserService{}
	serviceFail.On("GetUserByID", mock.Anything, 0).Return(nil, ErrUserNotFound)
	t.Run("Fail Case Not found", RunTest(serviceFail, mockUtils, GetUserByIDNotFoundCases))
}

func RunTest(service UserService, utils utils.Utils, testCases []TestCases) func(t *testing.T) {
	return func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		h := NewUserHandler(service, utils)
		r.POST("/users", toGinHandlerFunc(h.CreateUser))
		r.GET("/users", toGinHandlerFunc(h.GetListUser))
		r.GET("/users/:id", toGinHandlerFunc(h.GetUserByID))

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
