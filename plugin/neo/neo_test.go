package neo

//go:generate mockgen -destination ../../mock/neo4j.go -package mock github.com/neo4j/neo4j-go-driver/neo4j Record,Node,Driver

import (
	"errors"
	"reflect"
	"testing"

	"github.com/freecloudio/server/domain/models/fcerror"
	"github.com/freecloudio/server/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testModel struct {
	Prop1   string `json:"prop1"`
	Prop2   string `json:"prop2" fc_neo:"changed2"`
	Prop3   string `fc_neo:"changed3,unique"`
	DontUse string `fc_neo:"-"`
}

type uniqueModel struct {
	Unique    string `fc_neo:"uniqueProp,unique"`
	NotUnique string `fc_neo:"notUniqueProp"`
}

type optionalModel struct {
	Optional    string `fc_neo:"uniqueProp,optional"`
	NotOptional string `fc_neo:"notOptionalProp"`
}

type indexModel struct {
	Index    string `fc_neo:"indexProp,index"`
	NotIndex string `fc_neo:"indexProp"`
}

func TestCloseNeo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Close().Return(nil).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	fcerr := CloseNeo()
	assert.Nil(t, fcerr, "Closing neo driver failed")
}

func TestCloseNeoError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	neoMock := mock.NewMockDriver(mockCtrl)
	neoMock.EXPECT().Close().Return(errors.New("Some error")).Times(1)
	neo = neoMock
	defer func() { neo = nil }()

	fcerr := CloseNeo()
	assert.NotNil(t, fcerr, "Closing neo driver succeed but should fail")
	assert.Equal(t, fcerror.ErrDBCloseFailed, fcerr.ID, "Wrong error id for failed db close")
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

func TestModelToMap(t *testing.T) {
	inputModel := &testModel{
		Prop1:   "value1",
		Prop2:   "value2",
		Prop3:   "value3",
		DontUse: "valueDontUse",
	}

	expectedMap := map[string]interface{}{
		"prop1":    inputModel.Prop1,
		"changed2": inputModel.Prop2,
		"changed3": inputModel.Prop3,
	}

	actualMap := modelToMap(inputModel)

	assert.Equal(t, expectedMap, actualMap, "Expected model-map does not match actual one")
}

func TestRecordToModel(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	expectedModel := &testModel{
		Prop1: "value1",
		Prop2: "value2",
		Prop3: "value3",
	}

	inputKey := "key"

	inputMap := map[string]interface{}{
		"prop1":    expectedModel.Prop1,
		"changed2": expectedModel.Prop2,
		"changed3": expectedModel.Prop3,
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
		} else if it == 1 {
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

		if it == 0 {
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
		} else if it == 1 {
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
	}

	for _, test := range tests {
		if test.neoErr == nil {
			assert.Nil(t, neoToFcError(test.neoErr, notFoundErr, otherErr), "Got err for nil input")
		} else {
			assert.Equal(t, test.fcErr, neoToFcError(test.neoErr, notFoundErr, otherErr).ID, "Unexpected neo to fc error conversion")
		}
	}
}
