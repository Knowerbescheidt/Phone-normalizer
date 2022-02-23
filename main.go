package main

import (
	"database/sql"
	"fmt"
	"regexp"

	//import runs init function which is necessary to register this db driver
	"db"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "gophercise_db"
)

//11:26 inserting records

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(db.Reset("postgres", psqlInfo, dbname))

	must(createPHTable(db))

	_, err = insertPhoneN(db, "1234567890")
	must(err)
	_, err = insertPhoneN(db, "123 456 7891")
	must(err)
	_, err = insertPhoneN(db, "(123) 456 7892")
	must(err)
	_, err = insertPhoneN(db, "(999) 000-1234")
	must(err)
	_, err = insertPhoneN(db, "123-456-7894")
	must(err)
	_, err = insertPhoneN(db, "123-456-7890")
	must(err)
	_, err = insertPhoneN(db, "1234567892")
	must(err)
	_, err = insertPhoneN(db, "(123)456-7892)")
	must(err)

	//phone, err := getPhone(db, id)
	//must(err)
	//fmt.Printf("we found this phone number in the db %s", phone)
	phones, err := getallPhones(db)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Value)
		if number != p.Value {
			fmt.Println("Updating or removing...", number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				must(deletePhone(db, p.Id))
			} else {
				p.Value = number
				err := updatePhone(db, p)
				must(err)
			}

		} else {
			fmt.Println("No changes required!")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func dropDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE " + name)
	return err
}

func getPhone(db *sql.DB, id int) (string, error) {
	statement := "SELECT value FROM phone_numbers WHERE id=$1"
	var phoneNumber string
	err := db.QueryRow(statement, id).Scan(&phoneNumber)
	if err != nil {
		return "", err
	}
	return phoneNumber, nil
}

func findPhone(db *sql.DB, number string) (*phone, error) {
	statement := "SELECT * FROM phone_numbers WHERE value=$1"
	var p phone
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

func updatePhone(db *sql.DB, p phone) error {
	statement := "UPDATE phone_numbers SET value=$2 WHERE id=$1"
	_, err := db.Exec(statement, p.Id, p.Value)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := "DELETE FROM phone_numbers WHERE id=$1"
	_, err := db.Exec(statement, id)
	return err
}

type phone struct {
	Id    int
	Value string
}

func getallPhones(db *sql.DB) ([]phone, error) {
	statement := "SELECT id, value FROM phone_numbers"
	rows, err := db.Query(statement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []phone
	for rows.Next() {
		var p phone
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

func createPHTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)`
	_, err := db.Exec(statement)
	return err
}

// my own working solution
// func normalize(tn string) string {
// 	r := strings.NewReplacer(" ", "", "(", "", ")", "", "-", "")
// 	tn = r.Replace(tn)
// 	return tn
// }

func normalize(tn string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(tn, "")

	// //bytes buffer is very strong in performance
	// var buf bytes.Buffer
	// for _, ch := range tn {
	// 	if ch >= '1' && ch <= '9' {
	// 		buf.WriteRune(ch)
	// 	}
	// }
	// return buf.String()
}
