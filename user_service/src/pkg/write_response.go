package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user_service/src/pkg/logger"
)

type StandardResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

func WriteResponse(w http.ResponseWriter, statuscode int, resp StandardResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "ERROR : ", http.StatusInternalServerError)
	} else {
		fmt.Print("\n")
		logger.Log.Info(resp.Message)
	}
}
