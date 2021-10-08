package service

import "embed"

//go:embed static/dspace/oai.xsl
//go:embed static/dspace/img/lyncode.png
//go:embed static/bootstrap/css/bootstrap.min.css
//go:embed static/bootstrap/css/bootstrap.min.css.map
//go:embed static/bootstrap/css/datatables.min.css
//go:embed static/bootstrap/js/bootstrap.bundle.min.js
//go:embed static/bootstrap/js/bootstrap.bundle.min.js.map
//go:embed static/bootstrap/icons/box-arrow-up-right.svg
//go:embed static/bootstrap/icons/file-binary.svg
//go:embed static/dspace/js/jquery.js
//go:embed static/dspace/js/bootstrap.min.js
//go:embed static/bootstrap/js/datatables.min.js
//go:embed static/dspace/css/style.css
//go:embed static/dspace/css/bootstrap.min.css
//go:embed static/dspace/css/bootstrap-theme.min.css
//go:embed static/img/DOI_logo.svg
var staticFS embed.FS

//go:embed template/oai.gohtml
//go:embed template/partition.gohtml
//go:embed template/dataviewer.gohtml
//go:embed template/deleted.gohtml
var templateFS embed.FS

var templateFiles = map[string]string{
	"deleted":    "template/deleted.gohtml",
	"partition":  "template/partition.gohtml",
	"oai":        "template/oai.gohtml",
	"dataviewer": "template/dataviewer.gohtml",
}
