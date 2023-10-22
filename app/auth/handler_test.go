package auth

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
	"github.com/stretchr/testify/suite"
)

type TestCase struct {
	name           string
	url            string
	method         string
	reqBody        string
	expectedStatus int
	expectedBody   string
}

var mockAuthResponseData = &AuthResponse{
	AccessToken:          "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc5Nzk2ODcsImZpcnN0bmFtZSI6IkFtaXlhIiwiaWF0IjoxNjk3OTU4MDg3LCJpc3MiOiJnby1yZXN0YXBpIiwibGFzdG5hbWUiOiJBcm1zdHJvbmciLCJyb2xlIjoiYWRtaW4iLCJzdWIiOjEzLCJ1c2VySUQiOjEzLCJ1c2VybmFtZSI6ImFkbWluIn0.RETaz4RWH1VHloJXv-NVrDY1VdgRTXbPY6dxaEXVlF2kiqsbHMdPkY8KT-wjV0k06Jwc3KtYgcSqrT4iWEa5Ej8JnoF3ag56OXkbnTH0cvdA9oTgltSgSUd1OlafUK3IPS-8XbFFSHbV-3oCN8tiUQMLAmF78DMBPB3H2B2JsegL5No0P-WUb3ZOzVEVyOXRs5sh57EXrFWQb-t-ZPRha6P-gskE0uyt2BQ3FmTNuk17yk2ssA_iGW21wdwsYAOmAbeEjaHgcyIQPOVkkMBNBvY6qlOKzjJYzYC3sScF0o-3EvTqd6UPpTPBhqY7R5mpNGD6NylE4sZI9K840aKs6w",
	AccessTokenExpireAt:  "2021-08-24 15:13:07",
	RefreshToken:         "fcd277b6-562c-49f6-8146-051bb339fb8c",
	RefreshTokenExpireAt: "2021-08-25 15:13:07",
}

var LoginSuccessCases = []TestCase{
	{
		name:           "Should return 200",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"admin","password":"password"}`,
		expectedStatus: 200,
		expectedBody:   `{"status":"SUCCESS","message":"","data":{"accessToken":"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTc5Nzk2ODcsImZpcnN0bmFtZSI6IkFtaXlhIiwiaWF0IjoxNjk3OTU4MDg3LCJpc3MiOiJnby1yZXN0YXBpIiwibGFzdG5hbWUiOiJBcm1zdHJvbmciLCJyb2xlIjoiYWRtaW4iLCJzdWIiOjEzLCJ1c2VySUQiOjEzLCJ1c2VybmFtZSI6ImFkbWluIn0.RETaz4RWH1VHloJXv-NVrDY1VdgRTXbPY6dxaEXVlF2kiqsbHMdPkY8KT-wjV0k06Jwc3KtYgcSqrT4iWEa5Ej8JnoF3ag56OXkbnTH0cvdA9oTgltSgSUd1OlafUK3IPS-8XbFFSHbV-3oCN8tiUQMLAmF78DMBPB3H2B2JsegL5No0P-WUb3ZOzVEVyOXRs5sh57EXrFWQb-t-ZPRha6P-gskE0uyt2BQ3FmTNuk17yk2ssA_iGW21wdwsYAOmAbeEjaHgcyIQPOVkkMBNBvY6qlOKzjJYzYC3sScF0o-3EvTqd6UPpTPBhqY7R5mpNGD6NylE4sZI9K840aKs6w","accessTokenExpireAt":"2021-08-24 15:13:07","refreshToken":"fcd277b6-562c-49f6-8146-051bb339fb8c","refreshTokenExpireAt":"2021-08-25 15:13:07"}}`,
	},
}

var LoginFailValidateCases = []TestCase{
	{
		name:           "Should return 400 when request body is invalid (Bind)",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"","password":"password"`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
	{
		name:           "Should return 400 when request body is invalid (Validate)",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"admin","password":""}`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
	{
		name:           "Should return 500 when service return error",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"admin","password":"password"}`,
		expectedStatus: 500,
		expectedBody:   `{"status":"ERROR","message":"The server encountered an unexpected condition which prevented it from fulfilling the request."}`,
	},
}

var LoginFailServiceNotFoundCases = []TestCase{
	{
		name:           "Should return 404 when user not found",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"admin","password":"password"}`,
		expectedStatus: 404,
		expectedBody:   `{"status":"ERROR","message":"The requested resource could not be found but may be available in the future."}`,
	},
}

var LoginFailServicePasswordNotMatchCases = []TestCase{
	{
		name:           "Should return 400 when password not match",
		url:            "/login",
		method:         "POST",
		reqBody:        `{"username":"admin","password":"password"}`,
		expectedStatus: 400,
		expectedBody:   `{"status":"ERROR","message":"Invalid request body, Please check your request body and try again!"}`,
	},
}

// -------------------------------------

type testHandlerSuite struct {
	suite.Suite
}

func (s *testHandlerSuite) SetupTest() {

}

func toGinHandlerFunc(f func(ctx app.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.New()
		ctx := app.NewContext(c, l.Handler())
		f(ctx)
	}
}

func RunTest(service AuthService, testCases []TestCase) func(t *testing.T) {
	return func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		r := gin.Default()

		handler := NewAuthHandler(service)
		r.POST("/login", toGinHandlerFunc(handler.Login))

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

func (s *testHandlerSuite) TestLoginHandler() {
	authSvc := &mockAuthService{}
	authSvc.On("Login", mock.Anything).Return(mockAuthResponseData, nil)
	s.T().Run("Success Case", RunTest(authSvc, LoginSuccessCases))

	authFailSvc := &mockAuthService{}
	authFailSvc.On("Login", mock.Anything).Return(nil, errors.New("error"))
	s.T().Run("Fail Validate Case", RunTest(authFailSvc, LoginFailValidateCases))

	authFailNotfoundSvc := &mockAuthService{}
	authFailNotfoundSvc.On("Login", mock.Anything).Return(nil, ErrUserNotFound)
	s.T().Run("Fail Not found Case", RunTest(authFailNotfoundSvc, LoginFailServiceNotFoundCases))

	authFailPasswordNotMatchSvc := &mockAuthService{}
	authFailPasswordNotMatchSvc.On("Login", mock.Anything).Return(nil, ErrPasswordNotMatch)
	s.T().Run("Fail Password not match Case", RunTest(authFailPasswordNotMatchSvc, LoginFailServicePasswordNotMatchCases))
}

func TestAuthHandler(t *testing.T) {
	suite.Run(t, new(testHandlerSuite))
}
