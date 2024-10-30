package dto

import "gilab.com/estate-agency-api/internal/domain/entity"

type ApartmentView struct {
	entity.Apartment
	RealtorView
}
