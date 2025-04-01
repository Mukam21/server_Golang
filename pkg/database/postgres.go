package database

import (
	"database/sql"
	"fmt"

	"github.com/Mukam21/server_Golang/pkg/config"
	"github.com/Mukam21/server_Golang/pkg/model"

	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(person *model.Person) (int64, error) {
	query := `
        INSERT INTO persons (name, surname, patronymic, age, gender, nationality)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	var id int64
	err := r.db.QueryRow(query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) GetByID(id int64) (*model.Person, error) {
	query := `
        SELECT id, name, surname, patronymic, age, gender, nationality
        FROM persons
        WHERE id = $1`

	person := &model.Person{}
	err := r.db.QueryRow(query, id).Scan(
		&person.ID,
		&person.Name,
		&person.Surname,
		&person.Patronymic,
		&person.Age,
		&person.Gender,
		&person.Nationality,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return person, nil
}

func (r *Repository) GetAll(page, limit int, nameFilter string) ([]*model.Person, error) {
	offset := (page - 1) * limit
	query := `
        SELECT id, name, surname, patronymic, age, gender, nationality
        FROM persons
        WHERE name ILIKE $1
        LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, "%"+nameFilter+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var persons []*model.Person
	for rows.Next() {
		person := &model.Person{}
		err := rows.Scan(
			&person.ID,
			&person.Name,
			&person.Surname,
			&person.Patronymic,
			&person.Age,
			&person.Gender,
			&person.Nationality,
		)
		if err != nil {
			return nil, err
		}
		persons = append(persons, person)
	}

	return persons, nil
}

func (r *Repository) Update(person *model.Person) error {
	query := `
        UPDATE persons
        SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6
        WHERE id = $7`

	result, err := r.db.Exec(query,
		person.Name,
		person.Surname,
		person.Patronymic,
		person.Age,
		person.Gender,
		person.Nationality,
		person.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	query := `DELETE FROM persons WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
