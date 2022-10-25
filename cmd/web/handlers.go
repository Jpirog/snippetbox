package main

import (
	"errors"
	"fmt"

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

	data := app.newTemplateData(r) // diff than book
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	log.Print("Responding to a request at ", time.Now())
	log.Print(" -- Host", r.Host, r.URL.Path)
	log.Print(" -- RemoteAddr", r.RemoteAddr)
	log.Print(w.Header())
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetView request", r.URL.Path)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)

}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetCreate request", r.URL.Path)
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "Post")
		// w.WriteHeader(405)
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
