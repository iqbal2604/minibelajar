// database/ent/schema/notification.go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Notification struct{ ent.Schema }

func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("title"),
		field.String("description").Optional(),
		field.String("url").Optional(),
		field.Bool("is_published").Default(false),
		field.Time("publish_schedule"),
	}
}

func (Notification) Mixin() []ent.Mixin {
	return []ent.Mixin{TimeMixin{}, DeleteMixin{}}
}
