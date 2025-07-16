package response

import (
	"encoding/json"
	"net/http"
)

func Json(res http.ResponseWriter, data any, status int) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(data)

}
