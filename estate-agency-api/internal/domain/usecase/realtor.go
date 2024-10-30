package usecase

import (
	"context"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

type RealtorService interface {
	GetAll(ctx context.Context, page int, pageSize int) (realtors []*entity.Realtor, err error)
	GetByID(ctx context.Context, id int) (realtor *entity.Realtor, err error)
	Create(ctx context.Context, realtor *entity.Realtor) (id int64, err error)
	Update(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error)
	Delete(ctx context.Context, id int) error
}

type realtorUsecase struct {
	realtorService RealtorService
	/*apartmentService ApartmentService*/
}

func NewRealtorUsecase(realtorService RealtorService /*, apartmentService ApartmentService*/) *realtorUsecase {
	return &realtorUsecase{realtorService: realtorService /*, apartmentService: apartmentService*/}
}

func (u *realtorUsecase) GetAllRealtor(ctx context.Context, page int, pageSize int) ([]*entity.Realtor, error) {
	return u.realtorService.GetAll(ctx, page, pageSize)
}

func (u *realtorUsecase) GetRealtorByID(ctx context.Context, id int) (realtor *entity.Realtor, err error) {
	return u.realtorService.GetByID(ctx, id)
}

func (u *realtorUsecase) CreateRealtor(ctx context.Context, realtor *entity.Realtor) (id int64, err error) {
	return u.realtorService.Create(ctx, realtor)
}

func (u *realtorUsecase) UpdateRealtor(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error) {
	return u.realtorService.Update(ctx, id, realtor)
}

func (u *realtorUsecase) DeleteRealtor(ctx context.Context, id int) error {
	return u.realtorService.Delete(ctx, id)
}
