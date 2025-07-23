package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) domain.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) GetByEmail(ctx context.Context, email string) (*domain.Customer, error) {
	var customer domain.Customer
	if err := r.db.WithContext(ctx).First(&customer, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}

func (r *customerRepository) GetOrCreate(ctx context.Context, customer *domain.Customer) (*domain.Customer, error) {
	// Try to find existing customer by email
	existing, err := r.GetByEmail(ctx, customer.Email)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		// Update existing customer with new information
		existing.CompanyName = customer.CompanyName
		existing.PersonName = customer.PersonName
		existing.IPAddress = customer.IPAddress

		if err := r.Update(ctx, existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	// Create new customer
	if err := r.Create(ctx, customer); err != nil {
		return nil, err
	}
	return customer, nil
}
