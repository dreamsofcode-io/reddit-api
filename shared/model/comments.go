package model

import "time"

type Comment struct {
	ID        string    `json:"id"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"time"`
	Content   string    `json:"comment"`
	Points    int       `json:"points"`
	IsDeleted bool      `json:"isDeleted"`
	Children  []Comment `json:"children"`
}
