package mysql

import (
	"context"
	"database/sql"
	"time"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

const (
	contextTimeGetAllApartment = 2
	contextTimeGetOneApartment = 1
	contextTimeCreateApartment = 1
	contextTimeUpdateApartment = 1
	contextTimeDeleteApartment = 1
)

type apartmentStorage struct {
	db      *sql.DB
	context context.Context
}

func NewApartmentStorage(db *sql.DB) *apartmentStorage {
	return &apartmentStorage{db: db, context: context.Background()}
}

func (as *apartmentStorage) GetAll(page int, pageSize int) (apartments []*entity.Apartment, err error) {

	q := `SELECT * FROM apartments LIMIT ?,?`

	context, close := context.WithTimeout(as.context, contextTimeGetAllApartment*time.Second)
	defer close()

	if err = as.db.PingContext(context); err != nil {
		return
	}

	stmt, err := as.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	rows, err := stmt.QueryContext(context, page*pageSize, pageSize)
	if err != nil {
		return
	}

	for rows.Next() {
		apartment := &entity.Apartment{}
		err = rows.Scan(&apartment.ID, &apartment.Title, &apartment.Price, &apartment.City, &apartment.Rooms, &apartment.Address, &apartment.Square, &apartment.IDRealtor, &apartment.UpdateTime, &apartment.CreateTime)
		if err != nil {
			return
		}
		apartments = append(apartments, apartment)
	}

	rows.Close()

	return apartments, nil
}

func (as *apartmentStorage) GetByID(id int) (apartment *entity.Apartment, err error) {

	q := `SELECT * FROM apartments WHERE id=?`

	context, close := context.WithTimeout(as.context, contextTimeGetOneApartment*time.Second)
	defer close()

	if err = as.db.PingContext(context); err != nil {
		return
	}

	stmt, err := as.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	apartment = &entity.Apartment{}
	if err = stmt.QueryRowContext(context, id).Scan(&apartment.ID, &apartment.Title, &apartment.Price, &apartment.City, &apartment.Rooms, &apartment.Address, &apartment.Square, &apartment.IDRealtor, &apartment.UpdateTime, &apartment.CreateTime); err != nil {
		return
	}

	return apartment, nil
}

func (as *apartmentStorage) Create(apartment *entity.Apartment) (id int64, err error) {

	q := `INSERT INTO apartments (title, price, city, rooms, address, square, id_realtor, update_time, create_time) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	context, close := context.WithTimeout(as.context, contextTimeCreateApartment*time.Second)
	defer close()

	if err = as.db.PingContext(context); err != nil {
		return
	}

	stmt, err := as.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	row, err := stmt.ExecContext(context, apartment.Title, apartment.Price, apartment.City, apartment.Rooms, apartment.Address, apartment.Square, apartment.IDRealtor, apartment.UpdateTime, apartment.CreateTime)
	if err != nil {
		return
	}

	id, err = row.LastInsertId()

	return id, err
}

func (as *apartmentStorage) Update(apartment *entity.Apartment) (aff int64, err error) {

	q := `UPDATE apartments SET title=?, price=?, city=?, rooms=?, address=?, square=?, id_realtor=?, update_time=?, create_time=? WHERE id=?`

	context, close := context.WithTimeout(as.context, contextTimeUpdateApartment*time.Second)
	defer close()

	if err = as.db.PingContext(context); err != nil {
		return aff, err
	}

	stmt, err := as.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	result, err := stmt.ExecContext(context, apartment.Title, apartment.Price, apartment.City, apartment.Rooms, apartment.Address, apartment.Square, apartment.IDRealtor, apartment.UpdateTime, apartment.CreateTime, apartment.ID)
	if err != nil {
		return
	}

	aff, err = result.RowsAffected()

	return aff, err
}

func (as *apartmentStorage) Delete(id int) error {

	q := `DELETE FROM apartments WHERE id=?`

	context, close := context.WithTimeout(as.context, contextTimeDeleteApartment*time.Second)
	defer close()

	if err := as.db.PingContext(context); err != nil {
		return err
	}

	stmt, err := as.db.PrepareContext(context, q)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(context, "", id)

	return err
}
