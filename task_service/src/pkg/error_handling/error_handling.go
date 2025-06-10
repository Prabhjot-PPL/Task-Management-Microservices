package errorhandling

import (
	"net/http"
	pkgresponse "task_service/src/pkg/response"
)

func HandlerError(w http.ResponseWriter, msg string, statusCode int, err error) {

	response := pkgresponse.StandardResponse{
		Status:  "FAILURE",
		Message: msg,
	}
	pkgresponse.WriteResponse(w, statusCode, response)

}
