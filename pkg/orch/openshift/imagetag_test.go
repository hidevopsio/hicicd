package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

const (
	name = "red-shop"
	namespace = "demo-test"
	fromNamespace = "demo-dev"
	version = "v2"
	fullName = name + ":" + version
)


func TestCreateTags(t *testing.T) {
	ist, err := NewImageStreamTags(name, version, namespace)
	assert.Equal(t, nil, err)
	is, err := ist.Create(fromNamespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, fullName, is.Name)
}

func TestGetTags(t *testing.T) {
	ist, err := NewImageStreamTags(name, version, namespace)
	assert.Equal(t, nil, err)
	is, err := ist.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, fullName, is.Name)
}

func TestDeleteTag(t *testing.T)  {
	ist, err := NewImageStreamTags(name, version, namespace)
	assert.Equal(t, nil, err)
	err = ist.Delete()
	assert.Equal(t, nil, err)
}

func TestUpdateTag(t *testing.T)  {
	ist, err := NewImageStreamTags(name, version, namespace)
	assert.Equal(t, nil, err)
	_, err = ist.Update(fromNamespace)
	assert.Equal(t, nil, err)
}