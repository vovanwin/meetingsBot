package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User — участник Telegram-чата.
type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("telegram_id").Unique().Comment("ID пользователя в Telegram"),
		field.String("username").Optional(),
		field.String("first_name").Optional(),
		field.String("last_name").Optional(),
		field.Bool("is_owner").Default(false),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// связи через Membership
		edge.To("memberships", Membership.Type),
		//edge.To("chats", Chat.Type).
		//	Through("memberships", "chat"),
		// голоса
		edge.To("votes", Vote.Type),
		//edge.To("gathers", Gather.Type).
		//	Through("votes", "gather"),
	}
}
