package ticket

import (
	"errors"
	"fmt"
	_const "mini/app/tools/const"
	"strings"

	"github.com/gin-gonic/gin"
)

type Form struct {
	CustomerName  string   `json:"customer_name"`
	CustomerPhone string   `json:"customer_phone"`
	Address       string   `json:"address"`
	Issue         string   `json:"issue"`
	Priority      string   `json:"priority"`
	Type          string   `json:"type"`
	TechnicianIDs []string `json:"technician_ids"`
}

func (f *Form) Validate() error {
	if f.CustomerName == "" {
		return errors.New("customer name is required")
	}
	if f.Issue == "" {
		return errors.New("issue is required")
	}
	if f.Type == "" {
		return errors.New("type is required")
	}
	if f.Type != _const.TicketTypeInstallation && f.Type != _const.TicketTypeTroubleshoot {
		return errors.New("type must be installation or troubleshoot")
	}
	if f.Priority != "" &&
		f.Priority != _const.TicketPriorityUrgent &&
		f.Priority != _const.TicketPriorityHigh &&
		f.Priority != _const.TicketPriorityMedium &&
		f.Priority != _const.TicketPriorityLow {
		return errors.New("priority is not valid")
	}
	if f.Priority == "" {
		f.Priority = _const.TicketPriorityLow
	}
	return nil
}

type UpdateForm struct {
	Issue              string   `json:"issue"`
	Priority           string   `json:"priority"`
	Status             string   `json:"status"`
	TechnicianResponse string   `json:"technician_response"`
	CancelReason       string   `json:"cancel_reason"`
	TechnicianIDs      []string `json:"technician_ids"`
}

func (f *UpdateForm) Validate() error {
	if f.Issue == "" {
		return errors.New("issue is required")
	}
	if f.Status == "" {
		return errors.New("status is required")
	}
	validStatuses := []string{
		_const.TicketStatusNew,
		_const.TicketStatusProgress,
		_const.TicketStatusHold,
		_const.TicketStatusDone,
		_const.TicketStatusProblem,
		_const.TicketStatusCancel,
	}
	valid := false
	for _, s := range validStatuses {
		if f.Status == s {
			valid = true
			break
		}
	}
	if !valid {
		return errors.New("status is not valid")
	}
	if f.Status == _const.TicketStatusCancel && f.CancelReason == "" {
		return errors.New("cancel_reason is required when status is cancel")
	}
	return nil
}

// QueryParam dipakai untuk filter list ticket
type QueryParam struct {
	Search   string
	Status   string
	Type     string
	Priority string
	Sort     string
	Order    string
	Page     int
	Limit    int
}

func NewQueryParam(c *gin.Context) *QueryParam {
	q := &QueryParam{
		Search:   c.Query("search"),
		Status:   c.Query("status"),
		Type:     c.Query("type"),
		Priority: c.Query("priority"),
		Sort:     c.DefaultQuery("sort", _const.TicketSortByCreatedAt),
		Order:    c.DefaultQuery("order", _const.Descending),
		Page:     1,
		Limit:    10,
	}

	// Validasi sort & order
	if q.Sort != _const.TicketSortByCreatedAt {
		q.Sort = _const.TicketSortByCreatedAt
	}
	if q.Order != _const.Ascending && q.Order != _const.Descending {
		q.Order = _const.Descending
	}

	return q
}

type Response struct {
	ID                 string `json:"id"`
	Code               string `json:"code"`
	Type               string `json:"type"`
	Priority           string `json:"priority"`
	Status             string `json:"status"`
	Issue              string `json:"issue"`
	TechnicianResponse string `json:"technician_response"`
	CancelReason       string `json:"cancel_reason"`
	CustomerName       string `json:"customer_name"`
	CustomerPhone      string `json:"customer_phone"`
	Address            string `json:"address"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`

	Technicians []TechnicianInfo `json:"technicians"`
}

type TechnicianInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

// GenerateCode buat kode tiket otomatis contoh TKT-00001
func GenerateCode(seq int) string {
	return "TKT-" + strings.ToUpper(fmt.Sprintf("%05d", seq))
}
