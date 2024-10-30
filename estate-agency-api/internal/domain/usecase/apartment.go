package usecase

import (
	"context"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

type ApartmentService interface {
	GetAll(ctx context.Context, page int, pageSize int) (apartments []*entity.Apartment, err error)
	GetByID(ctx context.Context, id int) (apartment *entity.Apartment, err error)
	Create(ctx context.Context, apartment *entity.Apartment) (id int64, err error)
	Update(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error)
	Delete(ctx context.Context, id int) error
}

type apartmentUsecase struct {
	apartmentService ApartmentService
	realtorService   RealtorService
}

func NewApartmentUsecase(apartmentService ApartmentService, realtorService RealtorService) *apartmentUsecase {
	return &apartmentUsecase{apartmentService: apartmentService, realtorService: realtorService}
}

func (u *apartmentUsecase) GetAllApartment(ctx context.Context, page int, pageSize int) ([]*entity.Apartment, error) {
	return u.apartmentService.GetAll(ctx, page, pageSize)
}

func (u *apartmentUsecase) GetApartmentByID(ctx context.Context, id int) (apartment *entity.Apartment, realtor *entity.Realtor, err error) {
	apartment, err = u.apartmentService.GetByID(ctx, id)
	if err != nil {
		return
	}
	realtor, err = u.realtorService.GetByID(ctx, apartment.IDRealtor)
	return apartment, realtor, err
}

func (u *apartmentUsecase) CreateApartment(ctx context.Context, apartment *entity.Apartment) (id int64, err error) {
	return u.apartmentService.Create(ctx, apartment)
}

func (u *apartmentUsecase) UpdateApartment(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error) {
	if apartment.IDRealtor != 0 {
		realtor, err := u.realtorService.GetByID(ctx, apartment.IDRealtor)
		if err != nil {
			return aff, err
		}
		_ = realtor
	}

	return u.apartmentService.Update(ctx, id, apartment)
}

func (u *apartmentUsecase) DeleteApartment(ctx context.Context, id int) error {
	return u.apartmentService.Delete(ctx, id)
}
