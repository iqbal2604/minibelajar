package repository

import (
	"context"
	"fmt"
	"mini/database/ent"
	"mini/database/ent/ticket"

	ticketDomain "mini/app/ticket"

	"github.com/google/uuid"
)

type TicketRepository interface {
	Create(ctx context.Context, client *ent.Client, f *ticketDomain.Form, code string) (*ent.Ticket, error)
	Update(ctx context.Context, client *ent.Client, id uuid.UUID, f *ticketDomain.UpdateForm) (*ent.Ticket, error)
	Delete(ctx context.Context, client *ent.Client, id uuid.UUID) error
	FindByID(ctx context.Context, client *ent.Client, id uuid.UUID) (*ent.Ticket, error)
	FindAll(ctx context.Context, client *ent.Client, q *ticketDomain.QueryParam) ([]*ent.Ticket, int, error)
}

type ticketRepository struct{}

func NewTicketRepository() TicketRepository {
	return &ticketRepository{}
}

func (r *ticketRepository) Create(ctx context.Context, client *ent.Client, f *ticketDomain.Form, code string) (*ent.Ticket, error) {
	q := client.Ticket.Create().
		SetCode(code).
		SetTicketType(ticket.TicketType(f.Type)).
		SetStatus(ticket.StatusNew).
		SetIssue(f.Issue).
		SetCustomerName(f.CustomerName).
		SetCustomerPhone(f.CustomerPhone)

	if f.Address != "" {
		q.SetAddress(f.Address)
	}

	return q.Save(ctx)
}

func (r *ticketRepository) FindByID(ctx context.Context, client *ent.Client, id uuid.UUID) (*ent.Ticket, error) {
	return client.Ticket.Query().
		Where(ticket.ID(id)).
		WithTechnicians().
		Only(ctx)
}

func (r *ticketRepository) FindAll(ctx context.Context, client *ent.Client, q *ticketDomain.QueryParam) ([]*ent.Ticket, int, error) {
	query := client.Ticket.Query().WithTechnicians()

	if q.Search != "" {
		query = query.Where(
			ticket.Or(
				ticket.CustomerNameContainsFold(q.Search),
				ticket.CodeContainsFold(q.Search),
			),
		)
	}
	if q.Status != "" {
		query = query.Where(ticket.StatusEQ(ticket.Status(q.Status)))
	}
	if q.Type != "" {
		query = query.Where(ticket.TicketTypeEQ(ticket.TicketType(q.Type)))
	}
	if q.Priority != "" {
		query = query.Where(ticket.PriorityEQ(ticket.Priority(q.Priority)))
	}

	total, err := query.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	offset := (q.Page - 1) * q.Limit
	tickets, err := query.
		Order(ent.Desc(ticket.FieldCreatedAt)).
		Limit(q.Limit).
		Offset(offset).
		All(ctx)

	return tickets, total, err

}

func (r *ticketRepository) Update(ctx context.Context, client *ent.Client, id uuid.UUID, f *ticketDomain.UpdateForm) (*ent.Ticket, error) {
	u := client.Ticket.UpdateOneID(id).
		SetIssue(f.Issue).
		SetStatus(ticket.Status(f.Status)).
		SetTechnicianResponse(f.TechnicianResponse).
		SetCancelReason(f.CancelReason)

	if f.Priority != "" {
		u.SetPriority(ticket.Priority(f.Priority))
	}
	return u.Save(ctx)
}

func (r *ticketRepository) Delete(ctx context.Context, client *ent.Client, id uuid.UUID) error {
	return client.Ticket.DeleteOneID(id).Exec(ctx)
}

func (r *ticketRepository) CountAll(ctx context.Context, client ent.Client) (int, error) {
	return client.Ticket.Query().Count(ctx)
}

// buildCode buat kode tiket berdasarkan total existing ticket
func buildCode(seq int) string {
	return fmt.Sprintf("TK-%05d", seq+1)
}
