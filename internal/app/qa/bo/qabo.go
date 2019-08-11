package bo

type QuestionAnswer struct {
	Id           string   `json:"id"`
	Question     []string `json:"question"`
	Answer       []string `json:"answer"`
	Tag          []string `json:"tag"`
	QuestionTime uint64   `json:"questionTime"`
	AnswerTime   uint64   `json:"answerTime"`
}
