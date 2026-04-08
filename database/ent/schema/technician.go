package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Technician struct{ ent.Schema }

func (Technician) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("name"),
		field.String("phone").Optional(),
		field.String("position").Optional(),
	}
}

func (Technician) Mixin() []ent.Mixin {
	return []ent.Mixin{TimeMixin{}, DeleteMixin{}}
}

func (Technician) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("tickets", Ticket.Type).Ref("technicians"),
	}
}
