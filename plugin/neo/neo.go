package neo

import (
	"github.com/freecloudio/server/application/persistence"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func init() {
	persistence.RegisterPluginInitialization("neo", InitializeNeo)
}

// InitializeNeo connects this plugin to the neo4j database
func InitializeNeo() (err error) {
	configForNeo4j40 := func(conf *neo4j.Config) {
		conf.Encrypted = false
	}
	_, err = neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""), configForNeo4j40)
	if err != nil {
		return err
	}
	return
}
