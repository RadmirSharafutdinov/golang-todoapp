package users_transport_http

import (
	"fmt"
	"net/http"
	core_logger "study/internal/core/logger"
	core_http_response "study/internal/core/transport/http/response"
	core_http_utils "study/internal/core/transport/http/utils"
)

type GetUsersResponse []UserDTOResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)
	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponce(err, "failed to get `limit/offset` query param")
		return
	}
	userDomain, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to get users",
		)
	}
	response := GetUsersResponse(UsersDTOFromDomains(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	limit, err := core_http_utils.GetINTQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` query param:%w", err)
	}
	offset, err := core_http_utils.GetINTQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` query param:%w", err)
	}
	return limit, offset, nil
}
