package user

import "go-restapi/app"

type UserHandler interface {
	CreateUser(ctx app.Context)
}

type userHandler struct {
	userSvc UserService
}

func NewUserHandler(userService UserService) UserHandler {
	return &userHandler{
		userSvc: userService,
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
		ctx.StoreError(err)
		return
	}

	ctx.OK(nil)
}
