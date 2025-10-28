package entity

type GitlabWebhookRequest struct {
	ProjectId        int                                  `json:"project_id"`
	ObjectAttributes GitlabWebhookRequestObjectAttributes `json:"object_attributes"`
	User             GitlabWebhookRequestUser             `json:"user"`
	Note             GitlabWebhookRequestNote             `json:"note"`
}

type GitlabWebhookRequestObjectAttributes struct {
	Action       string `json:"action"`
	Name         string `json:"name"`
	AwardedOnUrl string `json:"awarded_on_url"`
}

type GitlabWebhookRequestUser struct {
	Username string `json:"username"`
}

type GitlabWebhookRequestNote struct {
	Description string `json:"description"`
}
