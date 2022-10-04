package contact

import (
   "context"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"time"
)

// Service encapsulates usecase logic for contacts.
type Service interface {
	Get(ctx context.Context, id string) (contact, error)
	Query(ctx context.Context, offset, limit int) ([]contact, error)
	Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreatecontactRequest) (contact, error)
	Update(ctx context.Context, id string, input UpdatecontactRequest) (contact, error)
	Delete(ctx context.Context, id string) (contact, error)
}

// contact represents the data about an contact.
type contact struct {
	entity.contact
}

// CreatecontactRequest represents an contact creation request.
type CreatecontactRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreatecontactRequest fields.
func (m CreatecontactRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdatecontactRequest represents an contact update request.
type UpdatecontactRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreatecontactRequest fields.
func (m UpdatecontactRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new contact service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the contact with the specified the contact ID.
func (s service) Get(ctx context.Context, id string) (contact, error) {
	contact, err := s.repo.Get(ctx, id)
	if err != nil {
		return contact{}, err
	}
	return contact{contact}, nil
}

// Create creates a new contact.
func (s service) Create(ctx context.Context, req CreatecontactRequest) (contact, error) {
	if err := req.Validate(); err != nil {
		return contact{}, err
	}
	id := entity.GenerateID()
	now := time.Now()
	err := s.repo.Create(ctx, entity.contact{
		ID:        id,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return contact{}, err
	}
	return s.Get(ctx, id)
}

// Update updates the contact with the specified ID.
func (s service) Update(ctx context.Context, id string, req UpdatecontactRequest) (contact, error) {
	if err := req.Validate(); err != nil {
		return contact{}, err
	}

	contact, err := s.Get(ctx, id)
	if err != nil {
		return contact, err
	}
	contact.Name = req.Name
	contact.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, contact.contact); err != nil {
		return contact, err
	}
	return contact, nil
}

// Delete deletes the contact with the specified ID.
func (s service) Delete(ctx context.Context, id string) (contact, error) {
	contact, err := s.Get(ctx, id)
	if err != nil {
		return contact{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return contact{}, err
	}
	return contact, nil
}

// Count returns the number of contacts.
func (s service) Count(ctx context.Context) (int, error) {
	return s.repo.Count(ctx)
}

// Query returns the contacts with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]contact, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []contact{}
	for _, item := range items {
		result = append(result, contact{item})
	}
	return result, nil
}
