package data

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "WinterIsComing!"
	dbname   = "mon"
)

func ProcessFileIDData(ids []string, date_string string) (num_updated, num_added int) {

	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// close database
	defer db.Close()

	// set up update statement
	updateStmt := `update sf.source_data_studies set last_revised = $1 where sd_id = $2`
	selectStmt := `select * from sf.source_data_studies where sd_id = $1`
	insertStmt := `insert into sf.source_data_studies 
	 				(source_id, sd_id, remote_url, last_revised, download_status) 
					 values (100123, $1, $2, $3, 0)`

	// do the updates
	for _, t := range ids {
		id := t[:14]

		// check if this record exists
		rows, err := db.Query(selectStmt, id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		defer rows.Close()

		if rows.Next() {
			// there is an existing record
			_, err = db.Exec(updateStmt, date_string, id)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			num_updated++
		} else {
			// need to add a new record
			link_country := t[15:]
			if len(link_country) > 2 && link_country[:3] == "Out" {
				link_country = "3rd"
			}
			link_id := t[:14] + "/" + link_country
			remote_link := "https://www.clinicaltrialsregister.eu/ctr-search/trial/" + link_id
			_, err = db.Exec(insertStmt, id, remote_link, date_string)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			num_added++
		}
	}
	return
}
