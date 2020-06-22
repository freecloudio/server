package neo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/freecloudio/server/application/persistence"
	"github.com/freecloudio/server/config"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/sirupsen/logrus"
)

func init() {
	persistence.RegisterPluginInitialization(config.NeoPersistenceKey, InitializeNeo)
}

var neo neo4j.Driver

func InitializeNeo() (err error) {
	driver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "freecloud", ""), noEncrypted)
	if err != nil {
		return
	}
	neo = driver
	return
}

func noEncrypted(config *neo4j.Config) {
	config.Encrypted = false
}

type transactionCtx struct {
	session neo4j.Session
	neoTx   neo4j.Transaction
}

func (trCtx *transactionCtx) Commit() (err error) {
	txErr := trCtx.neoTx.Commit()
	if txErr != nil {
		logrus.WithError(txErr).Error("Failed to commit neo transaction - close session anyway")
		err = txErr
	}
	sessErr := trCtx.session.Close()
	if sessErr != nil {
		logrus.WithError(sessErr).Error("Failed to close neo session")
		if err == nil {
			err = sessErr
		}
	}
	return
}

func (trCtx *transactionCtx) Rollback() (err error) {
	txErr := trCtx.neoTx.Rollback()
	if txErr != nil {
		logrus.WithError(txErr).Error("Failed to rollback neo transaction - close session anyway")
		err = txErr
	}
	sessErr := trCtx.session.Close()
	if sessErr != nil {
		logrus.WithError(sessErr).Error("Failed to close neo session")
		if err == nil {
			err = sessErr
		}
	}
	return
}

func newTransactionContext(accessMode neo4j.AccessMode) (txCtx *transactionCtx, err error) {
	session, err := neo.NewSession(neo4j.SessionConfig{AccessMode: accessMode})
	if err != nil {
		err = fmt.Errorf("failed to create neo4j session: %v", err)
		return
	}
	neoTx, err := session.BeginTransaction()
	if err != nil {
		err = fmt.Errorf("failed to create neo4j transaction: %v", err)
		session.Close()
		return
	}
	txCtx = &transactionCtx{session, neoTx}
	return
}

// Convert given struct to a map with the 'fc_neo' / 'json' / field name as key and the field value as value
func modelToMap(model interface{}) map[string]interface{} {
	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()
	modelMap := make(map[string]interface{})

	for it := 0; it < modelValue.NumField(); it++ {
		valField := modelValue.Field(it)
		typeField := modelType.Field(it)

		dbName := getDBFieldName(typeField)
		if dbName == nil {
			continue
		}
		modelMap[*dbName] = valField.Interface()
	}

	return modelMap
}

func recordToModel(record neo4j.Record, key string, model interface{}) error {
	valInt, ok := record.Get(key)
	if !ok {
		return errors.New("value not found with key '" + key + "'")
	}
	valNode, ok := valInt.(neo4j.Node)
	if !ok {
		return errors.New("value with key '" + key + "' could not be converted to 'neo4j.Node'")
	}
	valProps := valNode.Props()

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for it := 0; it < modelValue.NumField(); it++ {
		valField := modelValue.Field(it)
		typeField := modelType.Field(it)

		dbNamePtr := getDBFieldName(typeField)
		if dbNamePtr == nil || !valField.CanSet() {
			continue
		}
		dbName := *dbNamePtr

		propInt, ok := valProps[dbName]
		if !ok {
			continue
		}
		propVal := reflect.ValueOf(propInt)
		valField.Set(propVal)
	}

	return nil
}

// Returns db field name based on tags of a struct field
// Returns nil if the field should not be stored in the db
// Uses own 'fc_neo' field tag but falls back to 'json' tags as these are automatically set from swagger
func getDBFieldName(typeField reflect.StructField) *string {
	var fieldTag string
	if fcNeoFieldTag := typeField.Tag.Get("fc_neo"); fcNeoFieldTag != "" {
		fieldTag = fcNeoFieldTag
	} else {
		fieldTag = strings.Split(typeField.Tag.Get("json"), ",")[0]
	}

	if fieldTag == "-" {
		return nil
	} else if fieldTag != "" {
		return &fieldTag
	} else {
		return &(typeField.Name)
	}
}
