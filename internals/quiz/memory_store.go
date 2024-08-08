package quiz

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type MemoryStore struct {
	sync.RWMutex
	Quizzes    map[int64]*Quiz
	QuizLastId int64
}

func NewStore() *MemoryStore {
	return &MemoryStore{
		Quizzes:    make(map[int64]*Quiz),
		QuizLastId: 0,
	}
}

func (s *MemoryStore) ListAllQuizzes() []*Quiz {
	s.RLock()
	defer s.RUnlock()

	quizzes := make([]*Quiz, 0, len(s.Quizzes))
	for _, quiz := range s.Quizzes {
		quizzes = append(quizzes, quiz)
	}
	return quizzes
}

func (s *MemoryStore) ListQuizzes(page, pageSize int) ([]*Quiz, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page or pageSize")
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	quizzes := s.ListAllQuizzes()
	if start >= len(quizzes) {
		return nil, fmt.Errorf("no more quizzes")
	}

	if end > len(quizzes) {
		end = len(quizzes)
	}

	return quizzes[start:end], nil
}

func (s *MemoryStore) FindQuizByID(id int64) (*Quiz, error) {
	s.RLock()
	defer s.RUnlock()

	quiz, found := s.Quizzes[id]
	if !found {
		return nil, fmt.Errorf("cannot find the quiz with the id of: %v", id)
	}
	return quiz, nil
}

func (s *MemoryStore) SearchQuiz(query string) ([]*Quiz, error) {
	s.RLock()
	defer s.RUnlock()

	results := make([]*Quiz, 0)
	for _, quiz := range s.Quizzes {
		if strings.Contains(strings.ToLower(quiz.Name), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(quiz.Description), strings.ToLower(query)) {
			results = append(results, quiz)
		}
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no quizzes found matching query: %v", query)
	}
	return results, nil
}

func (s *MemoryStore) ValidateQuiz(quiz *Quiz) error {
	if quiz.Name == "" {
		return fmt.Errorf("quiz name cannot be empty")
	}
	if len(quiz.Questions) == 0 {
		return fmt.Errorf("quiz must have at least one assignment")
	}
	return nil
}

func (s *MemoryStore) Store(quiz Quiz) error {
	if err := s.ValidateQuiz(&quiz); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	s.QuizLastId++
	s.Quizzes[s.QuizLastId] = &quiz
	return nil
}

func (s *MemoryStore) Update(id int64, quiz Quiz) error {
	if err := s.ValidateQuiz(&quiz); err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	_, found := s.Quizzes[id]
	if !found {
		return fmt.Errorf("cannot find the quiz with the id of: %v", id)
	}
	s.Quizzes[id] = &quiz
	return nil
}

func (s *MemoryStore) Delete(id int64) error {
	_, found := s.Quizzes[id]
	if !found {
		return fmt.Errorf("cannot find the quiz with the id of: %v", id)
	}
	delete(s.Quizzes, id)
	return nil
}

func (s *MemoryStore) ExportQuizzes(filename string) error {
	data, err := json.MarshalIndent(s.Quizzes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *MemoryStore) ImportQuizzes(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s.Quizzes)
}

func (s *MemoryStore) AddAssignment(quizID int64, assignment Question) error {
	quiz, err := s.FindQuizByID(quizID)
	if err != nil {
		return err
	}
	quiz.Questions = append(quiz.Questions, assignment)
	return nil
}

func (s *MemoryStore) RemoveAssignment(quizID int64, assignmentIndex int) error {
	quiz, err := s.FindQuizByID(quizID)
	if err != nil {
		return err
	}
	if assignmentIndex < 0 || assignmentIndex >= len(quiz.Questions) {
		return fmt.Errorf("invalid assignment index")
	}
	quiz.Questions = append(quiz.Questions[:assignmentIndex], quiz.Questions[assignmentIndex+1:]...)
	return nil
}
