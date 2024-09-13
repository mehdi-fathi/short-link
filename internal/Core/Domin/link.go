package Domin

import "database/sql"

// Character is one character from the database.
type Link struct {
	ID        int64
	Link      string
	ShortKey  string
	Visit     int
	UpdatedAt sql.NullTime
	Status    string
	CreatedAt string
}

const (
	LINK_STATUS_APPROVE = "approve"
	Link_STATUS_REJECT  = "reject"
	LINK_STATUS_PENDING = "pending"
)
