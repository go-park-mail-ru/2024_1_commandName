package domain

type AllStatistic struct {
	AllQuestionStatistic []OneQuestionStat
}

type OneQuestionStat struct {
	Grades        []int  `json:"grades"`
	QuestionID    int    `json:"question_id"`
	NSP           int    `json:"nsp,omitempty"`
	Type          string `json:"question_type"`
	CSAT          int    `json:"csap,omitempty"`
	QuestionTitle string `json:"title"`
}
