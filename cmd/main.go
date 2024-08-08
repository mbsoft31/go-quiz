package main

import (
	"github.com/mbsof31/go-quiz/internals"
	"github.com/mbsof31/go-quiz/internals/quiz"
	home "github.com/mbsof31/go-quiz/views/home"
	quizzes "github.com/mbsof31/go-quiz/views/quizzes"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	store, err := quiz.NewSQLiteStore("database/quiz.db")
	if err != nil {
		log.Fatalf("Error creating db store: %s", err.Error())
	}

	seedDB(store.DB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(internals.StoreMiddleware(store)) // Use the middleware

	// Home route
	r.Handle("/", templ.Handler(home.Home()))

	// Quiz routes
	r.Route("/quizzes", RegisterQuizRoutes)

	// Serve static files from the "public" directory
	fileServer(r, "/public", http.Dir("./public"))

	log.Println("Starting server on :4000")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		return
	}
}

func RegisterQuizRoutes(r chi.Router) {
	r.Get("/", quizListHandler)
	r.Get("/{quizID}", quizDetailsHandler)
	r.Get("/new", quizCreateHandler)
	r.Get("/{quizID}/edit", quizEditHandler)
}

func quizListHandler(w http.ResponseWriter, r *http.Request) {
	ctx := internals.GetAppContext(r)
	store := ctx.Store

	all, _ := store.ListAllQuizzes()
	err := quizzes.QuizListPage(all).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func quizCreateHandler(w http.ResponseWriter, r *http.Request) {
	err := quizzes.QuizFormPage(quiz.Quiz{}).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func quizDetailsHandler(w http.ResponseWriter, r *http.Request) {
	quizID := chi.URLParam(r, "quizID")
	// Here you'd normally fetch the q by ID from the database
	// This is just a sample q for demonstration
	var store, err = quiz.NewSQLiteStore("database/quiz.db")
	if err != nil {
		log.Fatalf("Error creating db store: %s", err.Error())
	}
	ID, err := strconv.Atoi(quizID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	q, err := store.FindQuizByID(uint(ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	err = quizzes.QuizDetailsPage(q).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func quizEditHandler(w http.ResponseWriter, r *http.Request) {
	//quizID := chi.URLParam(r, "quizID")
	// Here you'd normally fetch the q by ID from the database
	// This is just a sample q for demonstration
	q := quiz.Quiz{
		ID:          1,
		Name:        "Sample Quiz",
		Description: "This is a detailed description of the sample",
		Questions: []quiz.Question{
			{Content: "What is the capital of France?", Choices: []quiz.Choice{
				{Content: "Paris"},
				{Content: "London"},
				{Content: "Berlin"},
				{Content: "Madrid"},
			}},
			{Content: "What is the capital of UK?", Choices: []quiz.Choice{
				{Content: "Paris"},
				{Content: "London"},
				{Content: "Berlin"},
				{Content: "Madrid"},
			}},
			// Add more sample questions here
		},
	}

	err := quizzes.QuizFormPage(q).Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// fileServer conveniently sets up a http.FileServer handler to serve static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"/*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}

func seedDB(db *gorm.DB) {
	var quizCount int64
	db.Model(&quiz.Quiz{}).Count(&quizCount)
	if quizCount > 0 {
		log.Println("Quizzes already exist. No need to seed.")
		return
	}

	qs := []quiz.Quiz{
		{Name: "Quiz 1", Description: "This is a sample quiz."},
		{Name: "Quiz 2", Description: "This is a sample quiz 2."},
		{Name: "Quiz 3", Description: "This is a sample quiz 3."},
		{Name: "Quiz 4", Description: "This is a sample quiz 4."},
	}

	for i := range qs {
		for j := 0; j < 10; j++ {
			questionType := "single-choice"
			if j%2 == 0 {
				questionType = "multi-choice"
			} else if j%3 == 0 {
				questionType = "liquert-scale"
			}
			question := quiz.Question{
				Type:    questionType,
				Content: "Question " + strconv.Itoa(j+1),
				Choices: []quiz.Choice{
					{Content: "Choice 1", IsCorrect: j%2 == 0},
					{Content: "Choice 2", IsCorrect: j%3 == 0},
					{Content: "Choice 3", IsCorrect: j%5 == 0},
				},
			}
			qs[i].Questions = append(qs[i].Questions, question)
		}
	}

	for _, q := range qs {
		db.Create(&q)
	}

	log.Println("Database seeded successfully.")
}
