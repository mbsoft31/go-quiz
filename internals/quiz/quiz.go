package quiz

type Quiz struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `json:"name" form:"name"`
	Description string     `json:"description,omitempty" form:"description,omitempty"`
	Questions   []Question `gorm:"foreignKey:QuizID;references:ID" json:"questions,omitempty" form:"questions,omitempty"`
	Meta        JSONMap    `gorm:"type:json" json:"meta,omitempty" form:"meta,omitempty"`
}

var store = NewStore()

func NewQuiz() *Quiz {
	return &Quiz{
		Name:        "Untitled quiz",
		Description: "",
		Questions:   []Question{},
		Meta:        make(map[string]interface{}),
	}
}
