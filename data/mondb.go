package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func GetMaxId(dbname, schema_name, table_name string) (max_num int) {

	creds := GetCredentials("./data/db_settings.json")

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.User, creds.Password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// close database
	defer db.Close()

	selectStmt := "select max(Id) from " + schema_name + "." + table_name
	row := db.QueryRow(selectStmt)
	if err := row.Scan(&max_num); err != nil {
		if err == sql.ErrNoRows {
			max_num = 1
		} else {
			panic(err)
		}
	}
	return max_num
}
