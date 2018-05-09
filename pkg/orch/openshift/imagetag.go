package openshift

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
	image "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/openshift/api/image/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
)

const (
	ApiVersion = "apps.openshift.io/v1"
	Kind       = "ImageStreamTag"
)

type ImageStreamTag struct {
	Name      string
	FullName  string
	Namespace string
	Version   string
	Interface image.ImageStreamTagInterface
}

func NewImageStreamTags(name, version, namespace string) (*ImageStreamTag, error) {
	clientSet, err := image.NewForConfig(orch.Config)
	return &ImageStreamTag{
		Name:      name,
		Namespace: namespace,
		Version:   version,
		FullName:  name + ":" + version,
		Interface: clientSet.ImageStreamTags(namespace),
	}, err
}

func (ist *ImageStreamTag) Create(fromNamespace string) (*v1.ImageStreamTag, error) {
	log.Debug("ImageStreamTag Create")
	imageTag := &v1.ImageStreamTag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ist.FullName,
			Namespace: ist.Namespace,
		},
		Tag: &v1.TagReference{
			From: &corev1.ObjectReference{
				Name:      ist.FullName,
				Namespace: fromNamespace,
				Kind:      Kind,
			},
		},
	}
	img, err := ist.Interface.Create(imageTag)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (ist *ImageStreamTag) Get() (*v1.ImageStreamTag, error) {
	log.Debug("ImageStreamTag Get")
	option := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       Kind,
		},
		IncludeUninitialized: true,
	}
	img, err := ist.Interface.Get(ist.FullName, option)
	if err != nil {
		log.Println("imageStreamTag get", err)
		return nil, err
	}
	return img, nil
}

func (ist *ImageStreamTag) Delete() error {
	log.Debug("ImageStreamTag Delete")
	meta := &metav1.DeleteOptions{
	}
	err := ist.Interface.Delete(ist.FullName, meta)
	return err
}

func (ist *ImageStreamTag) Update() (*v1.ImageStreamTag, error) {
	img := &v1.ImageStreamTag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ist.Name,
			Namespace: ist.Namespace,
			Labels: map[string]string{
				"app":     ist.Name,
				"version": ist.Version,
			},
		},
	}
	return ist.Interface.Update(img)
}
