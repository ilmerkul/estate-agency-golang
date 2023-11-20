package mysql

import (
	"context"
	"database/sql"
	"time"

	"gilab.com/estate-agency-api/internal/domain/entity"
)

type realtorStorage struct {
	db *sql.DB
}

func NewRealtorStorage(db *sql.DB) *realtorStorage {
	return &realtorStorage{db: db}
}

func (rs *realtorStorage) GetAll(page int, pageSize int) (realtors []*entity.Realtor, err error) {
	context, close := context.WithTimeout(context.Background(), 3*time.Second)
	defer close()

	if err := rs.db.PingContext(context); err != nil {
		return nil, err
	}

	rows, err := rs.db.QueryContext(context, "SELECT * FROM realtors LIMIT ?,?", page*pageSize, pageSize)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		realtor := &entity.Realtor{}
		err = rows.Scan(&realtor.ID, &realtor.FirstName, &realtor.LastName, &realtor.Phone, &realtor.Email, &realtor.Rating, &realtor.Experience)
		if err != nil {
			return nil, err
		}
		realtors = append(realtors, realtor)
	}

	rows.Close()

	return realtors, nil
}

func (rs *realtorStorage) GetByID(id int) (realtor *entity.Realtor, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := rs.db.PingContext(context); err != nil {
		return realtor, err
	}

	realtor = &entity.Realtor{}
	err = rs.db.QueryRowContext(context, "SELECT * FROM realtors WHERE id=?", id).Scan(&realtor.ID, &realtor.FirstName, &realtor.LastName, &realtor.Phone, &realtor.Email, &realtor.Rating, &realtor.Experience)
	if err != nil {
		return realtor, err
	}

	return realtor, nil
}

func (rs *realtorStorage) Create(realtor *entity.Realtor) (id int64, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := rs.db.PingContext(context); err != nil {
		return id, err
	}

	row, err := rs.db.ExecContext(context, "INSERT INTO realtors (first_name, last_name, phone, email, rating, experience) VALUES (?, ?, ?, ?, ?, ?)", realtor.FirstName, realtor.LastName, realtor.Phone, realtor.Email, realtor.Rating, realtor.Experience)
	if err != nil {
		return id, err
	}

	id, err = row.LastInsertId()

	return id, err
}

func (rs *realtorStorage) Update(realtor *entity.Realtor) (aff int64, err error) {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err = rs.db.PingContext(context); err != nil {
		return aff, err
	}

	result, err := rs.db.ExecContext(context, "UPDATE realtors SET first_name=?, last_name=?, phone=?, email=?, rating=?, experience=? WHERE id=?", realtor.FirstName, realtor.LastName, realtor.Phone, realtor.Email, realtor.Rating, realtor.Experience, realtor.ID)
	if err != nil {
		return aff, err
	}

	aff, err = result.RowsAffected()

	return aff, err
}

func (rs *realtorStorage) Delete(id int) error {
	context, close := context.WithTimeout(context.Background(), 1*time.Second)
	defer close()

	if err := rs.db.PingContext(context); err != nil {
		return err
	}

	_, err := rs.db.ExecContext(context, "DELETE FROM realtors WHERE id=?", id)
	if err != nil {
		return err
	}

	return nil
}
