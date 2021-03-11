package neo

//go:generate mockgen -destination ../../mock/neo4j.go -package mock github.com/neo4j/neo4j-go-driver/neo4j Record,Node,Driver,Session,Transaction,Result

import (
	"errors"
	"reflect"
	"testing"

	"github.com/freecloudio/server/domain/models"
	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"

	"github.com/golang/mock/gomock"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"github.com/stretchr/testify/assert"
)

func createTrCtxMock(mockCtrl *gomock.Controller) (trCtx *transactionCtx, sessionMock *mock.MockSession, txMock *mock.MockTransaction) {
	txMock = mock.NewMockTransaction(mockCtrl)
	sessionMock = mock.NewMockSession(mockCtrl)
	trCtx = &transactionCtx{sessionMock, txMock}
	return
}

func setupMockNewTransactionContext(mockCtrl *gomock.Controller, accessMode neo4j.AccessMode) (sessionMock *mock.MockSession, txMock *mock.MockTransaction) {
	txMock = mock.NewMockTransaction(mockCtrl)

	sessionMock = mock.NewMockSession(mockCtrl)
	sessionMock.EXPECT().BeginTransaction(gomock.Any()).Return(txMock, nil).Times(1)

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Session(accessMode).Return(sessionMock, nil).Times(1)
	neo = neoMock
	return
}

type testModel struct {
	Prop1   string `json:"prop1"`
	Prop2   string `json:"prop2" fc_neo:"changed2"`
	Prop3   string `fc_neo:"changed3,unique"`
	DontUse string `fc_neo:"-"`
	DefName string
	Token   models.Token
}

type uniqueModel struct {
	Unique    string `fc_neo:"uniqueProp,unique"`
	NotUnique string `fc_neo:"notUniqueProp"`
	DontUse   string `fc_neo:"-"`
}

type optionalModel struct {
	Optional    string `fc_neo:"uniqueProp,optional"`
	NotOptional string `fc_neo:"notOptionalProp"`
	DontUse     string `fc_neo:"-"`
}

type indexModel struct {
	Index    string `fc_neo:"indexProp,index"`
	NotIndex string `fc_neo:"indexProp"`
	DontUse  string `fc_neo:"-"`
}

func TestCloseNeo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Close().Return(nil).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	fcerr := closeNeo()
	assert.Nil(t, fcerr, "Closing neo driver failed")
}

func TestCloseNeoError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Close().Return(errors.New("Some error")).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	fcerr := closeNeo()
	assert.NotNil(t, fcerr, "Closing neo driver succeed but should fail")
	assert.Equal(t, fcerror.ErrDBCloseFailed, fcerr.ID, "Wrong error id for failed db close")
}

func TestTransactionClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	trCtx, sessionMock, _ := createTrCtxMock(mockCtrl)
	sessionMock.EXPECT().Close().Return(nil).Times(1)

	fcerr := trCtx.Close()
	assert.Nil(t, fcerr, "Closing transaction failed")
}

func TestTransactionCloseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	trCtx, sessionMock, _ := createTrCtxMock(mockCtrl)
	sessionMock.EXPECT().Close().Return(errors.New("Some error")).Times(1)

	fcerr := trCtx.Close()
	assert.NotNil(t, fcerr, "Closing transaction failed")
	assert.Equal(t, fcerror.ErrDBCloseSessionFailed, fcerr.ID, "Wrong error id for failed close transaction")
}

func TestTransactionFinishRollback(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	trCtx, sessionMock, txMock := createTrCtxMock(mockCtrl)
	txMock.EXPECT().Rollback().Return(nil).Times(1)
	sessionMock.EXPECT().Close().Return(nil).Times(1)

	inputErrID := fcerror.ErrUnknown
	fcerr := trCtx.Finish(fcerror.NewError(inputErrID, nil))
	assert.NotNil(t, fcerr, "Finishing transaction succeeded but should not")
	assert.Equal(t, inputErrID, fcerr.ID, "Finish err does not match input err id")
}

func TestTransactionFinishCommit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	trCtx, sessionMock, txMock := createTrCtxMock(mockCtrl)
	txMock.EXPECT().Commit().Return(nil).Times(1)
	sessionMock.EXPECT().Close().Return(nil).Times(1)

	fcerr := trCtx.Finish(nil)
	assert.Nil(t, fcerr, "Finishing transaction failed")
}

