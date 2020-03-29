package main

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

var templates = template.Must(template.ParseFiles("./tmpl/edit.html", "./tmpl/view.html", "./tmpl/home.html"))

var validPath = regexp.MustCompile("^/(edit|save|view|home)/(([a-zA-Z0-9\\s]+)|([a-zA-Z0-9\\s]?))$")

//var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9\\s]+)|[a-zA-Z0-9]?)$")

type Page struct {
	Title string
	Body  []byte
	List  []string
}

func (p *Page) save() error {
	filename := "./data/" + p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "./data/" + title + ".txt"
	body, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Page{Title: title, Body: body, List: grabPages()}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, loadErr := loadPage(title)

	if loadErr != nil {
		fmt.Printf("Unable to load page=%s=", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, loadErr := loadPage(title)

	if loadErr != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	filename := "./tmpl/home.html"
	body, homeErr := ioutil.ReadFile(filename)

	if homeErr != nil {
		fmt.Println("Unable to load homepage (make sure there is a home.html in the folder you ran this from)")
	}

	p := &Page{Title: "Home", Body: body, List: grabPages()}

	renderTemplate(w, "home", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	saveErr := p.save()

	if saveErr != nil {
		fmt.Printf("Unable to save=%s=", p.Title)
		http.Error(w, saveErr.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	var rendErr error
	if tmpl == "home" {
		rendErr = templates.ExecuteTemplate(w, tmpl+".html", p.List)
	} else {
		rendErr = templates.ExecuteTemplate(w, tmpl+".html", p)
	}

	if rendErr != nil {
		fmt.Printf("Unable to render page=%s=", tmpl+".html")
		http.Error(w, rendErr.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)

	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}

	return m[2], nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			fmt.Println("makeHandler error")
			http.Redirect(w, r, "/home/", http.StatusFound)
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func grabPages() []string {
	files, readErr := ioutil.ReadDir("./data/")

	if readErr != nil {
		fmt.Printf("Cannot read `.` dir")
	}

	var list []string

	for _, f := range files {
		check, _ := regexp.MatchString(".txt$", f.Name())
		if check {
			fmt.Println(f.Name())
			list = append(list, f.Name()[0:len(f.Name())-4])
		}
	}

	return list
}

func main() {
	http.HandleFunc("/", makeHandler(homeHandler))
	http.HandleFunc("/home/", makeHandler(homeHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
