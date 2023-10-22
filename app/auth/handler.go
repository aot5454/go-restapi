package auth

import (
	"errors"
	"go-restapi/app"
)

type AuthHandler interface {
	Login(ctx app.Context)
	// Logout(ctx app.Context)
	// RefreshToken(ctx app.Context)
}

type authHandler struct {
	authSvc AuthService
}

func NewAuthHandler(authSvc AuthService) AuthHandler {
	return &authHandler{
		authSvc: authSvc,
	}
}

func (h *authHandler) Login(ctx app.Context) {
	var req AuthRequest
	if err := ctx.Bind(&req); err != nil {
		ctx.BadRequest(err)
		return
	}

	if _, err := ctx.Validate(&req); err != nil {
		ctx.BadRequest(err)
		return
	}

	res, err := h.authSvc.Login(req)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			ctx.NotFound()
			return
		}
		if errors.Is(err, ErrPasswordNotMatch) {
			ctx.BadRequest(err)
			return
		}
		ctx.InternalServerError(err)
		return
	}

	ctx.OK(res)
}

// func (h *authHandler) Logout(ctx app.Context) {
// 	ctx.OK(nil)
// }

// func (h *authHandler) RefreshToken(ctx app.Context) {
// 	ctx.OK(nil)
// }
