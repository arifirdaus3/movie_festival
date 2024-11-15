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
	Login(ctx context.Context, user model.User) (model.Token, error)
	Refresh(ctx context.Context, token string) (model.Token, error)
}

func NewHandlerHTTP(authUsecase authUsecase) *HandlerHTTP {
	return &HandlerHTTP{
		authUsecase: authUsecase,
	}
}

func (h *HandlerHTTP) InitRoute(e *echo.Echo) {
	e.POST("/auth/register", h.Register)
	e.POST("/auth/login", h.Login)
	e.POST("/auth/refresh", h.Refresh)
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
func (h *HandlerHTTP) Login(c echo.Context) error {
	email, password, ok := c.Request().BasicAuth()
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide email and password using basic auth")
	}
	user := model.User{
		Email:    email,
		Password: password,
	}
	err := validation.ValidateStruct(&user,
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Password, validation.Required, validation.Length(6, 30)),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := h.authUsecase.Login(c.Request().Context(), user)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = token.RefreshToken
	cookie.Expires = token.RefreshTokenExpiration
	c.SetCookie(cookie)
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data: map[string]interface{}{
			"accessToken": token.AcceessToken,
		},
	})
}

func (h *HandlerHTTP) Refresh(c echo.Context) error {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid refresh token")
	}
	err = validation.Validate(&refreshToken.Value, validation.Required.Error("Invalid refresh token"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	token, err := h.authUsecase.Refresh(c.Request().Context(), refreshToken.Value)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, model.Response{
		Status: model.ResponseStatusSuccess,
		Data: map[string]interface{}{
			"accessToken": token.AcceessToken,
		},
	})
}
