package secure

import (
	"github.com/google/go-safeweb/safehttp"
	"github.com/google/go-safeweb/safehttp/plugins/coop"
	"github.com/google/go-safeweb/safehttp/plugins/csp"
	"github.com/google/go-safeweb/safehttp/plugins/fetchmetadata"
	"github.com/google/go-safeweb/safehttp/plugins/hostcheck"
	"github.com/google/go-safeweb/safehttp/plugins/hsts"
	"github.com/google/go-safeweb/safehttp/plugins/staticheaders"
	"github.com/google/go-safeweb/safehttp/plugins/xsrf/xsrfhtml"

	"github.com/aiit2022-pbl-okuhara/incident-training/secure/auth"
	"github.com/aiit2022-pbl-okuhara/incident-training/storage"
)

// NewMuxConfig creates a safe ServeMuxConfig.
func NewMuxConfig(db *storage.DB, addr string) *safehttp.ServeMuxConfig {
	c := safehttp.NewServeMuxConfig(dispatcher{})
	c.Intercept(coop.Default(""))
	c.Intercept(csp.Default(""))
	c.Intercept(&fetchmetadata.Interceptor{})
	c.Intercept(hostcheck.New(addr))
	c.Intercept(hsts.Default())
	c.Intercept(staticheaders.Interceptor{})
	c.Intercept(&xsrfhtml.Interceptor{SecretAppKey: "secret-key-that-should-not-be-in-sources"})
	c.Intercept(auth.Interceptor{DB: db})
	return c
}
