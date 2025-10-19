package models

type Text struct {
	ID int64 `pg:",pk"`
	UserID int64
	Text string
	CreatedAt int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt int64 `pg:",default:extract(epoch from now())"`

	User *User `pg:"rel:has-one,fk:user_id"`
}