package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
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

func selectRows(session *gocql.Session, cql string) {
	iter := session.Query(cql).Iter()
	defer iter.Close()

	rows, err := iter.SliceMap()
	if err != nil {
		fmt.Printf("failed to query: %+v\n", err)
		return
	}

	for _, row := range rows {
		fmt.Printf("row: %+v\n", row)
	}
}

var (
	CHOST   = os.Getenv("CASSANDRA_HOST")
	CPORT   = os.Getenv("CASSANDRA_PORT")
	CUSER   = os.Getenv("CASSANDRA_USER")
	CPASS   = os.Getenv("CASSANDRA_PASSWORD")
	CCAPATH = os.Getenv("CASSANDRA_CA_PATH")
)

func getSslOptions() (opts *gocql.SslOptions) {
	config := &tls.Config{
		ServerName:         CHOST,
		InsecureSkipVerify: false,
	}
	opts = &gocql.SslOptions{
		Config:                 config,
		EnableHostVerification: true,
		CaPath:                 CCAPATH,
	}
	return
}

func getClusterConfig() (cluster *gocql.ClusterConfig) {
	port, _ := strconv.Atoi(CPORT)
	cluster = gocql.NewCluster(CHOST)
	cluster.CQLVersion = "3.4.4"
	cluster.Port = port
	cluster.Consistency = gocql.LocalOne
	cluster.SerialConsistency = gocql.LocalSerial
	cluster.ProtoVersion = 4
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: CUSER,
		Password: CPASS,
	}
	cluster.Timeout = 3 * time.Second
	cluster.ConnectTimeout = 3 * time.Second

	if CCAPATH != "" {
		cluster.SslOpts = getSslOptions()
	}

	// cluster.Keyspace = "..."
	return cluster
}

var cql = flag.String("cql", "", "specify cql statement")

func main() {
	flag.Parse()

	cluster := getClusterConfig()
	session, err := cluster.CreateSession()
	if err != nil {
		fmt.Printf("failed to create session: %+v\n", err)
		return
	}
	defer session.Close()

	if *cql != "" {
		fmt.Printf("%+v\n", *cql)
		selectRows(session, *cql)
	}
}
