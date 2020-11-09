package rest

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/go-chi/render"
	log "github.com/go-pkgz/lgr"
	R "github.com/go-pkgz/rest"
)

// ErrCode is used for client mapping and translation
type ErrCode int

// All error codes
const (
	ErrInternal   ErrCode = 0 // any internal error
	ErrDecode     ErrCode = 1 // failed to unmarshal incoming request
	ErrBadRequest ErrCode = 2 // request contains incorrect data or doesn't contain data
)

// SendErrorJSON makes {error: blah, details: blah} json body and responds with error code
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, errMsg string, errCode ErrCode) {
	render.Status(r, httpStatusCode)

	if err != nil {
		log.Printf("[WARN] %s", errDetailsMsg(r, httpStatusCode, err, errMsg))
		render.JSON(w, r, R.JSON{"error": err.Error(), "details": errMsg, "code": errCode})
		return
	}

	render.JSON(w, r, R.JSON{"error": nil, "details": errMsg, "code": errCode})
}

func errDetailsMsg(r *http.Request, code int, err error, msg string) string {
	q := r.URL.String()
	if qun, e := url.QueryUnescape(q); e == nil {
		q = qun
	}

	srcFileInfo := ""
	if pc, file, line, ok := runtime.Caller(2); ok {
		fnameElems := strings.Split(file, "/")
		funcNameElems := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		srcFileInfo = fmt.Sprintf(" [caused by %s:%d %s]", strings.Join(fnameElems[len(fnameElems)-3:], "/"),
			line, funcNameElems[len(funcNameElems)-1])
	}

	remoteIP := r.RemoteAddr
	if pos := strings.Index(remoteIP, ":"); pos >= 0 {
		remoteIP = remoteIP[:pos]
	}
	if err == nil {
		err = errors.New("no error")
	}
	return fmt.Sprintf("%s - %v - %d - %s - %s%s", msg, err, code, remoteIP, q, srcFileInfo)
}
