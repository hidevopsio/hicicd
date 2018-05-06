package openshift

import (
	"testing"
	"github.com/magiconair/properties/assert"
	"fmt"
	"log"
)

func TestCreateTags(t *testing.T) {
	tag := Tag{
		Name:        "demo-consumer",
		Namespace:   "demo-dev",
		NewName:     "demo-consumer",
		NewNamespce: "demo-test",
		Version:     "v1",
		NewVersion:  "v2",
	}
	name := tag.Name + ":" + tag.Version
	imageStreamtags, err := NewImageStreamTags(name, tag.NewNamespce)
	img, err := imageStreamtags.Create(tag)
	assert.Equal(t, nil, err)
	log.Printf("imageTag", img)
}

func TestGetTags(t *testing.T) {
	name := "hello-world:v1"
	namespace := "demo-dev"
	tag, err := NewImageStreamTags(name, namespace)
	img, err := tag.Get()
	assert.Equal(t, nil, err)
	fmt.Print(img)
}

func TestDeleteTag(t *testing.T)  {
	name := "demo-consumer:v2"
	namespace := "demo-test"
	tag, err := NewImageStreamTags(name, namespace)
	assert.Equal(t, nil, err)
	err = tag.Delete()
	assert.Equal(t, nil, err)
}