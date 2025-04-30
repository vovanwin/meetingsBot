// schema/chat.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Chat — Telegram-группа или супергруппа.
type Chat struct {
	ent.Schema
}

func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("telegram_id").
			Unique().
			Comment("ID группы в Telegram"),
		field.String("title").
			Optional().
			Comment("Название группы"),
	}
}

func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		// связи через Membership
		edge.To("memberships", Membership.Type),
		//edge.To("users", User.Type).
		//	Through("memberships", "user"),
		// сборы
		edge.To("gathers", Gather.Type),
	}
}
