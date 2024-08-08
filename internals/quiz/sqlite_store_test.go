package quiz_test

import (
	"testing"

	"github.com/mbsof31/go-quiz/internals/quiz"
	"github.com/stretchr/testify/assert"
)

func setupStore(t *testing.T) *quiz.SQLiteStore {
	store, err := quiz.NewSQLiteStore("test.db")
	assert.NoError(t, err)
	return store
}

func teardownStore(store *quiz.SQLiteStore) {
	store.DB.Exec("DROP TABLE quizzes; DROP TABLE questions; DROP TABLE choices;") // Clean up
}

func TestSQLiteStore_Quiz(t *testing.T) {
	store := setupStore(t)
	defer teardownStore(store)

	// Test storing a quiz
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz 1"
	q1.Meta = quiz.JSONMap{"key": "value"}
	q1.Questions = []quiz.Question{*quiz.NewQuestion()} // Ensure it has at least one assignment
	q1.Questions[0].Choices = []quiz.Choice{*quiz.NewChoice()}
	err := store.Store(*q1)
	assert.NoError(t, err)

	// Test finding a quiz by ID
	q, err := store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Quiz 1", q.Name)
	assert.Equal(t, quiz.JSONMap{"key": "value"}, q.Meta)

	// Test listing all quizzes
	quizzes, err := store.ListAllQuizzes()
	assert.NoError(t, err)
	assert.Len(t, quizzes, 1)

	// Test listing quizzes with pagination
	pagedQuizzes, err := store.ListQuizzes(1, 1)
	assert.NoError(t, err)
	assert.Len(t, pagedQuizzes, 1)

	// Test updating a quiz
	q1.Description = "Updated Description"
	err = store.Update(1, *q1)
	assert.NoError(t, err)
	q, err = store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Description", q.Description)

	// Test deleting a quiz
	err = store.Delete(1)
	assert.NoError(t, err)
	_, err = store.FindQuizByID(1)
	assert.Error(t, err)

	// Test search quizzes
	q2 := quiz.NewQuiz()
	q2.Name = "Quiz 2"
	q2.Meta = quiz.JSONMap{"key": "another value"}
	q2.Questions = []quiz.Question{*quiz.NewQuestion()} // Ensure it has at least one assignment
	q2.Questions[0].Choices = []quiz.Choice{*quiz.NewChoice()}
	err = store.Store(*q2)
	assert.NoError(t, err)
	results, err := store.SearchQuiz("Quiz 2")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestSQLiteStore_Question(t *testing.T) {
	store := setupStore(t)
	defer teardownStore(store)

	// Test adding an assignment
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz with Question"
	q1.Questions = []quiz.Question{*quiz.NewQuestion()}
	q1.Questions[0].Choices = []quiz.Choice{*quiz.NewChoice()}
	err := store.Store(*q1)
	assert.NoError(t, err)

	// Ensure there is only one question initially
	q, err := store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, q.Questions, 1)

	question := quiz.NewQuestion()
	question.Content = "Question 1"
	question.Choices = []quiz.Choice{*quiz.NewChoice()}
	err = store.AddAssignment(1, *question)
	assert.NoError(t, err)

	q, err = store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, q.Questions, 2)

	// Test removing an assignment
	err = store.RemoveAssignment(1, q.Questions[1].ID) // Remove the second question added
	assert.NoError(t, err)
	q, err = store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, q.Questions, 1)
}

func TestSQLiteStore_Choice(t *testing.T) {
	store := setupStore(t)
	defer teardownStore(store)

	// Test adding choices to a question
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz with Choices"
	q := quiz.NewQuestion()
	q.Content = "Question with Choices"
	q.Choices = []quiz.Choice{*quiz.NewChoice()}
	q1.Questions = []quiz.Question{*q}
	err := store.Store(*q1)
	assert.NoError(t, err)

	qz, err := store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, qz.Questions[0].Choices, 1)
	assert.Equal(t, "Choice 1", qz.Questions[0].Choices[0].Content)
}

/*func TestImportExportQuizzes(t *testing.T) {
	store := setupStore(t)
	defer teardownStore(store)

	// Create and store a quiz with questions and choices
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz 1"
	q1.Questions = []quiz.Question{
		{
			Content: "Question 1",
			Choices: []quiz.Choice{
				{Content: "Choice 1"},
				{Content: "Choice 2"},
			},
		},
	}
	err := store.Store(*q1)
	assert.NoError(t, err)

	// Export quizzes to a file
	err = store.ExportQuizzes("test_export.json")
	assert.NoError(t, err)

	// Import quizzes from the file
	err = store.ImportQuizzes("test_import.json")
	assert.NoError(t, err)

	// Verify the imported quizzes
	quizzes, err := store.ListAllQuizzes()
	assert.NoError(t, err)
	assert.Len(t, quizzes, 1)
	assert.Equal(t, "Quiz 1", quizzes[0].Name)
	assert.Len(t, quizzes[0].Questions, 1)
	assert.Len(t, quizzes[0].Questions[0].Choices, 2)
}
*/
