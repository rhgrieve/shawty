package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"main/urls"
)

type Entry struct {
	URL    *urls.ShortURL `json:"url"`
	Visits int       `json:"visits"`
}

func (e *Entry) JSON() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}

	return string(bytes)
}

func (e *Entry) IncrementVisits() {
	e.Visits += 1
}

type DB struct {
	Filename string            `json:"filename"`
	Entries  map[string]*Entry `json:"entries"`
}

func NewDB(path string) *DB {
	return &DB{
		Filename: path,
		Entries:  make(map[string]*Entry),
	}
}

func (d *DB) Load() {
	bytes, err := ioutil.ReadFile(d.Filename)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(bytes, &d)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DB) Add(url *urls.ShortURL) {
	d.Entries[url.Shortcode] = &Entry{
		URL:    url,
		Visits: 0,
	}
}

func (d *DB) Get(shortUrl string) (*Entry, error) {
	record := d.Entries[shortUrl]

	if record.URL.IsEmpty() {
		return &Entry{}, errors.New("Not found")
	} else {
		return record, nil
	}
}

func (d *DB) Commit() {
	bytes, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(d.Filename, bytes, 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DB) Transaction(callback func()) {
	callback()
	d.Commit()
}

func (d *DB) JSON() string {
  bytes, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}

  return string(bytes)
}

func (d *DB) Dump() {
	fmt.Println(d.JSON())
}