package create

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/vgo0/gotimepad/db"
	"github.com/vgo0/gotimepad/models"
	"github.com/vgo0/gotimepad/random"
	h "github.com/vgo0/gotimepad/handler"
)

func Exec(args []string) {
	createCmd := flag.NewFlagSet("init", flag.ExitOnError)
	createDB := createCmd.String("d", "gotimepad.gtp", "One time pad database file")
	createPages := createCmd.Uint("p", 500, "Number of pages")
	createPageSize := createCmd.Uint("s", 2000, "Page size (max number of characters per page)")
	createForce := createCmd.Bool("f", false, "Force create (overwrite existing file)")
	createCmd.Parse(args)

	err := checkIfExists(createDB, createForce)
	h.CheckError(err)

	db.Connect(*createDB)
	db.Migrate()
	
	err = generatePages(createPages, createPageSize)
	h.CheckError(err)
}

/*
Generates the one time pages based on provided inputs or defaults

Returns error if issue encountered
*/
func generatePages(createPages *uint, createPageSize *uint) error {
	var pages []models.OTPage
	var i uint

	for i = 0; i < *createPages; i++ {
		data, err := random.GenerateRandomBytes(*createPageSize)

		if err != nil {
			return fmt.Errorf("error creating random data on iteration %d: %s", i, err)
		}

		page := models.OTPage{
			Page: data,
		}

		pages = append(pages, page)
	}

	res := db.DB.Create(pages)

	if res.Error != nil {
		return fmt.Errorf("error inserting pages into database: %s", res.Error)
	}

	return nil
}

// Makes sure we don't overwrite an existing file unless specifically asked for
func checkIfExists(createDB *string, createForce *bool) error {
	if _, err := os.Stat(*createDB); err == nil {
		if !*createForce {
			return fmt.Errorf(`specified file already exists, if you wish to continue specify the -f flag to force create`)
		} else {
			// Attempt to remove existing file in force create mode
			log.Printf("Specified file exists, attempting to remove: %s\n", *createDB)

			err := os.Remove(*createDB)
			if err != nil {
				return fmt.Errorf("unable to remove the existing file %s: %s", *createDB, err)
			}
		}
	}

	return nil
}