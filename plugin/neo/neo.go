package neo

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/freecloudio/server/application/config"
	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/utils"
	"github.com/sirupsen/logrus"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var neo neo4j.Driver

type NeoEdition int

const (
	NeoEditionEnterprise NeoEdition = iota
	NeoEditionCommunity
)

type NeoConfigVariant int

const (
	NeoConfigUniqueConstraint = iota
	NeoConfigPropertyConstraint
	NeoConfigIndex
)

func initializeNeo(cfg config.Config) (fcerr *fcerror.Error) {
	logger := utils.CreateLogger(cfg.GetLoggingConfig())
	driver, err := neo4j.NewDriver(cfg.GetDBConnectionString(), neo4j.BasicAuth(cfg.GetDBUsername(), cfg.GetDBPassword(), ""), setConfig)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrDBInitializationFailed, err)
		return
	}
	err = driver.VerifyConnectivity()
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrDBUnavailable, err)
		return
	}

	neo = driver

	fcerr = initializeConstraintsAndIndexes(logger)
	if fcerr != nil {
		logger.WithError(fcerr).Error("Failed to initialize neo constraints")
		return
	}

	return
}

func closeNeo() (fcerr *fcerror.Error) {
	err := neo.Close()
	if err != nil {
		fcerr = fcerror.NewErrorSkipFunc(fcerror.ErrDBCloseFailed, nil)
		return
	}
	return
}

func setConfig(config *neo4j.Config) {
	config.Encrypted = false
	config.Log = neo4j.ConsoleLogger(neo4j.WARNING)
}

type transactionCtx struct {
	session neo4j.Session
	neoTx   neo4j.Transaction
	logger  utils.Logger
}

func (trCtx *transactionCtx) Close() (fcerr *fcerror.Error) {
	err := trCtx.session.Close()
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrDBCloseSessionFailed, err)
		trCtx.logger.WithError(fcerr).Error("Failed to close neo4j session")
		return
	}
	return
}

func (trCtx *transactionCtx) Finish(fcerr *fcerror.Error) *fcerror.Error {
	if fcerr != nil {
		trCtx.Rollback()
		return fcerr
	}

	return trCtx.Commit()
}

func (trCtx *transactionCtx) Commit() *fcerror.Error {
	var err error
	txErr := trCtx.neoTx.Commit()
	if txErr != nil {
		trCtx.logger.WithError(txErr).Error("Failed to commit neo transaction - close session anyway")
		err = txErr
	}
	sessErr := trCtx.session.Close()
	if sessErr != nil {
		trCtx.logger.WithError(sessErr).Error("Failed to close neo4j session")
		if err == nil {
			err = sessErr
		}
	}
	if err != nil {
		return fcerror.NewErrorSkipFunc(fcerror.ErrDBCommitFailed, err)
	} else {
		return nil
	}
}

func (trCtx *transactionCtx) Rollback() {
	err := trCtx.neoTx.Rollback()
	if err != nil {
		trCtx.logger.WithError(err).Error("Failed to rollback neo transaction - close session anyway")
	}
	err = trCtx.session.Close()
	if err != nil {
		trCtx.logger.WithError(err).Error("Failed to close neo session")
	}
}

func newTransactionContext(accessMode neo4j.AccessMode, logger utils.Logger) (txCtx *transactionCtx, fcerr *fcerror.Error) {
	session, err := neo.Session(accessMode)
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrDBTransactionCreationFailed, err)
		return
	}
	neoTx, err := session.BeginTransaction(neo4j.WithTxTimeout(10 * time.Second))
	if err != nil {
		fcerr = fcerror.NewError(fcerror.ErrDBTransactionCreationFailed, err)
		session.Close()
		return
	}
	txCtx = &transactionCtx{session, neoTx, logger}
	return
}

// Specify depending on the model tags which constraints should be set for a label
type labelModelMapping struct {
	label string
	model interface{}
}

// List of labels mapped to models filled in 'init' functions of each repository
var labelModelMappings []*labelModelMapping

func initializeConstraintsAndIndexes(logger utils.Logger) (fcerr *fcerror.Error) {
	neoEdition, fcerr := fetchNeoEdition(logger)
	if fcerr != nil {
		return
	}

	txCtx, fcerr := newTransactionContext(neo4j.AccessModeWrite, logger)
	if fcerr != nil {
		return
	}

	for _, constraint := range labelModelMappings {
		modelValue := reflect.ValueOf(constraint.model).Elem()
		modelType := modelValue.Type()

		for it := 0; it < modelType.NumField(); it++ {
			typeField := modelType.Field(it)
			dbNamePtr := getDBFieldName(typeField)
			if dbNamePtr == nil {
				continue
			}

			if isUniqueField(typeField) {
				insertConfig(txCtx, NeoConfigUniqueConstraint, constraint.label, *dbNamePtr)
			}
			if isIndexField(typeField) {
				insertConfig(txCtx, NeoConfigIndex, constraint.label, *dbNamePtr)
			}
			if neoEdition == NeoEditionEnterprise && !isOptionalField(typeField) {
				insertConfig(txCtx, NeoConfigPropertyConstraint, constraint.label, *dbNamePtr)
			}
		}
	}

	fcerr = txCtx.Commit()
	return
}

