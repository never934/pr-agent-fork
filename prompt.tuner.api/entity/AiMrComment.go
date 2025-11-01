package entity

import "time"

type AiMrComment struct {
	CommentPoints   []string  `json:"commentPoints"`
	CreateDate      time.Time `json:"createdDate"`
	GitlabProjectId string    `json:"gitlabProjectId"`
	Url             string    `json:"url"`
	LikesCount      int       `json:"likesCount"`
	DislikesCount   int       `json:"dislikesCount"`
}
