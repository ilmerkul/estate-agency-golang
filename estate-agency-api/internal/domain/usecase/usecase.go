package usecase

import (
	"context"

	"gilab.com/estate-agency-api/internal/entity"
)

type ApartmentService interface {
	GetAll(ctx context.Context, page int, pageSize int) (apartments []*entity.Apartment, err error)
	GetByID(ctx context.Context, id int) (apartment *entity.Apartment, err error)
	Create(ctx context.Context, apartment *entity.Apartment) (id int64, err error)
	Update(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error)
	Delete(ctx context.Context, id int) error
}

type RealtorService interface {
	GetAll(ctx context.Context, page int, pageSize int) (realtors []*entity.Realtor, err error)
	GetByID(ctx context.Context, id int) (realtor *entity.Realtor, err error)
	Create(ctx context.Context, realtor *entity.Realtor) (id int64, err error)
	Update(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error)
	Delete(ctx context.Context, id int) error
}

type usecase struct {
	apartmentService ApartmentService
	realtorService   RealtorService
}

func NewUsecase(apartmentService ApartmentService, realtorService RealtorService) *usecase {
	return &usecase{apartmentService: apartmentService, realtorService: realtorService}
}

func (u *usecase) GetAllRealtor(ctx context.Context, page int, pageSize int) ([]*entity.Realtor, error) {
	return u.realtorService.GetAll(ctx, page, pageSize)
}

func (u *usecase) GetRealtorByID(ctx context.Context, id int) (realtor *entity.Realtor, err error) {
	return u.realtorService.GetByID(ctx, id)
}

func (u *usecase) CreateRealtor(ctx context.Context, realtor *entity.Realtor) (id int64, err error) {
	return u.realtorService.Create(ctx, realtor)
}

func (u *usecase) UpdateRealtor(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error) {
	return u.realtorService.Update(ctx, id, realtor)
}

func (u *usecase) DeleteRealtor(ctx context.Context, id int) error {
	return u.realtorService.Delete(ctx, id)
}

func (u *usecase) GetAllApartment(ctx context.Context, page int, pageSize int) ([]*entity.Apartment, error) {
	return u.apartmentService.GetAll(ctx, page, pageSize)
}

func (u *usecase) GetApartmentByID(ctx context.Context, id int) (apartment *entity.Apartment, realtor *entity.Realtor, err error) {
	apartment, err = u.apartmentService.GetByID(ctx, id)
	if err != nil {
		return
	}
	realtor, err = u.realtorService.GetByID(ctx, apartment.IDRealtor)
	return apartment, realtor, err
}

func (u *usecase) CreateApartment(ctx context.Context, apartment *entity.Apartment) (id int64, err error) {
	return u.apartmentService.Create(ctx, apartment)
}

func (u *usecase) UpdateApartment(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error) {
	return u.apartmentService.Update(ctx, id, apartment)
}

func (u *usecase) DeleteApartment(ctx context.Context, id int) error {
	return u.apartmentService.Delete(ctx, id)
}
