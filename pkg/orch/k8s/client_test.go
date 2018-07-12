package k8s

import (
	"testing"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type KubernetesAPI struct {
	Suffix string
	Client  kubernetes.Interface
}

// NewNamespaceWithPostfix creates a new namespace with a stable postfix
func (k KubernetesAPI) NewNamespaceWithSuffix(namespace string) error {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("%s-%s", namespace, k.Suffix),
		},
	}

	_, err := k.Client.CoreV1().Namespaces().Create(ns)

	if err != nil {
		return err
	}

	return nil
}

func TestNewNamespaceWithSuffix(t *testing.T) {
	cases := []struct {
		ns string
	}{
		{
			ns: "test",
		},
	}

	api := &KubernetesAPI{
		Suffix: "unit-test",
		Client:  NewClientSet(),
	}

	for _, c := range cases {
		// create the postfixed namespace
		err := api.NewNamespaceWithSuffix(c.ns)
		assert.Equal(t, nil, err)

		_, err = api.Client.CoreV1().Namespaces().Get("test-unit-test", metav1.GetOptions{})
		assert.Equal(t, nil, err)

	}
}
