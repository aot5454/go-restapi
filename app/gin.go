package app

import (
	"go-restapi/logger"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Context interface {
	Bind(any) error
	Validate(any) ([]ErrorField, error)
	OK(any)
	OKWithPaging(any, Paging)
	BadRequest(err error)
	StoreError(err error)
	InternalServerError(err error)
	Conflict(err error)
	NotFound()
	GetAllHeader() http.Header
	GetHeader(string) string
	GetQuery(string) string
	GetParam(string) string
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

func (c *context) GetQuery(key string) string {
	return c.Context.Query(key)
}

func (c *context) GetParam(key string) string {
	return c.Context.Param(key)
}

func (c *context) GetAllHeader() http.Header {
	return c.Request.Header
}

func (c *context) GetHeader(key string) string {
	return c.Context.GetHeader(key)
}

func (c *context) OK(data any) { // 200
	c.Context.JSON(http.StatusOK, Response{
		Status: Success,
		Data:   data,
	})
}

func (c *context) OKWithPaging(data any, paging Paging) { // 200
	c.Context.JSON(http.StatusOK, Response{
		Status: Success,
		Paging: paging,
		Data:   data,
	})
}

func (c *context) BadRequest(err error) { // 400
	logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusBadRequest, Response{
		Status:  Fail,
		Message: BadRequestMsg,
	})
}

func (c *context) NotFound() { // 404
	// logger.AppErrorf(c.logHandler, "%s", err)
	c.Context.JSON(http.StatusNotFound, Response{
		Status:  Fail,
		Message: NotFoundMsg,
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
		transationId := c.Request.Header.Get("transaction-id")
		if transationId == "" {
			transationId = uuid.NewString()
			c.Request.Header.Set("transaction-id", transationId)
		}
		handler(NewContext(c, logger.Handler().WithAttrs([]slog.Attr{slog.String("transaction-id", transationId)})))
		// handler(NewContext(c, logger.With(zap.String("transaction-id", c.Request.Header.Get("transaction-id")))))
	}
}

type Router struct {
	*gin.Engine
	logger *slog.Logger
}

func NewRouter(logger *slog.Logger, conf Config) *Router {
	if conf.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
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

func (r *Router) NoRoute() {
	handler := func(c Context) { c.NotFound() }
	r.Engine.NoRoute(NewGinHandler(handler, r.logger))
}

type RouterGroup struct {
	*gin.RouterGroup
	logger *slog.Logger
}

func (r *Router) Group(path string) *RouterGroup {
	return &RouterGroup{
		RouterGroup: r.Engine.Group(path),
		logger:      r.logger,
	}
}

func (rg *RouterGroup) GET(path string, handler func(Context)) {
	rg.RouterGroup.GET(path, NewGinHandler(handler, rg.logger))
}

func (rg *RouterGroup) POST(path string, handler func(Context)) {
	rg.RouterGroup.POST(path, NewGinHandler(handler, rg.logger))
}

func (rg *RouterGroup) PUT(path string, handler func(Context)) {
	rg.RouterGroup.PUT(path, NewGinHandler(handler, rg.logger))
}

func (rg *RouterGroup) PATCH(path string, handler func(Context)) {
	rg.RouterGroup.PATCH(path, NewGinHandler(handler, rg.logger))
}

func (rg *RouterGroup) DELETE(path string, handler func(Context)) {
	rg.RouterGroup.DELETE(path, NewGinHandler(handler, rg.logger))
}
