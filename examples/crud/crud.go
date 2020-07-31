package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/jameycribbs/hare"
)

// IMPORTANT!
// Your record's struct MUST have an "ID" field and it must
// also implement the 3 methods below in order for it to
// satisfy the hare.Record interface.

type record struct {
	ID               int       `json:"id"`
	Season           int       `json:"season"`
	Episode          int       `json:"episode"`
	Film             string    `json:"film"`
	Shorts           []string  `json:"shorts"`
	YearFilmReleased int       `json:"year_film_released"`
	DateEpisodeAired time.Time `json:"date_episode_aired"`
	Host             string    `json:"host"`
}

func (r *record) GetID() int {
	return r.ID
}

func (r *record) SetID(id int) {
	r.ID = id
}

func (r *record) AfterFind() {
	*r = record(*r)
}

func main() {
	setupExampleDB()

	// Open the database and return a handle to it.
	db, err := hare.OpenDB("../example_data")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Let's grab a handle to the MST3K Episodes table so we can play with it.
	tbl, err := db.GetTable("mst3k_episodes")
	if err != nil {
		panic(err)
	}

	//----- CREATE -----

	// Here's how to create a new record.

	// To create a record, you pass a populated struct that satisfies the
	// hare.Record interface.  Do NOT populate the ID attribute!
	// It will be supplied by Hare when it creates the record.
	// You simply pass the new record to the Create method.  Hare will
	// add the record to the table and return the new record's ID.
	recID, err := tbl.Create(&record{
		Season:           6,
		Episode:          19,
		Film:             "Red Zone Cuba",
		Shorts:           []string{"Speech:  Platform, Posture, and Appearance"},
		YearFilmReleased: 1966,
		DateEpisodeAired: time.Date(1994, 12, 17, 0, 0, 0, 0, time.UTC),
		Host:             "Mike",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("New record id is:", recID)

	//----- READ -----

	// To read a specific record, you use the Find method.
	// You must know the record ID.

	// For more general querying, check out the example program in the
	// "querying" directory.
	rec := record{}

	err = tbl.Find(3, &rec)
	if err != nil {
		panic(err)
	}

	fmt.Println("Found record is:", rec.Film)

	//----- UPDATE -----

	// Here's how to update a record.

	// To update a record, you must have a populated record, including the
	// ID field.  You simply change the desired attributes, and then
	// pass the changed record to the Update method.

	// If you take a look inside ../data/mst3k_episodes.json after running
	// this example program, you will see a line of "X"'s right above the
	// line holding the line for id 4.  That is called a dummy record in
	// Hare, and it gets created when Hare has to update a record and the
	// changes have increased the record length.  Because the change we made
	// to the record with id 3 increased the record length, Hare dummied out
	// the existing version of record 3 and wrote the changed version of the
	// record at the bottom of the file.  Hare will attempt to re-use the
	// dummy record's space the next time it needs more space for a new or
	// updated record.  You can also run the database's "Compact" method
	// to remove all dummy records.

	rec.Film = "The Skydivers - The Final Cut"
	if err = tbl.Update(&rec); err != nil {
		panic(err)
	}

	//----- DELETE -----

	// To delete a record, you use the Destroy method.
	// You must know the record ID.

	// If you take a look inside ../data/mst3k_episodes.json after running
	// this example program, you will see that the first line of the file,
	// which was where the record with id 1 existed, has been replaced
	// with a dummy line.  This is how Hare deletes a record.  Hare will
	// attempt to re-use the dummy record's space the next time it needs
	// more space for a new or updated record.  You can also run the
	// database's "Compact" method to remove all dummy records.
	err = tbl.Destroy(1)
	if err != nil {
		panic(err)
	}
}

func setupExampleDB() {
	cmd := exec.Command("cp", "../example_data/mst3k_episodes_default.json", "../example_data/mst3k_episodes.json")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
