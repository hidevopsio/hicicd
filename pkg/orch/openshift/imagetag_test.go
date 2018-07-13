package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

const (
	name = "hiweb"
	namespace = "hidevopsio"
	fromNamespace = "hidevopsio-dev"
	version = "v1"
	fullName = name + ":" + version
)
func TestCrudTags(t *testing.T) {
	ist, err := NewImageStreamTags(name, version, namespace)
	assert.Equal(t, nil, err)

	is, err := ist.Create(fromNamespace)
	assert.Equal(t, nil, err)
	assert.Equal(t, fullName, is.Name)

	is, err = ist.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, fullName, is.Name)

	_, err = ist.Update(fromNamespace)
	assert.Equal(t, nil, err)

	err = ist.Delete()
	assert.Equal(t, nil, err)


}
