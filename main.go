package main

import (
	"database/sql"
	"fmt"
	"regexp"

	//import runs init function which is necessary to register this db driver
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "gophercise_db"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = createDB(db, dbname)
	if err != nil {
		panic(err)
	}
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(db.Ping())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func dropDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE " + name)
	return err
}

func resetdb(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
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
