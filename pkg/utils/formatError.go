package utils

import "github.com/labstack/echo/v4"

type JSONErrorMessage struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func FormatError(c echo.Context, message string, status int) error {
	return c.JSON(status, &JSONErrorMessage{Message: message, Status: status})
}
