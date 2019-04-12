package couchbase

import (
	"log"
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
	settings := &Settings{Username: "Administrator", Password: "password", BucketName: "test", BucketPassword: "", Server: "http://localhost:8091"}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Nil(t, err)

	tc := test.NewActivityContext(act.Metadata())
	tc.SetInput("key", 1234567889)

	tc.SetInput("data", `{"name":"foo"}`)

	tc.SetInput("method", `Insert`)
	tc.SetInput("expiry", 0)

	_, insertError := act.Eval(tc)
	assert.Nil(t, insertError)
	log.Println("TestInsert successful")
}
