package entity

type Reaction struct {
	Type            string `json:"type"`
	AiComment       string `json:"aiComment"`
	CreateDate      string `json:"createdDate"`
	GitlabProjectId string `json:"gitlabProjectId"`
	ReactionUrl     string `json:"reactionUrl"`
}

const (
	PositiveReaction = "PositiveReaction"
	NegativeReaction = "NegativeReaction"
)
