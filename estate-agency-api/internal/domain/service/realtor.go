package service

import (
	"context"
	"fmt"
	"os"

	"gilab.com/estate-agency-api/internal/entity"
)

type RealtorStorage interface {
	GetAll(page int, pageSize int) (realtors []*entity.Realtor, err error)
	GetByID(id int) (realtor *entity.Realtor, err error)
	Create(realtor *entity.Realtor) (id int64, err error)
	Update(realtor *entity.Realtor) (aff int64, err error)
	Delete(id int) error
}

type realtorService struct {
	storage RealtorStorage
}

func NewRealtorService(storage RealtorStorage) *realtorService {
	return &realtorService{storage: storage}
}

func (s *realtorService) GetAll(ctx context.Context, page int, pageSize int) ([]*entity.Realtor, error) {
	return s.storage.GetAll(page, pageSize)
}

func (s *realtorService) GetByID(ctx context.Context, id int) (realtor *entity.Realtor, err error) {
	return s.storage.GetByID(id)
}

func (s *realtorService) Create(ctx context.Context, realtor *entity.Realtor) (id int64, err error) {
	return s.storage.Create(realtor)
}

func (s *realtorService) Update(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error) {
	r, err := s.storage.GetByID(id)
	if err != nil || len(r.FirstName) == 0 {
		return aff, err
	}
	realtor.ID = id

	if len(realtor.FirstName) == 0 {
		realtor.FirstName = r.FirstName
	}
	if len(realtor.LastName) == 0 {
		realtor.LastName = r.LastName
	}
	if len(realtor.Phone) == 0 {
		realtor.Phone = r.Phone
	}
	if len(realtor.Email) == 0 {
		realtor.Email = r.Email
	}
	if realtor.Rating == 0 {
		realtor.Rating = r.Rating
	}
	if realtor.Experience == 0 {
		realtor.Experience = r.Experience
	}

	return s.storage.Update(realtor)
}

func (s *realtorService) Delete(ctx context.Context, id int) error {
	if err := s.storage.Delete(id); err != nil {
		return err
	}

	return os.Remove(fmt.Sprintf("./../../internal/images/realtor/%d.png", id))
}
