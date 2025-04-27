package models

import "time"

type Poll struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Options   []string  `json:"options"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
}
