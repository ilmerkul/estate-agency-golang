package adapterRedis

import (
	"context"
	"strconv"
	"time"

	"gilab.com/estate-agency-api/internal/entity"
	"github.com/go-redis/cache/v9"
)

type apartmentAdapter struct {
	cacheClient *cache.Cache
}

func NewApartmentAdapter(cacheClient *cache.Cache) *apartmentAdapter {
	return &apartmentAdapter{cacheClient: cacheClient}
}

func (a *apartmentAdapter) SetApartment(apartment *entity.Apartment, id int) error {
	err := a.cacheClient.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   strconv.Itoa(id),
		Value: apartment,
		TTL:   time.Hour,
	})

	return err
}

func (a *apartmentAdapter) GetApartment(id int) (*entity.Apartment, error) {
	var apartment entity.Apartment

	err := a.cacheClient.Get(context.Background(), strconv.Itoa(id), &apartment)

	return &apartment, err
}