func TestTransactionCommit(t *testing.T) {
	tests := []struct {
		name        string
		commitErr   error
		closeErr    error
		expectedErr fcerror.ErrorID
	}{
		{name: "Success", commitErr: nil, closeErr: nil, expectedErr: fcerror.ErrorID(-1)},
		{name: "Commit Error", commitErr: errors.New("Some error"), closeErr: nil, expectedErr: fcerror.ErrDBCommitFailed},
		{name: "Close Error", commitErr: nil, closeErr: errors.New("Some error"), expectedErr: fcerror.ErrDBCommitFailed},
		{name: "Commit & Close Error", commitErr: errors.New("Some error"), closeErr: errors.New("Some error"), expectedErr: fcerror.ErrDBCommitFailed},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			trCtx, sessionMock, txMock := createTrCtxMock(mockCtrl)
			txMock.EXPECT().Commit().Return(test.commitErr).Times(1)
			sessionMock.EXPECT().Close().Return(test.closeErr).Times(1)

			fcerr := trCtx.Commit()
			if test.expectedErr == fcerror.ErrorID(-1) {
				assert.Nil(t, fcerr, "Commit transaction failed")
			} else {
				assert.NotNil(t, fcerr, "Commit transaction succeeded but should fail")
				assert.Equal(t, test.expectedErr, fcerr.ID, "Commit err does not match expected err id")
			}
		})
	}
}

func TestTransactionRollback(t *testing.T) {
	tests := []struct {
		name      string
		commitErr error
		closeErr  error
	}{
		{name: "Success", commitErr: nil, closeErr: nil},
		{name: "Commit Error", commitErr: errors.New("Some error"), closeErr: nil},
		{name: "Close Error", commitErr: nil, closeErr: errors.New("Some error")},
		{name: "Commit & Close Error", commitErr: errors.New("Some error"), closeErr: errors.New("Some error")},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			trCtx, sessionMock, txMock := createTrCtxMock(mockCtrl)
			txMock.EXPECT().Rollback().Return(test.commitErr).Times(1)
			sessionMock.EXPECT().Close().Return(test.closeErr).Times(1)

			trCtx.Rollback()
		})
	}
}

func TestNewTransactionContext(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inputAccessMode := neo4j.AccessModeWrite

	sessionMock, txMock := setupMockNewTransactionContext(mockCtrl, inputAccessMode)
	defer func() { neo = nil }()

	txCtx, fcerr := newTransactionContext(inputAccessMode)
	assert.Nil(t, fcerr, "Failed to create transaction context")
	assert.Equal(t, sessionMock, txCtx.session, "Transaction session does not match mock")
	assert.Equal(t, txMock, txCtx.neoTx, "Transaction does not match mock")
}

func TestNewTransactionContextSessionError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inputAccessMode := neo4j.AccessModeWrite

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Session(inputAccessMode).Return(nil, errors.New("Some error")).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	_, fcerr := newTransactionContext(inputAccessMode)
	assert.NotNil(t, fcerr, "Create transaction context succeeded but should fail")
	assert.Equal(t, fcerror.ErrDBTransactionCreationFailed, fcerr.ID, "Err ID does not match expected one")
}

func TestNewTransactionContextTxError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inputAccessMode := neo4j.AccessModeWrite

	sessionMock := mock.NewMockSession(mockCtrl)
	sessionMock.EXPECT().BeginTransaction(gomock.Any()).Return(nil, errors.New("Some error")).Times(1)
	sessionMock.EXPECT().Close().Return(nil).Times(1)

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Session(inputAccessMode).Return(sessionMock, nil).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	_, fcerr := newTransactionContext(inputAccessMode)
	assert.NotNil(t, fcerr, "Create transaction context succeeded but should fail")
	assert.Equal(t, fcerror.ErrDBTransactionCreationFailed, fcerr.ID, "Err ID does not match expected one")
}

func TestBuildConfigNameContainsNeededInfo(t *testing.T) {
	variant := "equal"
	label := "User"
	property := "email"

	configName := buildConfigName(variant, label, property)

	assert.Contains(t, configName, variant, "Expect configName to contain variant")
	assert.Contains(t, configName, label, "Expect configName to contain label")
	assert.Contains(t, configName, property, "Expect configName to contain property")
}

