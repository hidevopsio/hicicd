package docker

import (
	"github.com/fsouza/go-dockerclient"
	"fmt"
)

type Image struct {
	Endpoint string `json:"endpoint"`
}

func (i *Image) AuthCheck() error {
	client, err := docker.NewClient(i.Endpoint)
	if err != nil {
		panic(err)
	}
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		panic(err)
	}
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	}
	return nil
}
