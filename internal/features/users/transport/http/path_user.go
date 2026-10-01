package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"
	"study/internal/core/domain"
	core_logger "study/internal/core/logger"
	core_http_types "study/internal/core/transport/http"
	core_http_request "study/internal/core/transport/http/request"
	core_http_response "study/internal/core/transport/http/response"
	core_http_utils "study/internal/core/transport/http/utils"
)

type PathUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PathUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`Full name` cant be Null")
		}
		fullNameLen := len([]rune(*r.FullName.ToDomain().Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("FullName must be between 3 and 100")
		}
		if r.PhoneNumber.Set {
			if r.PhoneNumber.Value != nil {
				phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
				if phoneNumberLen < 10 || phoneNumberLen > 15 {
					return fmt.Errorf("PhoneNumber must be between 10 and 15")
				}
			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must startswith `+`")
			}
		}

	}
	return nil
}

type PathUserPesponse UserDTOResponse

func (h *UsersHTTPHandler) PathUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to get path value",
		)
		return
	}

	var request PathUserRequest
	if err := core_http_request.DecodeAndValidareRequest(r, &request); err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to decode and validate http request",
		)
		return
	}
	userPath := userPatchFromRequest(request)
	userDomain, err := h.usersService.PathUser(ctx, userID, userPath)
	if err != nil {
		responseHandler.ErrorResponce(
			err,
			"failed to Path User",
		)
	}
	response := PathUserPesponse(UserdtoFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
	log.Debug(
		fmt.Sprintf(
			"PathUserRequest fields:\nFullName:`%v`\nPhoneNumber:`%v`",
			request.FullName,
			request.PhoneNumber,
		),
	)
	rw.WriteHeader(http.StatusOK)
}
func userPatchFromRequest(request PathUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
