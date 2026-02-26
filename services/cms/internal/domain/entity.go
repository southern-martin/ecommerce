package domain

import "time"

// Banner represents a promotional banner.
type Banner struct {
	ID             string
	Title          string
	ImageURL       string
	LinkURL        string
	Position       string
	SortOrder      int
	TargetAudience string
	StartsAt       time.Time
	EndsAt         *time.Time
	IsActive       bool
	CreatedAt      time.Time
}

// Page represents a CMS content page.
type Page struct {
	ID              string
	Title           string
	Slug            string
	ContentHTML     string
	MetaTitle       string
	MetaDescription string
	Status          PageStatus
	PublishedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// PageStatus represents the publication status of a page.
type PageStatus string

const (
	PageStatusDraft     PageStatus = "draft"
	PageStatusPublished PageStatus = "published"
	PageStatusScheduled PageStatus = "scheduled"
)

// ContentSchedule represents a scheduled content action.
type ContentSchedule struct {
	ID          string
	ContentType string
	ContentID   string
	Action      string
	ScheduledAt time.Time
	Executed    bool
	CreatedAt   time.Time
}
