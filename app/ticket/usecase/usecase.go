package usecase

import (
	"context"
	"errors"
	"fmt"
	ticketDomain "mini/app/ticket"
	ticketRepo "mini/app/ticket/repository"
	"mini/database/ent"

	"github.com/google/uuid"
)

type TicketUsecase interface {
	Create(ctx context.Context, f *ticketDomain.Form) (*ticketDomain.Response, error)
	GetByID(ctx context.Context, id string) (*ticketDomain.Response, error)
	GetAll(ctx context.Context, q *ticketDomain.QueryParam) ([]*ticketDomain.Response, int, error)
	Update(ctx context.Context, id string, f *ticketDomain.UpdateForm) (*ticketDomain.Response, error)
	Delete(ctx context.Context, id string) error
}

type ticketUsecase struct {
	db   *ent.Client
	repo ticketRepo.TicketRepository
}

func NewTicketUsecase(db *ent.Client, repo ticketRepo.TicketRepository) TicketUsecase {
	return &ticketUsecase{db: db, repo: repo}
}

func (uc *ticketUsecase) Create(ctx context.Context, f *ticketDomain.Form) (*ticketDomain.Response, error) {
	//Hitung Jumlah ticket untuk generate kode
	total, err := uc.repo.CountAll(ctx, uc.db)
	if err != nil {
		return nil, fmt.Errorf("failed count ticket: %w", err)
	}
	code := fmt.Sprintf("TKT-%05d", total+1)

	t, err := uc.repo.Create(ctx, uc.db, f, code)

	if err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return toResponse(t), nil
}

func (uc *ticketUsecase) GetByID(ctx context.Context, id string) (*ticketDomain.Response, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid ticket id")
	}

	t, err := uc.repo.FindByID(ctx, uc.db, uid)
	if err != nil {
		return nil, fmt.Errorf("ticket not found: %w", err)
	}

	return toResponse(t), nil
}

func (uc *ticketUsecase) GetAll(ctx context.Context, q *ticketDomain.QueryParam) ([]*ticketDomain.Response, int, error) {
	tickets, total, err := uc.repo.FindAll(ctx, uc.db, q)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all tickets: %w", err)
	}

	var result []*ticketDomain.Response
	for _, t := range tickets {
		result = append(result, toResponse(t))
	}

	return result, total, nil
}

func (uc *ticketUsecase) Update(ctx context.Context, id string, f *ticketDomain.UpdateForm) (*ticketDomain.Response, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid ticket id")
	}

	t, err := uc.repo.Update(ctx, uc.db, uid, f)
	if err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return toResponse(t), nil
}

func (uc *ticketUsecase) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid ticket id")
	}

	_, err = uc.repo.FindByID(ctx, uc.db, uid)
	if err != nil {
		return errors.New("ticket not found")
	}

	return uc.repo.Delete(ctx, uc.db, uid)
}

// toResponse convert ent.Ticket ke domain Response
func toResponse(t *ent.Ticket) *ticketDomain.Response {
	resp := &ticketDomain.Response{
		ID:                 t.ID.String(),
		Code:               t.Code,
		Type:               string(t.TicketType),
		Priority:           string(t.Priority),
		Status:             string(t.Status),
		Issue:              t.Issue,
		TechnicianResponse: t.TechnicianResponse,
		CancelReason:       t.CancelReason,
		CustomerName:       t.CustomerName,
		CustomerPhone:      t.CustomerPhone,
		Address:            t.Address,
		CreatedAt:          t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          t.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Technicians:        []ticketDomain.TechnicianInfo{},
	}

	for _, tech := range t.Edges.Technicians {
		resp.Technicians = append(resp.Technicians, ticketDomain.TechnicianInfo{
			ID:       tech.ID.String(),
			Name:     tech.Name,
			Position: tech.Position,
		})
	}

	return resp
}
