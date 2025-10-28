package entity

type Prompt struct {
	Text            string `json:"text"`
	GitlabProjectId string `json:"gitlabProjectId"`
}
