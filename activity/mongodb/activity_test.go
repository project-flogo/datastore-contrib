package mongodb

import (
	"fmt"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}
func TestInsert(t *testing.T) {
	settings := &Settings{URI: "localhost:27017"}
	iCtx := test.NewActivityInitContext(settings, nil)

	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("method", "GET")
	tc.SetInput("collection", "numbers")
	tc.SetInput("dbname", "numbers2")
	act.Eval(tc)
}

func TestUpdate(t *testing.T) {
	settings := &Settings{URI: "localhost:27017"}
	iCtx := test.NewActivityInitContext(settings, nil)

	act, err := New(iCtx)
	assert.Nil(t, err)
	data := map[string]interface{}{"foo": "bar23"}
	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("method", "UPDATE")
	tc.SetInput("collection", "numbers")
	tc.SetInput("dbname", "numbers")
	tc.SetInput("keyname", "foo")
	tc.SetInput("keyvalue", "bar")
	tc.SetInput("data", data)
	res, err := act.Eval(tc)
	assert.Nil(t, err)
	fmt.Println("Res", res)
}
