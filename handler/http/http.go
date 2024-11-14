package handlerhttp

import (
	"context"
	"moviefestival/model"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/labstack/echo/v4"
)

type HandlerHTTP struct {
	authUsecase authUsecase
}

type authUsecase interface {
	Register(ctx context.Context, user model.User) error
}

func NewHandlerHTTP(authUsecase authUsecase) *HandlerHTTP {
	return &HandlerHTTP{
		authUsecase: authUsecase,
	}
}

func (h *HandlerHTTP) InitRoute(e *echo.Echo) {
	e.POST("/auth/register", h.Register)
}
func (h *HandlerHTTP) Register(c echo.Context) error {
	var user model.User
	c.Bind(&user)
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Name, validation.Required, validation.Length(3, 30)),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 30)),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = h.authUsecase.Register(c.Request().Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
	})
}
