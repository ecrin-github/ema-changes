package data

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func GetMaxId(dbname, schema_name, table_name string) (max_num int) {

	// Get connection string and open database
	psqlconn := GetConnectionString(dbname)
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

func StoreSAFRecord(file_path string, start_time time.Time, num_updated, num_added int) {

	// Get connection string and open database
	psqlconn := GetConnectionString("mon")
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer db.Close()

	id := GetMaxId("mon", "sf", "saf_events") + 1
	g.InfoLogger.Printf("New Id in saf evets table: %d", id)

	end_time := time.Now()
	insertStmt := `insert into sf.saf_events
	 				(id, source_id, type_id, time_started, time_ended, num_records_checked, num_records_added, num_records_downloaded, comments) 
					values ($1, 100123, 305, $2, $3, $4, $5, 0, $6)`
	g.InfoLogger.Printf("\n%s\n", insertStmt)
	_, err = db.Exec(insertStmt, id, start_time, end_time, num_updated+num_added, num_added, "Source file: "+file_path)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}
}
