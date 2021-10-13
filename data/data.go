package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/scanhamman/ema-changes/custom_logger"

	_ "github.com/lib/pq"
)

var g *custom_logger.Logger

func init() {
	// Set up logger and file
	g = custom_logger.GetInstance(`C:\MDR_Logs\ema changes.txt`)
}

func ProcessFileIDData(ids []string, date_string string) (num_updated, num_added int) {

	creds := GetCredentials("./data/db_settings.json")
	dbname := "mon"
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.User, creds.Password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}

	// close database
	defer db.Close()

	// set up sql statements
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
			g.ErrorLogger.Println(err)
			os.Exit(1)
		}

		defer rows.Close()

		if rows.Next() {
			// there is an existing record
			_, err = db.Exec(updateStmt, date_string, id)
			if err != nil {
				g.ErrorLogger.Println(err)
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
				g.ErrorLogger.Println(err)
				os.Exit(1)
			}
			num_added++
		}
	}
	return
}

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
}

func GetCredentials(json_file string) (c Credentials) {

	content, err := ioutil.ReadFile(json_file)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(content, &c)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}
	return c
}
