package mysql

import (
	"context"
	"database/sql"
	"time"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

type apartmentStorage struct {
	db *sql.DB
}

func NewApartmentStorage(db *sql.DB) *apartmentStorage {
	return &apartmentStorage{db: db}
}

func (as *apartmentStorage) GetAll(page int, pageSize int) (apartments []*entity.Apartment, err error) {
	context, close := context.WithTimeout(context.Background(), 3*time.Second)
	defer close()

	if err := as.db.PingContext(context); err != nil {
		return nil, err
	}

	rows, err := as.db.QueryContext(context, "SELECT * FROM apartments LIMIT ?,?", page*pageSize, pageSize)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		apartment := &entity.Apartment{}
		err = rows.Scan(&apartment.ID, &apartment.Title, &apartment.Price, &apartment.City, &apartment.Rooms, &apartment.Address, &apartment.Square, &apartment.IDRealtor, &apartment.UpdateTime, &apartment.CreateTime)
		if err != nil {
			return nil, err
		}
		apartments = append(apartments, apartment)
	}

	rows.Close()

	return apartments, nil
}

func (as *apartmentStorage) GetByID(id int) (apartment *entity.Apartment, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := as.db.PingContext(context); err != nil {
		return apartment, err
	}

	apartment = &entity.Apartment{}
	err = as.db.QueryRowContext(context, "SELECT * FROM apartments WHERE id=?", id).Scan(&apartment.ID, &apartment.Title, &apartment.Price, &apartment.City, &apartment.Rooms, &apartment.Address, &apartment.Square, &apartment.IDRealtor, &apartment.UpdateTime, &apartment.CreateTime)
	if err != nil {
		return apartment, err
	}

	return apartment, nil
}

func (as *apartmentStorage) Create(apartment *entity.Apartment) (id int64, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := as.db.PingContext(context); err != nil {
		return id, err
	}

	row, err := as.db.ExecContext(context, "INSERT INTO apartments (title, price, city, rooms, address, square, id_realtor, update_time, create_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", apartment.Title, apartment.Price, apartment.City, apartment.Rooms, apartment.Address, apartment.Square, apartment.IDRealtor, apartment.UpdateTime, apartment.CreateTime)
	if err != nil {
		return id, err
	}

	id, err = row.LastInsertId()

	return id, err
}

func (as *apartmentStorage) Update(apartment *entity.Apartment) (aff int64, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err = as.db.PingContext(context); err != nil {
		return aff, err
	}

	result, err := as.db.ExecContext(context, "UPDATE apartments SET title=?, price=?, city=?, rooms=?, address=?, square=?, id_realtor=?, update_time=?, create_time=? WHERE id=?", apartment.Title, apartment.Price, apartment.City, apartment.Rooms, apartment.Address, apartment.Square, apartment.IDRealtor, apartment.UpdateTime, apartment.CreateTime, apartment.ID)
	if err != nil {
		return aff, err
	}

	aff, err = result.RowsAffected()

	return aff, err
}

func (as *apartmentStorage) Delete(id int) error {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := as.db.PingContext(context); err != nil {
		return err
	}

	_, err := as.db.ExecContext(context, "DELETE FROM apartments WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
