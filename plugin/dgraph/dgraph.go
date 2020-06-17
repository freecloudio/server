package dgraph

import (
	"context"
	"fmt"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/freecloudio/server/plugin/dgraph/schema"
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
	//TODO: Get vals from config
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
	//TODO: Improve schema concatenation
	schema := schema.Common + "\n" + schema.User
	return dg.Alter(context.TODO(), &api.Operation{Schema: schema})
}
