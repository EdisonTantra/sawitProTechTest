package sawithttp

import (
	"net/http"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func (h *Handler) MiddlewareLogging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.makeLogEntry(c).Info("incoming request")
		return next(c)
	}
}

func (h *Handler) MiddlewareError(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			report, ok := err.(*echo.HTTPError)
			if !ok {
				resp := generated.ErrorResponse{
					Message: http.StatusText(http.StatusInternalServerError),
					Details: []string{
						err.Error(),
					},
				}
				return c.JSON(http.StatusInternalServerError, resp)
			}

			msg := http.StatusText(report.Code)
			details := make([]string, 0)
			if errs, is := report.Message.(interface{ Unwrap() []error }); is {
				for _, err := range errs.Unwrap() {
					details = append(details, err.Error())
				}
			} else {
				details = append(details, report.Message.(string))
			}

			resp := generated.ErrorResponse{
				Message: msg,
				Details: details,
			}

			h.makeLogEntry(c).Error(report.Message)
			return c.JSON(report.Code, resp)
		}

		return nil
	}
}

func (h *Handler) makeLogEntry(c echo.Context) *log.Entry {
	const timeFormat = "2006-01-02 15:04:05"
	if c == nil {
		return h.logger.WithFields(log.Fields{
			"at": time.Now().Format(timeFormat),
		})
	}

	return h.logger.WithFields(log.Fields{
		"at":     time.Now().Format(timeFormat),
		"method": c.Request().Method,
		"uri":    c.Request().URL.String(),
		"ip":     c.Request().RemoteAddr,
	})
}
