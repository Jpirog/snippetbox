package main

import (
	"errors"
	"fmt"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"snippetbox.johnpirog.com/internal/models"
)

// this is the base home handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r) // diff than book
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	log.Print("Responding to a request at: ", time.Now())
	log.Print(" -- Host: ", r.Host, r.URL.Path)
	log.Print(" -- RemoteAddr: ", r.RemoteAddr)
	log.Print(w.Header())
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetView request", r.URL.Path)
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
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
	w.Write([]byte("Display the form for a new snippet..."))
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	log.Print("snippetCreate request", r.URL.Path)

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
