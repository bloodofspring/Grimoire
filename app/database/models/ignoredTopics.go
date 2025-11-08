package models

type IgnoredTopic struct {
	ID        int64 `pg:",pk"`
	ChatID int64
	ThreadID int
	CreatedAt int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt int64 `pg:",default:extract(epoch from now())"`

	UserID int64
	User *User `pg:"rel:has-one,fk:user_id"`
}