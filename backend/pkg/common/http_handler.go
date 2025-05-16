package common

import (
	echo "github.com/labstack/echo/v4"
)

type Handler interface {
	Handlers() map[HTTPMethod]echo.HandlerFunc
}

type BaseHandler struct{}

func (h *BaseHandler) Handlers() map[HTTPMethod]echo.HandlerFunc {
	return nil
}

func NewHandler(handlers map[HTTPMethod]echo.HandlerFunc) Handler {
	return &handler{handlers: handlers}
}

type handler struct {
	handlers map[HTTPMethod]echo.HandlerFunc
}

func (h *handler) Handlers() map[HTTPMethod]echo.HandlerFunc {
	return h.handlers
}

type ErrorResponse struct {
	Error string `json:"error"`
}
