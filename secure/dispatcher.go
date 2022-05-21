package secure

import (
	"net/http"

	"github.com/aiit2022-pbl-okuhara/incident-training/secure/responses"
	"github.com/aiit2022-pbl-okuhara/incident-training/secure/templates"
	"github.com/google/go-safeweb/safehttp"
)

type dispatcher struct {
	safehttp.DefaultDispatcher
}

func (d dispatcher) Error(rw http.ResponseWriter, resp safehttp.ErrorResponse) error {
	if ce, ok := resp.(responses.Error); ok {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(int(ce.Code()))
		return templates.All.ExecuteTemplate(rw, "error.tpl.html", ce.Message)
	}
	return d.DefaultDispatcher.Error(rw, resp)
}