func buildConfigName(variant, label, property string) string {
	return fmt.Sprintf("%s_%s_%s", variant, label, property)
}

func insertConfig(txCtx *transactionCtx, variant NeoConfigVariant, label, property string) {
	var query string
	var variantName string
	switch variant {
	case NeoConfigUniqueConstraint:
		query = "CREATE CONSTRAINT %s IF NOT EXISTS ON (n:%s) ASSERT n.%s IS UNIQUE"
		variantName = "unique"
	case NeoConfigPropertyConstraint:
		query = "CREATE CONSTRAINT %s IF NOT EXISTS ON (n:%s) ASSERT EXISTS (n.%s)"
		variantName = "property"
	case NeoConfigIndex:
		query = "CREATE INDEX %s IF NOT EXISTS FOR (n:%s) ON (n.%s)"
		variantName = "index"
	default:
		txCtx.logger.WithField("variant", variant).Error("Unknown neo config variant")
		return
	}

	name := buildConfigName(variantName, label, property)
	res, err := txCtx.neoTx.Run(fmt.Sprintf(query, name, label, property), nil)
	if err == nil {
		_, err = res.Consume()
	}
	if err != nil {
		txCtx.logger.WithError(err).WithFields(logrus.Fields{"variant": variantName, "label": label, "property": property}).Error("Failed to create constraint")
	}
}

func fetchNeoEdition(logger utils.Logger) (neoEdition NeoEdition, fcerr *fcerror.Error) {
	txCtx, fcerr := newTransactionContext(neo4j.AccessModeRead, logger)
	if fcerr != nil {
		return
	}
	defer txCtx.Close()

	record, err := neo4j.Single(txCtx.neoTx.Run("CALL dbms.components() yield edition", nil))
	if err != nil {
		fcerr = neoToFcError(err, fcerror.ErrDBReadFailed, fcerror.ErrDBReadFailed)
		return
	}

	editionInt, ok := record.Get("edition")
	if !ok {
		fcerr = fcerror.NewError(fcerror.ErrDBReadFailed, nil)
		return
	}

	if editionInt.(string) == "enterprise" {
		return NeoEditionEnterprise, nil
	}
	return NeoEditionCommunity, nil
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

func recordToModel(record neo4j.Record, key string, model interface{}) *fcerror.Error {
	valInt, ok := record.Get(key)
	if !ok {
		return fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("value not found with key '"+key+"'"))
	}
	valNode, ok := valInt.(neo4j.Node)
	if !ok {
		return fcerror.NewError(fcerror.ErrModelConversionFailed, errors.New("value with key '"+key+"' could not be converted to 'neo4j.Node'"))
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
		var propVal reflect.Value
		switch valField.Type() {
		case reflect.TypeOf((models.UserID)("")):
			propVal = reflect.ValueOf(models.UserID(propInt.(string)))
		case reflect.TypeOf((models.NodeID)("")):
			propVal = reflect.ValueOf(models.NodeID(propInt.(string)))
		case reflect.TypeOf((models.Token)("")):
			propVal = reflect.ValueOf(models.Token(propInt.(string)))
		case reflect.TypeOf((models.NodeMimeType)("")):
			propVal = reflect.ValueOf(models.NodeMimeType(propInt.(string)))
		case reflect.TypeOf((models.NodeType)(0)):
			propVal = reflect.ValueOf(models.NodeType(propInt.(string)))
		case reflect.TypeOf((models.ShareMode)(0)):
			propVal = reflect.ValueOf(models.ShareMode(propInt.(string)))
		default:
			propVal = reflect.ValueOf(propInt)
		}

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
		fieldTag = strings.Split(fcNeoFieldTag, ",")[0]
	}

	if fieldTag == "" {
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

func isUniqueField(typeField reflect.StructField) bool {
	return fieldHasNeoTag(typeField, "unique", false)
}

func isOptionalField(typeField reflect.StructField) bool {
	return fieldHasNeoTag(typeField, "optional", true)
}

func isIndexField(typeField reflect.StructField) bool {
	return fieldHasNeoTag(typeField, "index", false)
}

func fieldHasNeoTag(typeField reflect.StructField, tagName string, noDBTagDef bool) bool {
	if getDBFieldName(typeField) == nil {
		return noDBTagDef
	}

	tagParts := strings.SplitN(typeField.Tag.Get("fc_neo"), ",", 2)
	if len(tagParts) < 2 {
		return false
	}
	return strings.Contains(tagParts[1], tagName)
}

func isNotFoundError(err error) bool {
	return strings.Contains(err.Error(), "result contains no records")
}

func neoToFcError(err error, notfound fcerror.ErrorID, other fcerror.ErrorID) *fcerror.Error {
	if err == nil {
		return nil
	} else if isNotFoundError(err) {
		return fcerror.NewErrorSkipFunc(notfound, err)
	} else if neo4j.IsAuthenticationError(err) || neo4j.IsSecurityError(err) {
		return fcerror.NewErrorSkipFunc(fcerror.ErrDBAuthentication, err)
	} else if neo4j.IsServiceUnavailable(err) {
		return fcerror.NewErrorSkipFunc(fcerror.ErrDBAuthentication, err)
	} else {
		return fcerror.NewErrorSkipFunc(other, err)
	}
}
