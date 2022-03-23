package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDatabaseFunctional(t *testing.T) {
	assert := assert.New(t)
	want := "tmp.db"
	defer os.Remove("tmp.db")

	_, err := GetDatabase("tmp.db", false)
	assert.FileExists(want)
	assert.NoError(err)

	got, err := GetDatabase("", false)
	assert.Equal(got, &Database{})
	assert.Error(err)
}

func TestPutGetFunctional(t *testing.T) {
	assert := assert.New(t)
	db, _ := GetDatabase("./tmp.db", false)
	defer os.Remove("./tmp.db")

	empty, _ := Get(db, "myBucket", "key")
	assert.Equal(empty, "")

	err := Put(db, "myBucket", "key", "value")
	assert.NoError(err)
	got, _ := Get(db, "myBucket", "key")
	assert.Equal(got, "value")

	err = Put(db, "", "key", "value")
	assert.Error(err)

	err = Put(db, "myBucket", "", "value")
	assert.Error(err)
}
