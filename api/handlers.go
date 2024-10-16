package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi/v5"
)

var pathToTemplates = "./static/templates/"

type TemplateData struct {
	Data map[string]any
}

func render(w http.ResponseWriter, r *http.Request, t string, td *TemplateData) error {
	_ = r.Method

	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, t), path.Join(pathToTemplates, "/base.layout.gohtml"))
	if err != nil {
		return err
	}
	err = parsedTemplate.Execute(w, td)
	if err != nil {
		return err
	}
	return nil

}

// template rendering

func (app *Application) RenderAccueil(w http.ResponseWriter, r *http.Request) {

	td := TemplateData{}
	td.Data = make(map[string]any)
	_ = render(w, r, "/main.gohtml", &td)
}
func (app *Application) PostVideoId(w http.ResponseWriter, r *http.Request) {

	reader, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()
	body := string(reader)

	eqlsplit := strings.Split(body, "=")

	filename := FFmpegWrap(eqlsplit[1])
	w.Write([]byte(
		fmt.Sprintf(`
		<a href="/download/%s"> 
			  %s  Download 
		</a>
		`, filename, filename)))
	return

}

func (app *Application) DowloadSound(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "name")
	file, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
	}

	bytes, err := io.ReadAll(file)

	defer os.Remove(filename)
	w.Header().Set("Content-Disposition", "attachement")
	w.Header().Set("filename-parm", fmt.Sprintf("filename=%s", filename))
	w.Write(bytes)

}
