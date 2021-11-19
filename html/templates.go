package html

import (
	"html/template"
	"io/fs"
	"log"
	"os"
)

func templateDir() fs.FS {
	templates, err := fs.Sub(os.DirFS("."), "html/templates")
	if err != nil {
		log.Fatalln(err, "template dir error")
	}
	return templates
}

func TemplateHtml() *template.Template {
	return template.Must(template.ParseFS(templateDir(), "*.html"))
}
