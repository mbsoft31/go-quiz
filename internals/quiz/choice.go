package quiz

type Choice struct {
	ID         uint    `gorm:"primaryKey"`
	QuestionID uint    `gorm:"index"` // Foreign key
	Content    string  `json:"content" form:"content"`
	IsCorrect  bool    `json:"is_correct" form:"is_correct"`
	Thumb      []byte  `json:"thumb,omitempty" form:"thumb,omitempty"`
	Meta       JSONMap `gorm:"type:json" json:"meta,omitempty" form:"meta,omitempty"`
}

func NewChoice() *Choice {
	return &Choice{
		Content:   "Choice 1",
		IsCorrect: false,
		Thumb:     nil,
		Meta:      make(map[string]interface{}),
	}
}
