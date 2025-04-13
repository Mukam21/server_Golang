package database

import (
	"database/sql"
	"fmt"
	"strings"

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

func (r *Repository) GetAll(page, limit int, filters map[string]string) ([]*model.Person, error) {
	offset := (page - 1) * limit
	var conditions []string
	var args []interface{}
	argIndex := 1

	if name := filters["name"]; name != "" {
		conditions = append(conditions, fmt.Sprintf("name ILIKE $%d", argIndex))
		args = append(args, "%"+name+"%")
		argIndex++
	}
	if surname := filters["surname"]; surname != "" {
		conditions = append(conditions, fmt.Sprintf("surname ILIKE $%d", argIndex))
		args = append(args, "%"+surname+"%")
		argIndex++
	}
	if age := filters["age"]; age != "" {
		conditions = append(conditions, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, age)
		argIndex++
	}
	if gender := filters["gender"]; gender != "" {
		conditions = append(conditions, fmt.Sprintf("gender = $%d", argIndex))
		args = append(args, gender)
		argIndex++
	}
	if nationality := filters["nationality"]; nationality != "" {
		conditions = append(conditions, fmt.Sprintf("nationality = $%d", argIndex))
		args = append(args, nationality)
		argIndex++
	}

	query := "SELECT id, name, surname, patronymic, age, gender, nationality FROM persons"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
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

func (r *Repository) Patch(id int64, patch *model.PersonPatchRequest) error {
	var updates []string
	var args []interface{}
	argIndex := 1

	if patch.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *patch.Name)
		argIndex++
	}
	if patch.Surname != nil {
		updates = append(updates, fmt.Sprintf("surname = $%d", argIndex))
		args = append(args, *patch.Surname)
		argIndex++
	}
	if patch.Patronymic != nil {
		updates = append(updates, fmt.Sprintf("patronymic = $%d", argIndex))
		args = append(args, *patch.Patronymic)
		argIndex++
	}
	if patch.Age != nil {
		updates = append(updates, fmt.Sprintf("age = $%d", argIndex))
		args = append(args, *patch.Age)
		argIndex++
	}
	if patch.Gender != nil {
		updates = append(updates, fmt.Sprintf("gender = $%d", argIndex))
		args = append(args, *patch.Gender)
		argIndex++
	}
	if patch.Nationality != nil {
		updates = append(updates, fmt.Sprintf("nationality = $%d", argIndex))
		args = append(args, *patch.Nationality)
		argIndex++
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE persons SET %s WHERE id = $%d",
		strings.Join(updates, ", "), argIndex)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
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
