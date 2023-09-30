package app

import (
	"go-restapi/logger"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Context interface {
	Bind(v any) error
	Validate(v any) ([]ErrorField, error)
	OK(v any)
	BadRequest(err error)
	StoreError(err error)
	InternalServerError(err error)
	Conflict(err error)
}

type context struct {
	*gin.Context
	logHandler slog.Handler
	validator  *validator.Validate
}

func NewContext(c *gin.Context, logHandler slog.Handler) Context {
	validate := validator.New()
	return &context{
		Context:    c,
		logHandler: logHandler,
		validator:  validate,
	}
}

func (c *context) Bind(v any) error {
	return c.Context.ShouldBindJSON(v)
}

func (c *context) Validate(v any) ([]ErrorField, error) {
	if err := c.validator.Struct(v); err != nil {
		var fields []ErrorField
		for _, v := range err.(validator.ValidationErrors) {
			errField := ErrorField{
				Value: v.Param(),
				Field: v.Field(),
				Tag:   v.Tag(),
			}
			fields = append(fields, errField)
		}
		return fields, err
	}
	return nil, nil
}

func (c *context) OK(v any) { // 200
	c.Context.JSON(http.StatusOK, Response{
		Status: Success,
		Data:   v,
	})
}

func (c *context) BadRequest(err error) { // 400
	logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusBadRequest, Response{
		Status:  Fail,
		Message: BadRequestMsg,
	})
}

func (c *context) Conflict(err error) { // 409
	logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusConflict, Response{
		Status:  Fail,
		Message: ConflictMsg,
	})
}

func (c *context) StoreError(err error) { // 450
	logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusInsufficientStorage, Response{
		Status:  Fail,
		Message: StoreErrorMsg,
	})
}

func (c *context) InternalServerError(err error) { // 500
	logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusInternalServerError, Response{
		Status:  Fail,
		Message: InternalServerErrorMsg,
	})
}

func NewGinHandler(handler func(Context), logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewContext(c, logger.Handler().WithAttrs([]slog.Attr{slog.String("transaction-id", c.Request.Header.Get("transaction-id"))})))
		// handler(NewContext(c, logger.With(zap.String("transaction-id", c.Request.Header.Get("transaction-id")))))
	}
}

type Router struct {
	*gin.Engine
	logger *slog.Logger
}

func NewRouter(logger *slog.Logger) *Router {
	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"X-Requested-With", "Authorization", "Origin", "Content-Length", "Content-Type", "TransactionID"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))

	return &Router{Engine: r, logger: logger}
}

func (r *Router) GET(path string, handler func(Context)) {
	r.Engine.GET(path, NewGinHandler(handler, r.logger))
}

func (r *Router) POST(path string, handler func(Context)) {
	r.Engine.POST(path, NewGinHandler(handler, r.logger))
}

func (r *Router) PUT(path string, handler func(Context)) {
	r.Engine.PUT(path, NewGinHandler(handler, r.logger))
}

func (r *Router) PATCH(path string, handler func(Context)) {
	r.Engine.PATCH(path, NewGinHandler(handler, r.logger))
}

func (r *Router) DELETE(path string, handler func(Context)) {
	r.Engine.DELETE(path, NewGinHandler(handler, r.logger))
}
