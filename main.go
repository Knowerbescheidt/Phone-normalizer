package main

import (
	"fmt"
	"regexp"

	phonedb "github.com/Knowerbescheidt/Phone-normalizer/db"
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
	//reset db
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	must(phonedb.Reset("postgres", psqlInfo, dbname))

	//migrate db
	must(phonedb.Migrate("postgres", psqlInfo))

	//open connection to db
	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	//insert sample data
	if err := db.Seed(); err != nil {
		panic(err)
	}

	//run normalizer
	phones, err := db.AllPhones()
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Value)
		if number != p.Value {
			fmt.Println("Updating or removing...", number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.Id))
			} else {
				p.Value = number
				err := db.UpdatePhone(&p)
				must(err)
			}

		} else {
			fmt.Println("No changes required!")
		}
	}
}

//helper function
func must(err error) {
	if err != nil {
		panic(err)
	}
}

//normalize function
func normalize(tn string) string {
	re := regexp.MustCompile("[^0-9]")
	return re.ReplaceAllString(tn, "")
}
