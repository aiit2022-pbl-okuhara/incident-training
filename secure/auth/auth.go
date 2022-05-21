package auth

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/safehtml/template"

	"github.com/aiit2022-pbl-okuhara/incident-training/secure/responses"
	"github.com/aiit2022-pbl-okuhara/incident-training/storage"
)

const sessionCookie = "SESSION"

var unauthMsg = template.MustParseAndExecuteToHTML(`Please login before visiting this page.`)

// Interceptor is an auth (access control) interceptor.
type Interceptor struct {
	DB *storage.DB
}

// Before runs before the request is passed to the handler.
func (ip Interceptor) Before(w safehttp.ResponseWriter, r *safehttp.IncomingRequest, cfg safehttp.InterceptorConfig) safehttp.Result {
	// Identify the user.
	user := ip.userFromCookie(r)
	if user != "" {
		putUser(r.Context(), user)
	}

	if _, ok := cfg.(Skip); ok {
		// If the config says we should not perform auth, let's stop executing here.
		return safehttp.NotWritten()
	}

	if user == "" {
		// We have to perform auth, and the user was not identified, bail out.
		return w.WriteError(responses.Error{
			StatusCode:   safehttp.StatusUnauthorized,
			ErrorMessage: unauthMsg,
		})
	}
	return safehttp.NotWritten()
}

// Commit runs after the handler committed to a response.
func (ip Interceptor) Commit(w safehttp.ResponseHeadersWriter, r *safehttp.IncomingRequest, resp safehttp.Response, cfg safehttp.InterceptorConfig) {
	user := User(r)

	switch ctxSessionAction(r.Context()) {
	case clearSess:
		ip.DB.DelSession(user)
		w.AddCookie(safehttp.NewCookie(sessionCookie, ""))
	case setSess:
		token := ip.DB.GetToken(user)
		w.AddCookie(safehttp.NewCookie(sessionCookie, token))
	default:
		// do nothing
	}
}

func (Interceptor) Match(cfg safehttp.InterceptorConfig) bool {
	_, ok := cfg.(Skip)
	return ok
}

// User retrieves the user.
func User(r *safehttp.IncomingRequest) string {
	return ctxUser(r.Context())
}

func (ip Interceptor) userFromCookie(r *safehttp.IncomingRequest) string {
	sess, err := r.Cookie(sessionCookie)
	if err != nil || sess.Value() == "" {
		return ""
	}
	user, ok := ip.DB.GetUser(sess.Value())
	if !ok {
		return ""
	}
	return user
}

// ClearSession clears the session.
func ClearSession(r *safehttp.IncomingRequest) {
	putSessionAction(r.Context(), clearSess)
}

// CreateSession creates a session.
func CreateSession(r *safehttp.IncomingRequest, user string) {
	putSessionAction(r.Context(), setSess)
	putUser(r.Context(), user)
}

// Skip allows to mark an endpoint to skip auth checks.
type Skip struct{}
