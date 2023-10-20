package user

import (
	"go-restapi/app"
	"go-restapi/utils"
	"strconv"
)

type UserHandler interface {
	CreateUser(ctx app.Context)
	GetListUser(ctx app.Context)
	GetUserByID(ctx app.Context)
}

type userHandler struct {
	userSvc UserService
	utils   utils.Utils
}

func NewUserHandler(userService UserService, utils utils.Utils) UserHandler {
	return &userHandler{
		userSvc: userService,
		utils:   utils,
	}
}

func (h *userHandler) CreateUser(ctx app.Context) {
	var user = CreateUserRequest{}

	if err := ctx.Bind(&user); err != nil {
		ctx.BadRequest(err)
		return
	}

	if _, err := ctx.Validate(&user); err != nil {
		ctx.BadRequest(err)
		return
	}

	if err := h.userSvc.CreateUser(ctx, user); err != nil {
		if err == ErrUsernameAlreadyExists {
			ctx.BadRequest(err)
			return
		}

		ctx.StoreError(err)
		return
	}

	ctx.OK(nil)
}

func (h *userHandler) GetListUser(ctx app.Context) {
	page, _ := h.utils.GetPage(ctx)
	pageSize, _ := h.utils.GetPageSize(ctx)

	users, err := h.userSvc.GetListUser(ctx, page, pageSize)
	if err != nil {
		ctx.StoreError(err)
		return
	}

	totalRecord, err := h.userSvc.CountListUser(ctx)
	if err != nil {
		ctx.StoreError(err)
		return
	}

	paging := app.Paging{
		CurrentRecord: len(users),
		CurrentPage:   page,
		TotalRecord:   totalRecord,
		TotalPage:     h.utils.GetTotalPage(totalRecord, pageSize),
	}

	ctx.OKWithPaging(users, paging)
}

func (h *userHandler) GetUserByID(ctx app.Context) {
	paramId := ctx.GetParam("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		ctx.BadRequest(err)
		return
	}

	user, err := h.userSvc.GetUserByID(ctx, id)
	if err != nil {
		if err == ErrUserNotFound {
			ctx.NotFound()
			return
		}
		ctx.StoreError(err)
		return
	}

	ctx.OK(user)
}
