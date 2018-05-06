package openshift

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/hidevopsio/hiboot/pkg/log"
	image "github.com/openshift/client-go/image/clientset/versioned/typed/image/v1"
	"github.com/openshift/api/image/v1"
	"github.com/hidevopsio/hicicd/pkg/orch"
	corev1 "k8s.io/api/core/v1"
)

const (
	ApiVersion = "apps.openshift.io/v1"
	Kind       = "ImageStreamTag"
)

type Tag struct {
	Name string
	Namespace string
	NewName string
	NewNamespce string
	Version string
	NewVersion string
}

type ImageStreamTag struct {
	Name         string
	Namespace    string
	Version      string
	Interface    image.ImageStreamTagInterface
}

func (ist *ImageStreamTag) Create(tag Tag) (*v1.ImageStreamTag, error) {
	imageTag := &v1.ImageStreamTag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tag.NewName + ":" + tag.NewVersion,
			Namespace: tag.NewNamespce,
		},
		Tag: &v1.TagReference{
			From: &corev1.ObjectReference{
				Name:      tag.Name + ":" + tag.Version,
				Namespace: tag.Namespace,
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

func (tag *ImageStreamTag) Get() (*v1.ImageStreamTag, error) {
	option := metav1.GetOptions{
		TypeMeta: metav1.TypeMeta{
			APIVersion: ApiVersion,
			Kind:       Kind,
		},
		IncludeUninitialized: true,
	}
	image1, err := tag.Interface.Get(tag.Name, option)
	if err != nil {
		log.Println("imageStreamTag get", err)
		return nil, err
	}
	return image1, nil
}

func (tag *ImageStreamTag) Delete() error {
	log.Debug("ImageStreamTag Delete")
	meta := &metav1.DeleteOptions{
	}
	err := tag.Interface.Delete(tag.Name, meta)
	return err
}

func (tag *ImageStreamTag) Update() (*v1.ImageStreamTag, error) {
	image := &v1.ImageStreamTag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tag.Name,
			Namespace: tag.Namespace,
			Labels: map[string]string{
				"app":     tag.Name,
				"version": tag.Version,
			},
		},
	}
	img, err := tag.Interface.Update(image)
	return img, err
}

func NewImageStreamTags(name, namespace string) (*ImageStreamTag, error) {
	clientSet, err := image.NewForConfig(orch.Config)
	tag := clientSet.ImageStreamTags(namespace)
	return &ImageStreamTag{
		Name:      name,
		Namespace: namespace,
		Interface: tag,
	}, err
}
