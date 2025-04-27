package models

import "time"

type Poll struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Options   []string  `json:"options"`
	Tags      []string  `json:"tags"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PollRequest struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
	Tags    []string `json:"tags"`
}

// PollResults holds the results of a poll
type PollResults struct {
	PollID int            `json:"pollId"`
	Votes  []OptionResult `json:"votes"`
}

// OptionResult holds the details of a poll option and its vote count
type OptionResult struct {
	Option string `json:"option"`
	Count  int    `json:"count"`
}
