package main

import (
	"encoding/xml"
	"os"
	"regexp"

	"github.com/scanhamman/ema-changes/custom_logger"
	"github.com/scanhamman/ema-changes/data"
)

type trials struct {
	XMLName xml.Name `xml:"trials"`
	Trials  []trial  `xml:"trial"`
}

type trial struct {
	XMLName xml.Name `xml:"trial"`
	Main    MainData `xml:"main"`
}

type MainData struct {
	XMLName     xml.Name `xml:"main"`
	TrialId     string   `xml:"trial_id"`
	PublicTitle string   `xml:"public_title"`
}

var g *custom_logger.Logger

func init() {
	// Set up global logger, to write to a file

	g = custom_logger.GetInstance(`C:\MDR_Logs\ema changes.txt`)
}

func main() {
	fileName := os.Args[1]

	// Check filename has the right pattern
	_, err := regexp.MatchString("^ema [0-9]{8}.xml$", fileName)
	if err != nil {
		os.Exit(1)
	}

	// get date embedded in filename in ISO format.
	// It can then be used in database statements.
	date_string := fileName[4:8] + "-" + fileName[8:10] + "-" + fileName[10:12]
	g.InfoLogger.Printf("date string is : %s\n", date_string)

	// Build the location of the ema file.
	// filepath.Abs appends the file name to the default working directory.
	basePath := `C:\MDR_Sources\EUCTR\`
	trialsFilePath := basePath + fileName
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}

	// Open the ema file, with deferred closure.
	file, err := os.Open(trialsFilePath)
	if err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Read in XML file and decodes to defined structures.
	var foundTrials trials
	if err := xml.NewDecoder(file).Decode(&foundTrials); err != nil {
		g.ErrorLogger.Println(err)
		os.Exit(1)
	}

	// Obtain listed id for each trial
	// and add to a slice of strings.
	var ids []string
	for _, t := range foundTrials.Trials {
		m := t.Main
		ids = append(ids, m.TrialId)
	}

	// add these trial ids for testing purposes
	ids = append(ids, "2004-000007-18-SE")
	ids = append(ids, "2004-000012-13-CZ")
	ids = append(ids, "2004-000015-25-SK")

	// List the trials for checking purposes
	for _, id := range ids {
		g.InfoLogger.Printf("id: %s\n", id)
	}

	// Update source studies table in database.
	num_updated, num_added := data.ProcessFileIDData(ids, date_string)

	max_id := data.GetMaxId("mon", "sf", "saf_events")
	g.InfoLogger.Printf("Max Id in saf evets table: %d", max_id)
	g.InfoLogger.Printf("Number of records updated: %d", num_updated)
	g.InfoLogger.Printf("Num of records added: %d", num_added)
}
