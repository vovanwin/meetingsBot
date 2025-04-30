// schema/vote.go
package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/pkg/errors"
)

// Vote — голос пользователя в сборе.
type Vote struct {
	ent.Schema
}

func (Vote) Fields() []ent.Field {
	return []ent.Field{
		field.Int("count").
			Comment("От +1 до +5 или 0 для отказа").
			Validate(func(n int) error {
				if n < 0 || n > 5 {
					return errors.Errorf("некорректное значение голоса: %d", n)
				}
				return nil
			}),
		field.Time("voted_at").
			Default(time.Now),
	}
}

func (Vote) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("votes").
			Unique(), // один голос на одного user для одного gather
		edge.From("gather", Gather.Type).
			Ref("votes").
			Unique(),
	}
}
