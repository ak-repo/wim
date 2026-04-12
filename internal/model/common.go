package model

type PaginatedResult[T any] struct {
	Data  []T `json:"data"`
	Count int `json:"count"`
}

type BaseModel struct {
	ID        int     `db:"id" json:"id"`
	CreatedAt string  `db:"created_at" json:"createdAt"`
	UpdatedAt string  `db:"updated_at" json:"updatedAt"`
	DeletedAt *string `db:"deleted_at" json:"deletedAt,omitempty"`
}
