package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"reflect"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v2"
)

type Contents struct {
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	SourceId      int       `json:"source_id"`
	DatePublished time.Time `json:"date_published,omitempty"`
	DateAdded     time.Time `json:"date_added,omitempty"`
}

type Words []string

func (s *Words) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *Words) Set(v string) error {
	*s = append(*s, v)
	return nil
}

var urlParam = flag.String("url", "http://localhost:9200", "elasticsearch server url")
var indexParam = flag.String("index", "", "elasticsearch index")
var words Words

func NewElasticClient(url string) (*elastic.Client, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"url": url,
		}).Error("Failed to create elastic client")
	}
	return client, err
}

func PrintSearchResultEash(searchResult *elastic.SearchResult) {
	// optimistic way to consume search result. it ignores errors in serialization.
	log.Info("Start search result each")
	var contents Contents
	for _, item := range searchResult.Each(reflect.TypeOf(contents)) {
		if c, ok := item.(Contents); ok {
			log.WithFields(log.Fields{
				"title":          c.Title,
				"body":           c.Body,
				"source_id":      c.SourceId,
				"date_published": c.DatePublished,
				"date_added":     c.DateAdded,
			}).Info("Use search result eash")
		}
	}
	log.Info("End search result each")
}

func PrintSearchResultHits(searchResult *elastic.SearchResult) {
	// handle serialization by myself
	if searchResult.Hits != nil {
		log.WithFields(log.Fields{
			"TotalHits": searchResult.Hits.TotalHits,
		}).Info("Search Result Hits")

		for _, hit := range searchResult.Hits.Hits {
			var c Contents
			err := json.Unmarshal(*hit.Source, &c)
			if err != nil {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Failed to decode as json")
			}

			log.WithFields(log.Fields{
				"title":          c.Title,
				"body":           c.Body,
				"source_id":      c.SourceId,
				"date_published": c.DatePublished,
				"date_added":     c.DateAdded,
			}).Info("Use search result hits")
		}
	} else {
		log.Info("No search result hits")
	}
}

func main() {
	flag.Var(&words, "word", "search word")
	flag.Parse()

	url := strings.ToLower(*urlParam)
	index := strings.ToLower(*indexParam)

	if index == "" {
		log.Error("Specify index name")
		return
	}

	log.WithFields(log.Fields{
		"url":   url,
		"index": index,
		"words": words,
	}).Info("Command line parameters")

	client, err := NewElasticClient(url)
	if err != nil {
		return
	}

	wordsStr := words.String()
	query := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("title", wordsStr),
		elastic.NewMatchQuery("body", wordsStr),
	)
	searchResult, err := client.Search().
		Index(index).
		Query(&query).
		From(0).Size(10).
		Do()

	if err != nil {
		log.WithFields(log.Fields{
			"err":         err,
			"word string": wordsStr,
		}).Error("Failed to search")
		return
	}

	log.WithFields(log.Fields{
		"TotalHits":         searchResult.TotalHits(),
		"TookInMillis (ms)": searchResult.TookInMillis,
	}).Info("Search Result")

	PrintSearchResultEash(searchResult)
	PrintSearchResultHits(searchResult)
}
