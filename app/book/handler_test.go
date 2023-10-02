package book

import (
	"bytes"
	"errors"
	"go-restapi/app"
	"go-restapi/logger"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type bookServiceMockSuccess struct {
	BookService
}

func (m *bookServiceMockSuccess) GetAllBook() ([]Book, error) {
	return []Book{
		{
			ID:     1,
			Title:  "test",
			Author: "test",
		},
	}, nil
}

func (m *bookServiceMockSuccess) CreateBook(book BookRequest) error {
	return nil
}

func toGinHandlerFunc(f func(ctx app.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		l := logger.New()
		ctx := app.NewContext(c, l.Handler())
		f(ctx)
	}
}

var testSuccessCases = []struct {
	name           string
	url            string
	method         string
	reqBody        string
	expectedStatus int
	expectedBody   string
}{
	{
		name:           "GetAllBook: Should return array of book and success message",
		url:            "/books",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: http.StatusOK,
		expectedBody:   `{"status":"SUCCESS","message":"","data":[{"id":1,"title":"test","author":"test"}]}`,
	},
	{
		name:           "CreateBook: Should return success message",
		url:            "/books",
		method:         "POST",
		reqBody:        `{"title":"test","author":"test"}`,
		expectedStatus: http.StatusOK,
		expectedBody:   `{"status":"SUCCESS","message":""}`,
	},
}

func TestBookHandlerSuccessCase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	bookHandler := NewHandler(&bookServiceMockSuccess{})
	r.GET("/books", toGinHandlerFunc(bookHandler.GetAllBook))
	r.POST("/books", toGinHandlerFunc(bookHandler.CreateBook))

	for _, tc := range testSuccessCases {
		t.Run(tc.name, func(t *testing.T) {
			body := bytes.NewBufferString(tc.reqBody)

			req := httptest.NewRequest(tc.method, tc.url, body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d but got %d", tc.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tc.expectedBody {
				t.Errorf("expected body %s but got %s", tc.expectedBody, rec.Body.String())
			}
		})
	}
}

// --------------------------------------------------------------------------------------------

type bookServiceMockError struct {
	BookService
}

func (m *bookServiceMockError) GetAllBook() ([]Book, error) {
	return []Book{}, errors.New("error")
}

func (m *bookServiceMockError) CreateBook(book BookRequest) error {
	return errors.New("error")
}

var testErrorCases = []struct {
	name           string
	url            string
	method         string
	reqBody        string
	expectedStatus int
	expectedBody   string
}{
	{
		name:           "GetAllBook: Should return store error message",
		url:            "/books",
		method:         "GET",
		reqBody:        ``,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"` + app.StoreErrorMsg + `"}`,
	},
	{
		name:           "CreateBook: Should return BadRequest error message (Bind error)",
		url:            "/books",
		method:         "POST",
		reqBody:        `{"title":"test"`,
		expectedStatus: http.StatusBadRequest,
		expectedBody:   `{"status":"ERROR","message":"` + app.BadRequestMsg + `"}`,
	},
	{
		name:           "CreateBook: Should return BadRequest error message (Validate error)",
		url:            "/books",
		method:         "POST",
		reqBody:        `{"title":"test"}`,
		expectedStatus: http.StatusBadRequest,
		expectedBody:   `{"status":"ERROR","message":"` + app.BadRequestMsg + `"}`,
	},
	{
		name:           "CreateBook: Should return error message (Store error)",
		url:            "/books",
		method:         "POST",
		reqBody:        `{"title":"test","author":"test"}`,
		expectedStatus: 507,
		expectedBody:   `{"status":"ERROR","message":"` + app.StoreErrorMsg + `"}`,
	},
}

func TestBookHandlerErrorCase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	bookHandler := NewHandler(&bookServiceMockError{})
	r.GET("/books", toGinHandlerFunc(bookHandler.GetAllBook))
	r.POST("/books", toGinHandlerFunc(bookHandler.CreateBook))

	for _, tc := range testErrorCases {
		t.Run(tc.name, func(t *testing.T) {
			body := bytes.NewBufferString(tc.reqBody)

			req := httptest.NewRequest(tc.method, tc.url, body)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d but got %d", tc.expectedStatus, rec.Code)
			}

			if rec.Body.String() != tc.expectedBody {
				t.Errorf("expected body %s but got %s", tc.expectedBody, rec.Body.String())
			}
		})
	}
}
