package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"
	core_errors "study/internal/core/errors"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

type validateble interface {
	Validate() error
}

func DecodeAndValidareRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument)
	}
	v, ok := dest.(validateble)
	var (
		err error
	)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}
	if err!=nil{
		return fmt.Errorf(
			"request validation:%v:%w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}
	return nil

}
