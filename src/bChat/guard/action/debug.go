package action

import (
	"net/http"
)

// Debug just for test
func Debug(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Hello Code"))
}
