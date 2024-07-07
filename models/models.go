package models

import (
	"time"
)

//go:generate reform

// reform:News
type News struct {
	ID      int64     `reform:"id,pk"`
	Title   string    `reform:"title"`
	Content string    `reform:"content"`
	Created time.Time `reform:"created"`
	Updated time.Time `reform:"updated"`
}

type NewsCategories struct {
	NewsID     int64 `reform:"news_id"`
	CategoryID int64 `reform:"category_id"`
}
