package adapterSql

import (
	"context"
	"database/sql"
	"time"

	"gilab.com/estate-agency-api/internal/entity"
)

const (
	contextTimeGetAllRealtor = 2
	contextTimeGetOneRealtor = 1
	contextTimeCreateRealtor = 1
	contextTimeUpdateRealtor = 1
	contextTimeDeleteRealtor = 1
)

type realtorAdapter struct {
	db      *sql.DB
	context context.Context
}

func NewRealtorAdapter(db *sql.DB) *realtorAdapter {
	return &realtorAdapter{db: db, context: context.Background()}
}

func (rs *realtorAdapter) GetAll(page int, pageSize int) (realtors []*entity.Realtor, err error) {

	q := `SELECT * FROM realtors LIMIT ?,?`

	context, close := context.WithTimeout(rs.context, contextTimeGetAllRealtor*time.Second)
	defer close()

	if err = rs.db.PingContext(context); err != nil {
		return
	}

	stmt, err := rs.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	rows, err := stmt.QueryContext(context, page*pageSize, pageSize)
	if err != nil {
		return
	}

	for rows.Next() {
		realtor := &entity.Realtor{}
		err = rows.Scan(&realtor.ID, &realtor.FirstName, &realtor.LastName, &realtor.Phone, &realtor.Email, &realtor.Rating, &realtor.Experience)
		if err != nil {
			return
		}
		realtors = append(realtors, realtor)
	}

	rows.Close()

	return realtors, nil
}

func (rs *realtorAdapter) GetByID(id int) (realtor *entity.Realtor, err error) {

	q := `SELECT * FROM realtors WHERE id=?`

	context, close := context.WithTimeout(rs.context, contextTimeGetOneRealtor*time.Second)
	defer close()

	if err = rs.db.PingContext(context); err != nil {
		return
	}

	realtor = &entity.Realtor{}
	stmt, err := rs.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	if err = stmt.QueryRowContext(context, id).Scan(&realtor.ID, &realtor.FirstName, &realtor.LastName, &realtor.Phone, &realtor.Email, &realtor.Rating, &realtor.Experience); err != nil {
		return
	}

	return realtor, nil
}

func (rs *realtorAdapter) Create(realtor *entity.Realtor) (id int64, err error) {

	q := `INSERT INTO realtors (first_name, last_name, phone, email, rating, experience) VALUES (?, ?, ?, ?, ?, ?)`

	context, close := context.WithTimeout(rs.context, contextTimeCreateRealtor*time.Second)
	defer close()

	if err = rs.db.PingContext(context); err != nil {
		return
	}

	stmt, err := rs.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	row, err := stmt.ExecContext(context, realtor.FirstName, realtor.LastName, realtor.Phone, realtor.Email, realtor.Rating, realtor.Experience)
	if err != nil {
		return
	}

	id, err = row.LastInsertId()

	return id, err
}

func (rs *realtorAdapter) Update(realtor *entity.Realtor) (aff int64, err error) {

	q := `UPDATE realtors SET first_name=?, last_name=?, phone=?, email=?, rating=?, experience=? WHERE id=?`

	context, close := context.WithTimeout(rs.context, contextTimeUpdateRealtor*time.Second)
	defer close()

	if err = rs.db.PingContext(context); err != nil {
		return
	}

	stmt, err := rs.db.PrepareContext(context, q)
	if err != nil {
		return
	}

	result, err := stmt.ExecContext(context, realtor.FirstName, realtor.LastName, realtor.Phone, realtor.Email, realtor.Rating, realtor.Experience, realtor.ID)
	if err != nil {
		return
	}

	aff, err = result.RowsAffected()

	return aff, err
}

func (rs *realtorAdapter) Delete(id int) error {

	q := `DELETE FROM realtors WHERE id=?`

	context, close := context.WithTimeout(rs.context, contextTimeDeleteRealtor*time.Second)
	defer close()

	if err := rs.db.PingContext(context); err != nil {
		return err
	}

	stmt, err := rs.db.PrepareContext(context, q)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(context, id)

	return err
}
