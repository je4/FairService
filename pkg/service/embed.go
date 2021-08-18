package service

import "embed"

//go:embed static/dspace/oai.xsl
//go:embed static/dspace/img/lyncode.png
//go:embed static/bootstrap/css/bootstrap.min.css
//go:embed static/bootstrap/css/bootstrap.min.css.map
//go:embed static/bootstrap/js/bootstrap.bundle.min.js
//go:embed static/bootstrap/js/bootstrap.bundle.min.js.map
//go:embed static/dspace/js/jquery.js
//go:embed static/dspace/js/bootstrap.min.js
//go:embed static/dspace/css/style.css
//go:embed static/dspace/css/bootstrap.min.css
//go:embed static/dspace/css/bootstrap-theme.min.css
var staticFS embed.FS

//go:embed template/oai.gohtml
//go:embed template/partition.gohtml
var templateFS embed.FS

var templateFiles = map[string]string{
	"partition": "template/partition.gohtml",
	"oai":       "template/oai.gohtml",
}
