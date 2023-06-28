package main

import (
	"encoding/csv"
    "fmt"
    "log"
    "os"

	"github.com/USACE/go-consequences/compute"
	"github.com/USACE/go-consequences/hazardproviders"
	"github.com/USACE/go-consequences/resultswriters"
	"github.com/USACE/go-consequences/structureprovider"
)

func readCsvFile(filePath string) [][]string {
    f, err := os.Open(filePath)
    if err != nil {
        log.Fatal("Unable to read input file " + filePath, err)
    }
    defer f.Close()

    csvReader := csv.NewReader(f)
    records, err := csvReader.ReadAll()
    if err != nil {
        log.Fatal("Unable to parse file as CSV for " + filePath, err)
    }

    return records
}

func main() {
	fmt.Println("Start")

	// Read CSV file containing links to each simulation's data
	records := readCsvFile("run_table.csv")

	// Looping each record row for each simulation run. 
	// Skipping first row of headers.
	for i, s := range records[1:] {
		var sim_title string = s[0]
		var inventory_file string = s[2]
		var hazard_file string = s[3]
		var output_file string = s[4]

		fmt.Println(i, sim_title)

		// Define flood hazard, provide the file path to your desired flood
		hazard, e := hazardproviders.Init(hazard_file)
		// hazard, e := hazardproviders.Init("./data/Amite_Katrina2005_AORC_ADCIRC_2021Geometry.tif")
		if e != nil {
			panic(e)
		}
		defer hazard.Close()

		// // Define structure inventory, provide file path to structure inventory
		inventory, b := structureprovider.InitSHP(inventory_file)
		// inventory, b := structureprovider.InitSHP(inv_str)
		if b != nil {
			panic(b)
		}
		// fmt.Println(inventory)

		// // Define result file, provide file path where you want results to be saved
		results, c := resultswriters.InitGpkResultsWriter_Projected(output_file, "event_results", 9822)
		if c != nil {
			panic(c)
		}
		defer results.Close()

		// // Run the compute
		compute.StreamAbstract(hazard, inventory, results)
		fmt.Println("End")
	}

}
