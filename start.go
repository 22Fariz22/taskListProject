package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type TaskList struct {
	Name string
	Task []byte
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("/(edit|save|view)/([a-zA-Z0-9]+)$")

func (t *TaskList) save() error {
	filename := t.Name + ".txt"
	return os.WriteFile(filename, t.Task, 0600)
}
func getName(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Name")
	}
	return m[2], nil
}
func loadTaskList(name string) (*TaskList, error) {
	filename := name + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &TaskList{Name: name, Task: body}, nil
}
func viewHandler(w http.ResponseWriter, r *http.Request, name string) {
	//name, err := getName(w, r)
	//if err != nil {
	//	return
	//}
	p, err := loadTaskList(name)
	if err != nil {
		http.Redirect(w, r, "/edit/"+name, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}
func editHandler(w http.ResponseWriter, r *http.Request, name string) {
	//name, err := getName(w, r)
	//if err != nil {
	//	return
	//}
	p, err := loadTaskList(name)
	if err != nil {
		p = &TaskList{Name: name}
	}
	renderTemplate(w, "edit", p)
}
func saveHandler(w http.ResponseWriter, r *http.Request, name string) {
	//name, err := getName(w, r)
	//if err != nil {
	//	return
	//}
	body := r.FormValue("body")
	p := &TaskList{Name: name, Task: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+name, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *TaskList) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		m := validPath.FindStringSubmatch(request.URL.Path)
		if m == nil {
			http.NotFound(writer, request)
			return
		}
		fn(writer, request, m[2])
	}
}
func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