func TestInsertConfig(t *testing.T) {
	tests := []struct {
		name       string
		variant    NeoConfigVariant
		runErr     error
		consumeErr error
	}{
		{name: "unique", variant: NeoConfigUniqueConstraint, runErr: nil, consumeErr: nil},
		{name: "property", variant: NeoConfigPropertyConstraint, runErr: nil, consumeErr: nil},
		{name: "index", variant: NeoConfigIndex, runErr: nil, consumeErr: nil},
		{name: "unique run err", variant: NeoConfigUniqueConstraint, runErr: errors.New("Some error"), consumeErr: nil},
		{name: "property consume err", variant: NeoConfigPropertyConstraint, runErr: nil, consumeErr: errors.New("Some error")},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockResult := mock.NewMockResult(mockCtrl)
			if test.runErr == nil {
				mockResult.EXPECT().Consume().Return(nil, test.consumeErr).Times(1)
			}

			trCtx, _, txMock := createTrCtxMock(mockCtrl)
			txMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, test.runErr).Times(1)

			insertConfig(trCtx, test.variant, "label", "property")
		})
	}
}

func TestFetchNeoEdition(t *testing.T) {
	tests := []struct {
		name            string
		dbEdition       string
		expectedEdition NeoEdition
		dbErr           error
		recordSuccess   bool
	}{
		{name: "enterprise", dbEdition: "enterprise", expectedEdition: NeoEditionEnterprise, dbErr: nil, recordSuccess: true},
		{name: "community", dbEdition: "community", expectedEdition: NeoEditionCommunity, dbErr: nil, recordSuccess: true},
		{name: "unknown", dbEdition: "unknown", expectedEdition: NeoEditionCommunity, dbErr: nil, recordSuccess: true},
		{name: "enterprise db err", dbEdition: "unknown", expectedEdition: NeoEditionCommunity, dbErr: errors.New("Some error"), recordSuccess: true},
		{name: "enterprise record err", dbEdition: "unknown", expectedEdition: NeoEditionCommunity, dbErr: nil, recordSuccess: false},
	}

	for it := range tests {
		test := tests[it]
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockRecord := mock.NewMockRecord(mockCtrl)
			mockResult := mock.NewMockResult(mockCtrl)
			if test.dbErr == nil {
				mockRecord.EXPECT().Get("edition").Return(test.dbEdition, test.recordSuccess).Times(1)
				gomock.InOrder(
					mockResult.EXPECT().Next().Return(true).Times(1),
					mockResult.EXPECT().Record().Return(mockRecord).Times(1),
					mockResult.EXPECT().Err().Return(nil).Times(1),
					mockResult.EXPECT().Next().Return(false).Times(1),
				)
			}

			sessionMock, txMock := setupMockNewTransactionContext(mockCtrl, neo4j.AccessModeRead)
			sessionMock.EXPECT().Close().Return(nil).Times(1)
			txMock.EXPECT().Run(gomock.Any(), gomock.Any()).Return(mockResult, test.dbErr).Times(1)

			actualEdition, fcerr := fetchNeoEdition()
			if test.dbErr == nil && test.recordSuccess {
				assert.Equal(t, test.expectedEdition, actualEdition, "Got wrong neo edition")
			} else {
				assert.NotNil(t, fcerr, "Missing error for failed fetch neo edition")
				assert.Equal(t, fcerror.ErrDBReadFailed, fcerr.ID, "Got wrong error for failed neo edition fetch")
			}
		})
	}
}

func TestModelToMap(t *testing.T) {
	inputModel := &testModel{
		Prop1:   "value1",
		Prop2:   "value2",
		Prop3:   "value3",
		DontUse: "valueDontUse",
		DefName: "value4",
		Token:   models.Token("token"),
	}

	expectedMap := map[string]interface{}{
		"prop1":    inputModel.Prop1,
		"changed2": inputModel.Prop2,
		"changed3": inputModel.Prop3,
		"DefName":  inputModel.DefName,
		"Token":    inputModel.Token,
	}

	actualMap := modelToMap(inputModel)

	assert.Equal(t, expectedMap, actualMap, "Expected model-map does not match actual one")
}

