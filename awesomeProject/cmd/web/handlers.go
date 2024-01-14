package main

import (
	"errors"
	"fmt"
	"lina.net/aitunewstask/pkg/models"
	"log"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	newsList, err := app.news.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		NewsList: newsList,
	})
}

func (app *application) showNews(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.news.Get(id)
	if err != nil {
		log.Println(err) // or fmt.Println(err) for debugging
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		News: s,
	})
}

func (app *application) createNews(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling createNews request...")

	// Ensure that the request method is POST
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Extract form data
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	category := r.PostForm.Get("category")

	// Check if required fields are not empty
	if title == "" || content == "" || expires == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the category
	if category != "staff" && category != "students" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Call Insert with the additional category parameter
	id, err := app.news.Insert(title, content, expires, category)
	if err != nil {
		// Log the error and return a server error
		app.serverError(w, err)
		return
	}

	// Redirect to the news item page
	http.Redirect(w, r, fmt.Sprintf("/news?id=%d", id), http.StatusSeeOther)

	log.Println("createNews request handled.")
}

func (app *application) staffNews(w http.ResponseWriter, r *http.Request) {
	newsList, err := app.news.LatestByCategory("staff")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "for_staff.tmpl", &templateData{
		NewsList: newsList,
	})
}

func (app *application) studentsNews(w http.ResponseWriter, r *http.Request) {
	newsList, err := app.news.LatestByCategory("students")
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "for_students.tmpl", &templateData{
		NewsList: newsList,
	})
}
