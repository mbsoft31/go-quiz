package quiz

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type SQLiteStore struct {
	DB *gorm.DB
}

func NewSQLiteStore(dsn string) (*SQLiteStore, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	store := &SQLiteStore{DB: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) migrate() error {
	return s.DB.AutoMigrate(&Quiz{}, &Question{}, &Choice{})
}

func (s *SQLiteStore) ListAllQuizzes() ([]*Quiz, error) {
	var quizzes []*Quiz
	result := s.DB.Find(&quizzes)
	return quizzes, result.Error
}

func (s *SQLiteStore) ListQuizzes(page, pageSize int) ([]*Quiz, error) {
	if page < 1 || pageSize < 1 {
		return nil, fmt.Errorf("invalid page (%d) or pageSize (%d)", page, pageSize)
	}
	var quizzes []*Quiz
	offset := (page - 1) * pageSize
	result := s.DB.Limit(pageSize).Offset(offset).Find(&quizzes)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list quizzes: %w", result.Error)
	}
	return quizzes, nil
}

func (s *SQLiteStore) FindQuizByID(id uint) (*Quiz, error) {
	var quiz Quiz
	result := s.DB.Preload("Questions.Choices").First(&quiz, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &quiz, nil
}

func (s *SQLiteStore) SearchQuiz(query string) ([]*Quiz, error) {
	var quizzes []*Quiz
	result := s.DB.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+query+"%", "%"+query+"%").Find(&quizzes)
	if result.Error != nil {
		return nil, result.Error
	}
	return quizzes, nil
}

func (s *SQLiteStore) ValidateQuiz(quiz *Quiz) error {
	if quiz.Name == "" {
		return fmt.Errorf("quiz name cannot be empty")
	}
	if len(quiz.Questions) == 0 {
		return fmt.Errorf("quiz must have at least one assignment")
	}
	return nil
}

func (s *SQLiteStore) Store(quiz Quiz) error {
	if err := s.ValidateQuiz(&quiz); err != nil {
		return err
	}
	for _, question := range quiz.Questions {
		if err := s.ValidateQuestion(&question); err != nil {
			return err
		}
		for _, choice := range question.Choices {
			if err := s.ValidateChoice(&choice); err != nil {
				return err
			}
		}
	}
	return s.DB.Create(&quiz).Error
}

func (s *SQLiteStore) Update(id uint, quiz Quiz) error {
	if err := s.ValidateQuiz(&quiz); err != nil {
		return err
	}
	return s.DB.Model(&Quiz{}).Where("id = ?", id).Updates(quiz).Error
}

func (s *SQLiteStore) Delete(id uint) error {
	return s.DB.Delete(&Quiz{}, id).Error
}

func (s *SQLiteStore) ExportQuizzes(filename string) error {
	quizzes, err := s.ListAllQuizzes()
	if err != nil {
		return fmt.Errorf("failed to list quizzes: %w", err)
	}
	data, err := json.MarshalIndent(quizzes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal quizzes: %w", err)
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}
	return nil
}

func (s *SQLiteStore) ImportQuizzes(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	var quizzes []*Quiz
	if err := json.Unmarshal(data, &quizzes); err != nil {
		return fmt.Errorf("failed to unmarshal quizzes: %w", err)
	}
	for _, quiz := range quizzes {
		if err := s.Store(*quiz); err != nil {
			return fmt.Errorf("failed to store quiz: %w", err)
		}
	}
	return nil
}

func (s *SQLiteStore) ValidateQuestion(question *Question) error {
	if question.Content == "" {
		return fmt.Errorf("question content cannot be empty")
	}
	if len(question.Choices) == 0 {
		return fmt.Errorf("question must have at least one choice")
	}
	return nil
}

func (s *SQLiteStore) AddAssignment(quizID uint, assignment Question) error {
	if err := s.ValidateQuestion(&assignment); err != nil {
		return err
	}
	var quiz Quiz
	if err := s.DB.Preload("Questions").First(&quiz, quizID).Error; err != nil {
		return err
	}
	quiz.Questions = append(quiz.Questions, assignment)
	return s.DB.Save(&quiz).Error
}

func (s *SQLiteStore) RemoveAssignment(quizID uint, assignmentID uint) error {
	var assignment Question
	result := s.DB.Where("quiz_id = ? AND id = ?", quizID, assignmentID).First(&assignment)
	if result.Error != nil {
		return result.Error
	}
	return s.DB.Delete(&assignment).Error
}

func (s *SQLiteStore) ValidateChoice(choice *Choice) error {
	if choice.Content == "" {
		return fmt.Errorf("choice content cannot be empty")
	}
	return nil
}

func (s *SQLiteStore) AddChoice(questionID uint, choice Choice) error {
	if err := s.ValidateChoice(&choice); err != nil {
		return err
	}
	choice.QuestionID = questionID
	return s.DB.Create(&choice).Error
}

func (s *SQLiteStore) RemoveChoice(questionID uint, choiceID uint) error {
	var choice Choice
	result := s.DB.Where("question_id = ? AND id = ?", questionID, choiceID).First(&choice)
	if result.Error != nil {
		return result.Error
	}
	return s.DB.Delete(&choice).Error
}
