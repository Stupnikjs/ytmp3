package api

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
	if len(eqlsplit) < 1 {
		w.Write([]byte(`
		<p> error malformed </p>
		`))
	}

	filename, err := FFmpegWrap(eqlsplit[1])

	if err != nil {
		w.Write([]byte(`<div> Something wrong happened </div>`))
		return
	}
	tempdir := strings.Split(filename, "\\")[2]

	w.Write([]byte(
		fmt.Sprintf(`
		<p>Dowload mp3</p>
		<a href="/fileupload/%s"> 
			  %s 
		</a>
		`, tempdir, filename)))
	return

}

func (app *Application) DowloadSound(w http.ResponseWriter, r *http.Request) {
	tempDir := chi.URLParam(r, "tempdir") + "\\"
	path := filepath.Join(DownloadDir, tempDir)
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		w.Write([]byte("error opening tempdir"))
		return
	}
	for i, e := range dirEntry {
		if i > 0 {
			break
		}
		filep := filepath.Join(path, e.Name())
		file, err := os.Open(filep)
		if err != nil {
			w.Write([]byte("error opening file"))
			return
		}
		defer file.Close()
		bytes, err := io.ReadAll(file)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Disposition", "attachement")
		w.Header().Set("filename-parm", fmt.Sprintf("filename=%s", file.Name()))
		w.Write(bytes)
		// for some reason delete the temp created before
		err = os.RemoveAll(DownloadDir)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		return
	}

}
