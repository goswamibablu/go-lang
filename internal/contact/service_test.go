package contact

import (
	"context"
	"database/sql"
	"errors"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

var errCRUD = errors.New("error crud")

func TestCreatecontactRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     CreatecontactRequest
		wantError bool
	}{
		{"success", CreatecontactRequest{Name: "test"}, false},
		{"required", CreatecontactRequest{Name: ""}, true},
		{"too long", CreatecontactRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func TestUpdatecontactRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		model     UpdatecontactRequest
		wantError bool
	}{
		{"success", UpdatecontactRequest{Name: "test"}, false},
		{"required", UpdatecontactRequest{Name: ""}, true},
		{"too long", UpdatecontactRequest{Name: "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.Validate()
			assert.Equal(t, tt.wantError, err != nil)
		})
	}
}

func Test_service_CRUD(t *testing.T) {
	logger, _ := log.NewForTest()
	s := NewService(&mockRepository{}, logger)

	ctx := context.Background()

	// initial count
	count, _ := s.Count(ctx)
	assert.Equal(t, 0, count)

	// successful creation
	contact, err := s.Create(ctx, CreatecontactRequest{Name: "test"})
	assert.Nil(t, err)
	assert.NotEmpty(t, contact.ID)
	id := contact.ID
	assert.Equal(t, "test", contact.Name)
	assert.NotEmpty(t, contact.CreatedAt)
	assert.NotEmpty(t, contact.UpdatedAt)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// validation error in creation
	_, err = s.Create(ctx, CreatecontactRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	// unexpected error in creation
	_, err = s.Create(ctx, CreatecontactRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)

	_, _ = s.Create(ctx, CreatecontactRequest{Name: "test2"})

	// update
	contact, err = s.Update(ctx, id, UpdatecontactRequest{Name: "test updated"})
	assert.Nil(t, err)
	assert.Equal(t, "test updated", contact.Name)
	_, err = s.Update(ctx, "none", UpdatecontactRequest{Name: "test updated"})
	assert.NotNil(t, err)

	// validation error in update
	_, err = s.Update(ctx, id, UpdatecontactRequest{Name: ""})
	assert.NotNil(t, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// unexpected error in update
	_, err = s.Update(ctx, id, UpdatecontactRequest{Name: "error"})
	assert.Equal(t, errCRUD, err)
	count, _ = s.Count(ctx)
	assert.Equal(t, 2, count)

	// get
	_, err = s.Get(ctx, "none")
	assert.NotNil(t, err)
	contact, err = s.Get(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, "test updated", contact.Name)
	assert.Equal(t, id, contact.ID)

	// query
	contacts, _ := s.Query(ctx, 0, 0)
	assert.Equal(t, 2, len(contacts))

	// delete
	_, err = s.Delete(ctx, "none")
	assert.NotNil(t, err)
	contact, err = s.Delete(ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, id, contact.ID)
	count, _ = s.Count(ctx)
	assert.Equal(t, 1, count)
}

type mockRepository struct {
	items []entity.contact
}

func (m mockRepository) Get(ctx context.Context, id string) (entity.contact, error) {
	for _, item := range m.items {
		if item.ID == id {
			return item, nil
		}
	}
	return entity.contact{}, sql.ErrNoRows
}

func (m mockRepository) Count(ctx context.Context) (int, error) {
	return len(m.items), nil
}

func (m mockRepository) Query(ctx context.Context, offset, limit int) ([]entity.contact, error) {
	return m.items, nil
}

func (m *mockRepository) Create(ctx context.Context, contact entity.contact) error {
	if contact.Name == "error" {
		return errCRUD
	}
	m.items = append(m.items, contact)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, contact entity.contact) error {
	if contact.Name == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.ID == contact.ID {
			m.items[i] = contact
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	for i, item := range m.items {
		if item.ID == id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}
