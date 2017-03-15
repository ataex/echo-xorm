package version

import (
	"net/http"
	"time"

	"github.com/labstack/echo"

	"github.com/corvinusz/echo-xorm/ctx"
)

// Handler is a container for handlers and app data
type Handler struct {
	C *ctx.Context
}

// Result defines http response on GET /version
type Result struct {
	Result     string `json:"result"`
	Version    string `json:"version"`
	ServerTime int64  `json:"server_time"`
}

// GetVersion is a GET /version handler
func (h *Handler) GetVersion(c echo.Context) error {
	vr := Result{
		Result:     "OK",
		Version:    h.C.Config.Version,
		ServerTime: time.Now().UTC().Unix(),
	}
	return c.JSON(http.StatusOK, vr)
}
