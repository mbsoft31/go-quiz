package quiz

type Question struct {
	ID      uint     `gorm:"primaryKey"`
	QuizID  uint     `gorm:"index"` // Foreign key
	Type    string   `json:"type" form:"type"`
	Content string   `json:"content" form:"content"`
	Choices []Choice `gorm:"foreignKey:QuestionID" json:"choices,omitempty" form:"choices,omitempty"`
	Meta    JSONMap  `gorm:"type:json" json:"meta,omitempty" form:"meta,omitempty"`
}

func NewQuestion() *Question {
	return &Question{
		Type:    "multi-choice",
		Content: "Untitled question",
		Choices: []Choice{},
		Meta:    make(map[string]interface{}),
	}
}
