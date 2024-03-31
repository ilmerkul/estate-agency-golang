package handler

import (
	"context"

	"gilab.com/estate-agency-api/internal/entity"
)

type Usecase interface {
	GetAllApartment(ctx context.Context, page int, pageSize int) (apartments []*entity.Apartment, err error)
	GetApartmentByID(ctx context.Context, id int) (apartment *entity.Apartment, realtor *entity.Realtor, err error)
	CreateApartment(ctx context.Context, apartment *entity.Apartment) (id int64, err error)
	UpdateApartment(ctx context.Context, id int, apartment *entity.Apartment) (aff int64, err error)
	DeleteApartment(ctx context.Context, id int) error

	GetAllRealtor(ctx context.Context, page int, pageSize int) (realtors []*entity.Realtor, err error)
	GetRealtorByID(ctx context.Context, id int) (realtor *entity.Realtor, err error)
	CreateRealtor(ctx context.Context, realtor *entity.Realtor) (id int64, err error)
	UpdateRealtor(ctx context.Context, id int, realtor *entity.Realtor) (aff int64, err error)
	DeleteRealtor(ctx context.Context, id int) error
}
