package types

type RequestBase struct {
	Urls []string `json:"urls"`
}

type AnswerBase struct {
	Result string `json:"result"`
}
