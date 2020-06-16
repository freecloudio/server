package dgraph

import (
	"context"
	"fmt"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
)

func init() {
	persistence.RegisterPluginInitialization(config.DGraphPersistenceKey, InitializeDGraph)
}

var dg *dgo.Dgraph

// InitializeDGraph connects this plugin to the dgraph database
func InitializeDGraph() (err error) {
	grpcClient, err := grpc.Dial("localhost:9080", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	if err != nil {
		err = fmt.Errorf("failed to initialize grpc connection to dgraph: %v", err)
		return
	}

	dg = dgo.NewDgraphClient(api.NewDgraphClient(grpcClient))

	err = setSchema()
	return
}

func setSchema() error {
	return dg.Alter(context.Background(), &api.Operation{
		Schema: `
			created: dateTime .
			updated: dateTime .

			first_name: string @index(hash) .
			last_name:  string @index(hash) .
			email:      string @index(hash) .
			password:   string .
			is_admin:   bool @index(bool) .

			type User {
				first_name
				last_name
				email
				password
				is_admin
				created
				updated
			}

			token:       string @index(hash) .
			valid_until: dateTime .
			for_user:    uid .

			type Token {
				token
				valid_until
				for_user
			}
		`,
	})
}
