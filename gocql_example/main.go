package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

type PlayList struct {
	ID        gocql.UUID
	SongOrder int
	Album     string
	Artist    string
	Reviews   []string
	SongID    gocql.UUID
	Tags      []string
	Title     string
	Venue     map[time.Time]string
}

func (p *PlayList) FieldMap() map[string]interface{} {
	return map[string]interface{}{
		"id":         &p.ID,
		"song_order": &p.SongOrder,
		"album":      &p.Album,
		"artist":     &p.Artist,
		"reviews":    &p.Reviews,
		"song_id":    &p.SongID,
		"tags":       &p.Tags,
		"title":      &p.Title,
		"venue":      &p.Venue,
	}
}

type PlayLists []PlayList

type User struct {
	UserID    string
	Emails    []string
	FirstName string
	LastName  string
	Todo      map[time.Time]string
	TopPlaces []string
}

func (u *User) FieldMap() map[string]interface{} {
	return map[string]interface{}{
		"user_id":    &u.UserID,
		"emails":     &u.Emails,
		"first_name": &u.FirstName,
		"last_name":  &u.LastName,
		"todo":       &u.Todo,
		"top_places": &u.TopPlaces,
	}
}

type Users []User

func connect(host, keySpace string) *gocql.Session {
	cluster := gocql.NewCluster(host)
	cluster.Keyspace = keySpace
	cluster.Consistency = gocql.Quorum
	session, _ := cluster.CreateSession()
	return session
}

func selectPlayLists(session *gocql.Session) (playlists PlayLists) {
	var p PlayList

	q := `SELECT * FROM playlists`
	iter := session.Query(q).Iter()
	defer iter.Close()

	for iter.MapScan(p.FieldMap()) {
		playlists = append(playlists, p)
	}

	playListLength := len(playlists)
	if playListLength == 0 {
		log.Println("data not found")
		return
	}

	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	return
}

func selectUsersLength(session *gocql.Session) (count int, err error) {
	q := `SELECT count(*) FROM users`
	err = session.Query(q).Consistency(gocql.One).Scan(&count)
	return
}

func insertUser(session *gocql.Session, id, firstName, lastName string) (err error) {
	q := `INSERT INTO users (user_id, first_name, last_name) VALUES (?, ?, ?)`
	err = session.Query(q, id, firstName, lastName).Exec()
	return
}

func insertUserWithCounting(session *gocql.Session, userID string, names []string) {
	if cnt, err := selectUsersLength(session); err == nil {
		log.Println(fmt.Sprintf("before: users count is %d", cnt))
	}

	if err := insertUser(session, userID, names[0], names[1]); err != nil {
		log.Fatal(err)
	}
	log.Println(fmt.Sprintf("%s is inserted", userID))

	if cnt, err := selectUsersLength(session); err == nil {
		log.Println(fmt.Sprintf("after: users count is %d", cnt))
	}
}

var userID = flag.String("uid", "", "")
var name = flag.String("name", "", "firstName lastName")

func main() {
	flag.Parse()

	var names []string
	if *name != "" {
		names = strings.Split(*name, " ")
	}

	session := connect("127.0.0.1", "music")
	defer session.Close()

	// select
	playlists := selectPlayLists(session)
	for _, p := range playlists {
		log.Println(fmt.Sprintf("%s, %s, %s", p.ID, p.Title, p.Artist))
	}

	// insert
	if *userID != "" && len(names) != 0 {
		insertUserWithCounting(session, *userID, names)
	}
}
