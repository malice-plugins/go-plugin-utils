package database

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/maliceio/go-plugin-utils/utils"
	elastic "gopkg.in/olivere/elastic.v3"
)

var (
	name     string
	category string
)

// PluginResults plugin results
type PluginResults struct {
	ID   string `json:"id"`
	Data map[string]interface{}
}

// WriteToDatabase upserts plugin results into Database
func WriteToDatabase(results PluginResults) {

	// scanID := utils.Getopt("MALICE_SCANID", "")
	ElasticAddr := fmt.Sprintf("%s:9200", utils.Getopt("MALICE_ELASTICSEARCH", "elastic"))
	client, err := elastic.NewSimpleClient(elastic.SetURL(ElasticAddr))
	utils.Assert(err)

	getSample, err := client.Get().
		Index("malice").
		Type("samples").
		Id(results.ID).
		Do()

	fmt.Println(getSample)
	fmt.Println(err)
	if err != nil {

	}

	if getSample.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", getSample.Id, getSample.Version, getSample.Index, getSample.Type)
		updateScan := map[string]interface{}{
			"scan_date": time.Now().Format(time.RFC3339Nano),
			"plugins": map[string]interface{}{
				category: map[string]interface{}{
					name: results.Data,
				},
			},
		}
		update, err := client.Update().Index("malice").Type("samples").Id(getSample.Id).
			Doc(updateScan).
			Do()
		utils.Assert(err)

		log.Debugf("New version of sample %q is now %d\n", update.Id, update.Version)
		// return *update

	} else {

		scan := map[string]interface{}{
			// "id":      sample.SHA256,
			// "file":      sample,
			"plugins": map[string]interface{}{
				category: map[string]interface{}{
					name: results.Data,
				},
			},
			"scan_date": time.Now().Format(time.RFC3339Nano),
		}

		newScan, err := client.Index().
			Index("malice").
			Type("samples").
			OpType("create").
			// Id("1").
			BodyJson(scan).
			Do()
		utils.Assert(err)

		log.Debugf("Indexed sample %s to index %s, type %s\n", newScan.Id, newScan.Index, newScan.Type)
		// return *newScan
	}
}
