// schema/gather.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Gather — активный сбор в чате.
type Gather struct {
	ent.Schema
}

func (Gather) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("active").
			Default(true),
		field.Time("created_at").
			Default(time.Now),
		field.Time("closed_at").
			Optional(),
	}
}

func (Gather) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("gathers").
			Unique(),
		edge.To("votes", Vote.Type),
	}
}
