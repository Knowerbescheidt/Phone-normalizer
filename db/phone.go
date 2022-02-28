package db

import "database/sql"

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

type DB struct {
	db *sql.DB
}

func (db *DB) Close() error {
	return (db.db.Close())
}

func (db *DB) Seed() error {
	phones := []string{"1234567890", "123 456 7891", "(123) 456 7892", "(999) 000-1234", "123-456-7894", "123-456-7890", "1234567892", "(123)456-7892)"}
	for _, phone := range phones {
		if _, err := insertPhoneN(db.db, phone); err != nil {
			return err
		}
	}
	return nil
}

func Migrate(driverName, dataSource string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	createPhoneNumbersTable(db)
	if err != nil {
		return err
	}
	return db.Close()
}

func Reset(driverName, dataSource, dbName string) error {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return err
	}
	err = resetdb(db, dbName)
	if err != nil {
		return err
	}
	return db.Close()

}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func resetdb(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createPhoneNumbersTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

func insertPhoneN(db *sql.DB, phone string) (int, error) {

	// mit dem $1 umgeht ma sql injection
	statement := `INSERT INTO phone_numbers (value) VALUES($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	//psql returned nicht automatisch ie id wenn geinserted wird
	return id, nil
}

//Represents the phone_numbers
type Phone struct {
	Id    int
	Value string
}

func (db *DB) AllPhones() ([]Phone, error) {
	return allPhones(db.db)
}

func allPhones(db *sql.DB) ([]Phone, error) {
	statement := "SELECT id, value FROM phone_numbers"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []Phone
	for rows.Next() {
		var p Phone
		if err := rows.Scan(&p.Id, &p.Value); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil

}

func findPhone(db *sql.DB, number string) (*Phone, error) {
	statement := "SELECT * FROM phone_numbers WHERE value=$1"
	var p Phone
	row := db.QueryRow(statement, number)
	err := row.Scan(&p.Id, &p.Value)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}
