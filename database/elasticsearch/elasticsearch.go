package database

import (
	"encoding/json"
	"fmt"
	"reflect"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	elastic "gopkg.in/olivere/elastic.v3"
)

var (
	name     string
	category string
)

type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// PluginResults plugin results
type PluginResults struct {
	ID   string `json:"id" gorethink:"id,omitempty"`
	Data map[string]interface{}
}

// Create a new index.
const mapping = `{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "tweet":{
            "properties":{
                "tags":{
                    "type":"string"
                },
                "location":{
                    "type":"geo_point"
                },
                "suggest_field":{
                    "type":"completion",
                    "payloads":true
                }
            }
        }
    }
}`

// WriteToDatabase upserts plugin results into Database
func WriteToDatabase(results map[string]interface{}) {

	// Create a client
	elasitcURL := fmt.Sprintf("%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elasticsearch"))
	client, err := elastic.NewClient(
		elastic.SetURL(elasitcURL),
		elastic.SetMaxRetries(10))
	utils.Assert(err)

	// Create an index
	createIndex, err := client.CreateIndex("malice").Do()
	utils.Assert(err)
	if !createIndex.Acknowledged {
		// Not acknowledged
		log.Error("Couldn't create Index.")
	}

	// Add a document to the index
	tweet := Tweet{User: "olivere", Message: "Take Five"}
	_, err = client.Index().
		Index("twitter").
		Type("tweet").
		Id("1").
		BodyJson(tweet).
		Refresh(true).
		Do()
	utils.Assert(err)

	// Search with a term query
	termQuery := elastic.NewTermQuery("user", "olivere")
	searchResult, err := client.Search().
		Index("twitter").   // search in index "twitter"
		Query(termQuery).   // specify the query
		Sort("user", true). // sort by "user" field, ascending
		From(0).Size(10).   // take documents 0-9
		Pretty(true).       // pretty print request and response JSON
		Do()                // execute
	utils.Assert(err)

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t Tweet
			utils.Assert(json.Unmarshal(*hit.Source, &t))

			// Work with tweet
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Delete the index again
	_, err = client.DeleteIndex("twitter").Do()
	utils.Assert(err)
}
