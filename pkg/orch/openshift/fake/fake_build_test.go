package fake

import (
	"testing"
	image "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/magiconair/properties/assert"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/orch"
)

func TestImageTag(t *testing.T){
	namespace := "demo-stage"
	clientset,_ := image.NewForConfig(orch.Config)
	i := clientset.ImageStreamTags(namespace)
	it, err := i.Get("admin", meta_v1.GetOptions{})
	assert.Equal(t, nil, err)
	log.Info("", it)
}


/*func TestFakeImageTag(t *testing.T){
	namespace := "demo-stage"
	clientset,_ := image.NewForConfig(orch.Config)
	fakeImage := clientset.
	client := new(fake.FakeImageStreamTags)
	client.Get()
}*/