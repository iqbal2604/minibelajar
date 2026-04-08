package schema

import (
	_const "mini/app/tools/const"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Ticket holds the schema definition for the Ticket entity.
type Ticket struct {
	ent.Schema
}

// Fields of the Ticket.
func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable().Unique(),
		field.String("code"),
		field.Enum("ticket_type").Values(_const.TicketTypeInstallation, _const.TicketTypeTroubleshoot),
		field.Enum("priority").Values(
			_const.TicketPriorityUrgent,
			_const.TicketPriorityHigh,
			_const.TicketPriorityMedium,
			_const.TicketPriorityLow,
		).Default(_const.TicketPriorityLow),
		field.Enum("status").Values(
			_const.TicketStatusNew,
			_const.TicketStatusProgress,
			_const.TicketStatusHold,
			_const.TicketStatusDone,
			_const.TicketStatusProblem,
			_const.TicketStatusCancel,
		).Default(_const.TicketStatusNew),
		field.Text("issue"),
		field.Text("technician_response").Default(""),
		field.Text("cancel_reason").Default(""),
		field.String("customer_name"),
		field.String("customer_phone"),
		field.String("address").Optional(),
	}
}

func (Ticket) Mixin() []ent.Mixin {
	return []ent.Mixin{TimeMixin{}, DeleteMixin{}}
}

func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("technicians", Technician.Type),
	}
}
