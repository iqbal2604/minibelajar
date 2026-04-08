package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Blog struct{ ent.Schema }

func (Blog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("title"),
		field.Text("content"),
		field.Bool("is_published").Default(false),
		field.Time("publish_schedule").Optional().Nillable(),
	}
}

func (Blog) Mixin() []ent.Mixin {
	return []ent.Mixin{TimeMixin{}, DeleteMixin{}}
}
