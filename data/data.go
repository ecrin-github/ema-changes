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

var (
	g        *custom_logger.Logger // reference to logger
	settings string                // string for relative path of settings file
)

func init() {
	g = custom_logger.GetInstance(`C:\MDR_Logs\ema changes.txt`)
	settings = "./data/db_settings.json"
}

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
}

func GetConnectionString(db_name string) string {

	content, err := ioutil.ReadFile(settings)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}
	var c Credentials
	err = json.Unmarshal(content, &c)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, db_name)
}

func ProcessFileIDData(ids []string, date_string string) (num_updated, num_added int) {

	// Get connection string and open database
	psqlconn := GetConnectionString("mon")
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}

	// close database
	defer db.Close()

	// set up sql statements
	updateStmt := `update sf.source_data_studies set last_revised = $1, download_status = 0 where sd_id = $2`
	selectStmt := `select id from sf.source_data_studies where sd_id = $1`
	insertStmt := `insert into sf.source_data_studies 
	 				(source_id, sd_id, remote_url, last_revised, download_status) 
					 values (100123, $1, $2, $3, 0)`

	// do the updates
	for _, t := range ids {

		g.InfoLogger.Printf("  id in file: %s\n", t)
		id := t[:14]

		link_country := t[15:]
		if len(link_country) > 2 && link_country[:3] == "Out" {
			link_country = "3rd"
		}
		link_id := t[:14] + "/" + link_country
		remote_link := "https://www.clinicaltrialsregister.eu/ctr-search/trial/" + link_id

		// check if this record exists
		var table_id int
		err := db.QueryRow(selectStmt, id).Scan(&table_id)
		if err != nil {
			if err == sql.ErrNoRows {
				// no record - needs to be inserted
				_, err = db.Exec(insertStmt, id, remote_link, date_string)
				if err != nil {
					g.ErrorLogger.Println(err)
					os.Exit(1)
				}
				num_added++
			} else {
				g.ErrorLogger.Println(err)
				os.Exit(1)
			}
		} else {
			// record present - needs to be updated
			_, err = db.Exec(updateStmt, date_string, id)
			if err != nil {
				g.ErrorLogger.Println(err)
				os.Exit(1)
			}
			num_updated++
		}
	}
	return
}
