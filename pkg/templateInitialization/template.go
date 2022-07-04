package templateInitialization

import (
	"html/template"
	"os"
)

var templateNames = [5]string{ //TODO cfg
	"welcome",
	"login",
	"layout",
	"list",
	"movie",
}

func LoadTemplates(folder string) (map[string]*template.Template, error) {
	var templates = make(map[string]*template.Template, len(templateNames))
	for _, name := range templateNames {
		t, err := template.ParseFiles(folder+string(os.PathSeparator)+"layout.html",
			folder+string(os.PathSeparator)+name+".html")
		if err == nil {
			templates[name] = t
		} else {
			return nil, err
		}
	}
	return templates, nil
}
