package templates

import (
	"embed"

	"github.com/google/go-safeweb/safehttp/plugins/htmlinject"
	"github.com/google/safehtml/template"
)

// go:embed で埋め込める型は string / []byte / embed.FS の三種類
// embed.FS は階層型に読み込める

//go:embed *.tpl.html
var templatesFS embed.FS

var All *template.Template

func init() {
	var err error

	tplSrc := template.TrustedSourceFromConstant("*.tpl.html")
	All, err = htmlinject.LoadGlobEmbed(nil, htmlinject.LoadConfig{}, tplSrc, templatesFS)
	if err != nil {
		// TODO: 適切に error 処理を行う
		panic(err)
	}
}
