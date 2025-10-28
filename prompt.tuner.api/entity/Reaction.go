package entity

type Reaction struct {
	Type            string   `json:"type"`
	AiComments      []string `json:"aiComments"`
	CreateDate      string   `json:"createdDate"`
	GitlabProjectId string   `json:"gitlabProjectId"`
}

const (
	PositiveReaction = "PositiveReaction"
	NegativeReaction = "NegativeReaction"
)
