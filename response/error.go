package response

import (
	"fmt"
	"net/http"
)

func Error(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(fmt.Sprintf(`{
		"message": "%s"
	}`, msg)))
}
