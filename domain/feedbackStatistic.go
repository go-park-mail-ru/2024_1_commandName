package domain

type AllStatistic struct {
	AllQuestionStatistic []OneQuestionStat
}

type OneQuestionStat struct {
	Grades        []int  `json:"grades"`
	QuestionID    int    `json:"question_id"`
	NSP           int    `json:"nsp,omitempty"`
	CSAP          int    `json:"csap,omitempty"`
	QuestionTitle string `json:"title"`
}
