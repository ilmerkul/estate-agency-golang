package httpModel

import "gilab.com/estate-agency-api/internal/entity"

type ApartmentView struct {
	entity.Apartment
	RealtorView
}
