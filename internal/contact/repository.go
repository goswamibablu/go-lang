package contact

import (
	"context"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/dbcontext"
	"github.com/qiangxue/go-rest-api/pkg/log"
)

// Repository encapsulates the logic to access contacts from the data source.
type Repository interface {
	// Get returns the contact with the specified contact ID.
	Get(ctx context.Context, id string) (entity.contact, error)
	// Count returns the number of contacts.
	Count(ctx context.Context) (int, error)
	// Query returns the list of contacts with the given offset and limit.
	Query(ctx context.Context, offset, limit int) ([]entity.contact, error)
	// Create saves a new contact in the storage.
	Create(ctx context.Context, contact entity.contact) error
	// Update updates the contact with given ID in the storage.
	Update(ctx context.Context, contact entity.contact) error
	// Delete removes the contact with given ID from the storage.
	Delete(ctx context.Context, id string) error
}

// repository persists contacts in database
type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

// NewRepository creates a new contact repository
func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

// Get reads the contact with the specified ID from the database.
func (r repository) Get(ctx context.Context, id string) (entity.contact, error) {
	var contact entity.contact
	err := r.db.With(ctx).Select().Model(id, &contact)
	return contact, err
}

// Create saves a new contact record in the database.
// It returns the ID of the newly inserted contact record.
func (r repository) Create(ctx context.Context, contact entity.contact) error {
	return r.db.With(ctx).Model(&contact).Insert()
}

// Update saves the changes to an contact in the database.
func (r repository) Update(ctx context.Context, contact entity.contact) error {
	return r.db.With(ctx).Model(&contact).Update()
}

// Delete deletes an contact with the specified ID from the database.
func (r repository) Delete(ctx context.Context, id string) error {
	contact, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.db.With(ctx).Model(&contact).Delete()
}

// Count returns the number of the contact records in the database.
func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("contact").Row(&count)
	return count, err
}

// Query retrieves the contact records with the specified offset and limit from the database.
func (r repository) Query(ctx context.Context, offset, limit int) ([]entity.contact, error) {
	var contacts []entity.contact
	err := r.db.With(ctx).
		Select().
		OrderBy("id").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&contacts)
	return contacts, err
}
