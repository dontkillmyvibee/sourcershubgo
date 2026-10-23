package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	coreerrors "github.com/dontkillmyvibee/sourcershubgo/internal/core/errors"
	"github.com/dontkillmyvibee/sourcershubgo/internal/core/logger"
	"go.uber.org/zap"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	w   http.ResponseWriter
}

func NewHTTPResponseHandler(log *logger.Logger, w http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		w:   w,
	}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.w.WriteHeader(statusCode)
	if err := json.NewEncoder(h.w).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.w.WriteHeader(http.StatusOK)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, coreerrors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, coreerrors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JSONResponse(response, statusCode)
}
