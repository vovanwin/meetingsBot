// schema/membership.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Membership — связь «пользователь–чат» с ролью.
type Membership struct {
	ent.Schema
}

func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_admin").
			Default(false).
			Comment("Флаг администратора чата"),
	}
}

func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("memberships").
			Unique(),
		edge.From("chat", Chat.Type).
			Ref("memberships").
			Unique(),
	}
}
