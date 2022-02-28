package main

import (
	"database/sql"
	"fmt"
	"regexp"

	//import runs init function which is necessary to register this db driver
	phonedb "github.com/Knowerbescheidt/Phone-normalizer/db"

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
	//reset db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	//minute 7
	must(phonedb.Migrate("postgresql", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	if err := db.Seed(); err != nil {
		panic(err)
	}

	phones, err := db.AllPhones()
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Value)
		if number != p.Value {
			fmt.Println("Updating or removing...", number)
			// existing, err := findPhone(db, number)
			// must(err)
			// if existing != nil {
			// 	must(deletePhone(db, p.Id))
			// } else {
			// 	p.Value = number
			// 	err := updatePhone(db, p)
			// 	must(err)
			// }

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

// func updatePhone(db *sql.DB, p phone) error {
// 	statement := "UPDATE phone_numbers SET value=$2 WHERE id=$1"
// 	_, err := db.Exec(statement, p.Id, p.Value)
// 	return err
// }

func deletePhone(db *sql.DB, id int) error {
	statement := "DELETE FROM phone_numbers WHERE id=$1"
	_, err := db.Exec(statement, id)
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
