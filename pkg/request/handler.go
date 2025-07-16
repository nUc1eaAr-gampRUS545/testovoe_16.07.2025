package request

import (
	"net/http"
	"testovoe_16.07.2025/pkg/response"
)

func HandlerBody[T any](w *http.ResponseWriter, req *http.Request) (*T, error) {
	body, err := Decode[T](req.Body)
	if err != nil {
		response.Json(*w, err.Error(), 402)
		return nil, err
	}
	err = isValid[T](body)
	if err != nil {
		response.Json(*w, err.Error(), 422)
		return nil, err
	}
	return &body, nil

}
