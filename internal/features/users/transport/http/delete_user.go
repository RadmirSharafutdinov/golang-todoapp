package users_transport_http

import (
	"net/http"
	core_logger "study/internal/core/logger"
	core_http_response "study/internal/core/transport/http/response"
	core_http_utils "study/internal/core/transport/http/utils"
)

func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to get userID path value",
		)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to delete user",
		)
		return
	}
	responseHandler.NoContentResponse()

}
