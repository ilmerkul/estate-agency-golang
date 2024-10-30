package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

type ApartmentStorage interface {
	GetAll(page int, pageSize int) (apartments []*entity.Apartment, err error)
	GetByID(id int) (apartment *entity.Apartment, err error)
	Create(apartment *entity.Apartment) (id int64, err error)
	Update(apartment *entity.Apartment) (aff int64, err error)
	Delete(id int) error
}

type apartmentService struct {
	storage ApartmentStorage
}

func NewApartmentService(storage ApartmentStorage) *apartmentService {
	return &apartmentService{storage: storage}
}

func (s *apartmentService) GetAll(ctx context.Context, page int, pageSize int) ([]*entity.Apartment, error) {
	return s.storage.GetAll(page, pageSize)
}

func (s *apartmentService) GetByID(ctx context.Context, id int) (realtor *entity.Apartment, err error) {
	return s.storage.GetByID(id)
}

func (s *apartmentService) Create(ctx context.Context, apartment *entity.Apartment) (id int64, err error) {
	apartment.UpdateTime = time.Now().Format("02.01.2006 15:04:05")
	apartment.CreateTime = time.Now().Format("02.01.2006 15:04:05")

	return s.storage.Create(apartment)
}

func (s *apartmentService) Update(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error) {
	r, err := s.storage.GetByID(id)
	if err != nil {
		return aff, err
	}
	apartment.ID = id

	if len(apartment.Title) == 0 {
		apartment.Title = r.Title
	}
	if apartment.Price == 0 {
		apartment.Price = r.Price
	}
	if len(apartment.City) == 0 {
		apartment.City = r.City
	}
	if apartment.Rooms == 0 {
		apartment.Rooms = r.Rooms
	}
	if len(apartment.Address) == 0 {
		apartment.Address = r.Address
	}
	if apartment.Square == 0 {
		apartment.Square = r.Square
	}
	if apartment.IDRealtor == 0 {
		apartment.IDRealtor = r.IDRealtor
	}
	apartment.UpdateTime = time.Now().Format("02.01.2006 15:04:05")
	apartment.CreateTime = r.CreateTime

	return s.storage.Update(apartment)
}

func (s *apartmentService) Delete(ctx context.Context, id int) error {
	if err := s.storage.Delete(id); err != nil {
		return err
	}

	return os.Remove(fmt.Sprintf("./../../internal/images/apartment/%d.png", id))
}
