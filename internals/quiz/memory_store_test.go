package quiz_test

import (
	"github.com/mbsof31/go-quiz/internals/quiz"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQuiz(t *testing.T) {
	q := quiz.NewQuiz()
	assert.NotNil(t, q)
	assert.Equal(t, "Untitled quiz", q.Name)
	assert.Equal(t, "", q.Description)
	assert.NotNil(t, q.Questions)
	assert.NotNil(t, q.Meta)
}

func TestNewQuestion(t *testing.T) {
	q := quiz.NewQuestion()
	assert.NotNil(t, q)
	assert.Equal(t, "multi-choice", q.Type)
	assert.Equal(t, "Untitled question", q.Content)
	assert.NotNil(t, q.Choices)
	assert.NotNil(t, q.Meta)
}

func TestNewChoice(t *testing.T) {
	c := quiz.NewChoice()
	assert.NotNil(t, c)
	assert.Equal(t, "Choice 1", c.Content)
	assert.False(t, c.IsCorrect)
	assert.NotNil(t, c.Meta)
}

func TestMemoryStore(t *testing.T) {
	store := quiz.NewStore()

	// Test storing a quiz
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz 1"
	q1.Questions = []quiz.Question{*quiz.NewQuestion()} // Ensure it has at least one assignment
	err := store.Store(*q1)
	assert.NoError(t, err)

	// Test finding a quiz by ID
	q, err := store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "Quiz 1", q.Name)

	// Test listing all quizzes
	quizzes := store.ListAllQuizzes()
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

	// Test adding an assignment
	question := quiz.NewQuestion()
	question.Content = "Question 1"
	err = store.AddAssignment(1, *question)
	assert.NoError(t, err)
	q, err = store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, q.Questions, 2)

	// Test removing an assignment
	err = store.RemoveAssignment(1, 0)
	assert.NoError(t, err)
	q, err = store.FindQuizByID(1)
	assert.NoError(t, err)
	assert.Len(t, q.Questions, 1)

	// Test deleting a quiz
	err = store.Delete(1)
	assert.NoError(t, err)
	_, err = store.FindQuizByID(1)
	assert.Error(t, err)

	// Test search quizzes
	q2 := quiz.NewQuiz()
	q2.Name = "Quiz 2"
	q2.Questions = []quiz.Question{*quiz.NewQuestion()} // Ensure it has at least one assignment
	err = store.Store(*q2)
	assert.NoError(t, err)
	results, err := store.SearchQuiz("Quiz 2")
	assert.NoError(t, err)
	assert.Len(t, results, 1)
}

func TestMemoryStore_ImportExport(t *testing.T) {
	store := quiz.NewStore()

	// Create and store a quiz
	q1 := quiz.NewQuiz()
	q1.Name = "Quiz for Export"
	q1.Questions = []quiz.Question{*quiz.NewQuestion()} // Ensure it has at least one assignment
	err := store.Store(*q1)
	assert.NoError(t, err)

	// Export quizzes to a file
	filename := "quizzes.json"
	err = store.ExportQuizzes(filename)
	assert.NoError(t, err)

	// Ensure file is created
	_, err = os.Stat(filename)
	assert.NoError(t, err)

	// Create a new store and import quizzes
	newStore := quiz.NewStore()
	err = newStore.ImportQuizzes(filename)
	assert.NoError(t, err)

	// Verify imported quizzes
	importedQuizzes := newStore.ListAllQuizzes()
	assert.Len(t, importedQuizzes, 1)
	assert.Equal(t, "Quiz for Export", importedQuizzes[0].Name)

	// Cleanup
	err = os.Remove(filename)
	assert.NoError(t, err)
}

func TestMemoryStore_ValidateQuiz(t *testing.T) {
	store := quiz.NewStore()

	// Test empty name
	q1 := quiz.NewQuiz()
	q1.Name = ""
	err := store.ValidateQuiz(q1)
	assert.Error(t, err)
	assert.Equal(t, "quiz name cannot be empty", err.Error())

	// Test empty assignments
	q2 := quiz.NewQuiz()
	q2.Questions = []quiz.Question{}
	err = store.ValidateQuiz(q2)
	assert.Error(t, err)
	assert.Equal(t, "quiz must have at least one assignment", err.Error())

	// Test valid quiz
	q3 := quiz.NewQuiz()
	q3.Name = "Valid Quiz"
	q3.Questions = []quiz.Question{*quiz.NewQuestion()}
	err = store.ValidateQuiz(q3)
	assert.NoError(t, err)
}
