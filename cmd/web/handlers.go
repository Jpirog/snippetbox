package main

import (
	"errors"
	"fmt"
	"html/template"

	"log"
	"net/http"
	"strconv"
	"time"

	"snippetbox.johnpirog.com/internal/models"
)

// this is the base home handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}

	// Initialize a slice containing the paths to the two files. It's important
	// to note that the file containing our base template must be the *first*
	// file in the slice.
	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/home.tmpl",
	// }

	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file
	// paths as a variadic parameter?
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.errorLog.Print("1" + err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }

	// err = ts.ExecuteTemplate(w, "base", nil)
	// if err != nil {
	// 	app.errorLog.Print("2" + err.Error())
	// 	app.serverError(w, err)
	// 	return
	// }

	// w.Header().Set("foo", "bar")
	// // w.Header().Set("ngrok-skip-browser-warning", "true")
	// w.Header()["ngrok-skip-browser-warning"] = []string{"true"}
	// w.Header().Add("user-agent", "ngrok-go")
	// w.Header().Add("foo", "bar")

	// 	headers := http.Header{
	// 		"ngrok-skip-browser-warning": []string{"true"},
	// 		"Accept": []string{"text/plain", "text/html"},
	// }
	// r.Header["ngrok-skip-browser-warning"] = []string{"true"}

	log.Print("Responding to a request at ", time.Now())
	log.Print(" -- Host", r.Host, r.URL.Path)
	log.Print(" -- RemoteAddr", r.RemoteAddr)
	log.Print(w.Header())
	// w.Write([]byte("Hello from Snippetbox, current timestamp is " + time.Now().String()))
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetView request", r.URL.Path)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	// fmt.Fprintf(w, "Displaying a specific snippet (%d)...\n", id)
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// fmt.Fprintf(w, "%+v", snippet)
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{
		Snippet: snippet,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetCreate request", r.URL.Path)
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "Post")
		w.WriteHeader(405)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
