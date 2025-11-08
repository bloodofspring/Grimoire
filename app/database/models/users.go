package models

type User struct {
	TgID int64 `pg:",pk"`
	FullName string
	Username string
	CreatedAt int64 `pg:",default:extract(epoch from now())"`
	UpdatedAt int64 `pg:",default:extract(epoch from now())"`

	WrittenTexts []*Text `pg:"rel:has-many,join_fk:user_id"`
	IgnoredTopics []*IgnoredTopic `pg:"rel:has-many,join_fk:user_id"`
}