func TestRecordToModel(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedModel := &testModel{
		Prop1:   "value1",
		Prop2:   "value2",
		Prop3:   "value3",
		DefName: "value4",
		Token:   models.Token("token"),
	}

	inputKey := "key"

	inputMap := map[string]interface{}{
		"prop1":    expectedModel.Prop1,
		"changed2": expectedModel.Prop2,
		"changed3": expectedModel.Prop3,
		"DefName":  expectedModel.DefName,
		"Token":    string(expectedModel.Token),
	}
	inputNode := mock.NewMockNode(mockCtrl)
	inputNode.EXPECT().Props().Return(inputMap).Times(1)

	inputRecord := mock.NewMockRecord(mockCtrl)
	inputRecord.EXPECT().Get(inputKey).Return(inputNode, true).Times(1)

	actualModel := &testModel{}
	fcerr := recordToModel(inputRecord, "key", actualModel)
	assert.Nil(t, fcerr, "Could not get model from record")
	assert.Equal(t, expectedModel, actualModel, "Model from record does not match expected model")
}

func TestRecordToModelWrongKey(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inputKey := "key"

	inputRecord := mock.NewMockRecord(mockCtrl)
	inputRecord.EXPECT().Get(inputKey).Return(nil, false).Times(1)

	actualModel := &testModel{}
	fcerr := recordToModel(inputRecord, inputKey, actualModel)
	assert.Error(t, fcerr, "Record to model did not fail with wrong key")
}

func TestRecordToModelWrongNodeType(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	inputKey := "key"

	inputRecord := mock.NewMockRecord(mockCtrl)
	inputRecord.EXPECT().Get(inputKey).Return("No Node", true).Times(1)

	actualModel := &testModel{}
	fcerr := recordToModel(inputRecord, inputKey, actualModel)
	assert.Error(t, fcerr, "Record to model did not fail with wrong node type")
}

func TestIsUniqueField(t *testing.T) {
	model := &uniqueModel{}

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for it := 0; it < modelValue.NumField(); it++ {
		typeField := modelType.Field(it)

		if it == 0 {
			assert.True(t, isUniqueField(typeField), "Field should be unique")
		} else if it > 0 {
			assert.False(t, isUniqueField(typeField), "Field should not be unique")
		}
	}
}

func TestIsOptionalField(t *testing.T) {
	model := &optionalModel{}

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for it := 0; it < modelValue.NumField(); it++ {
		typeField := modelType.Field(it)

		if it == 0 || it == 2 {
			assert.True(t, isOptionalField(typeField), "Field should be optional")
		} else if it == 1 {
			assert.False(t, isOptionalField(typeField), "Field should not be optional")
		}
	}
}

func TestIsIndexField(t *testing.T) {
	model := &indexModel{}

	modelValue := reflect.ValueOf(model).Elem()
	modelType := modelValue.Type()

	for it := 0; it < modelValue.NumField(); it++ {
		typeField := modelType.Field(it)

		if it == 0 {
			assert.True(t, isIndexField(typeField), "Field should be an index")
		} else if it > 0 {
			assert.False(t, isIndexField(typeField), "Field should not be an index")
		}
	}
}

func TestIsNotFoundError(t *testing.T) {
	assert.True(t, isNotFoundError(errors.New("pipapo result contains no records blub")), "error should be a not found error")
	assert.False(t, isNotFoundError(errors.New("random error")), "error should be not be an not found error")
}

func TestNeoToFcError(t *testing.T) {
	notFoundErr := fcerror.ErrUserNotFound
	otherErr := fcerror.ErrUnknown
	tests := []struct {
		neoErr error
		fcErr  fcerror.ErrorID
	}{
		{errors.New("result contains no records"), notFoundErr},
		{errors.New("some error"), otherErr},
		{nil, fcerror.ErrUnknown},
	}

	for _, test := range tests {
		if test.neoErr == nil {
			assert.Nil(t, neoToFcError(test.neoErr, notFoundErr, otherErr), "Got err for nil input")
		} else {
			assert.Equal(t, test.fcErr, neoToFcError(test.neoErr, notFoundErr, otherErr).ID, "Unexpected neo to fc error conversion")
		}
	}
}
