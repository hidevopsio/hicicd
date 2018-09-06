package gitlab

import (
	"net/http"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/xanzy/go-gitlab"
	"github.com/hidevopsio/hicicd/pkg/scm"
	"github.com/jinzhu/copier"
)

type Repository struct {
	scm.TreeNode
}

type TreeNode struct {
	scm.TreeNode
}


func (r *Repository) GetRepository(baseUrl, token, filePath, ref string, pid int) (string, error) {
	log.Debug("Repository.Repository()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	opt := &gitlab.GetFileOptions{
		Ref: &ref,
		FilePath: &filePath,
	}
	file, _, err := c.RepositoryFiles.GetFile(pid, opt)
	if err != nil {
		return "", err
	}
	return file.Content, nil
}

func (r *Repository) ListTree(baseUrl, token, ref string, pid int)  ([]scm.TreeNode, error){
	log.Debug("Repository.ListTree()")
	log.Debugf("url: %v", baseUrl)
	c := gitlab.NewClient(&http.Client{}, token)
	c.SetBaseURL(baseUrl + ApiVersion)
	opt := &gitlab.ListTreeOptions{
		RefName: &ref,
	}
	tree, _, err := c.Repositories.ListTree(pid, opt)
	if err != nil {
		return nil, err
	}
	log.Info(tree)
	var treeNodes []scm.TreeNode
	for _, tr := range tree{
		treeNode := scm.TreeNode{}
		copier.Copy(&treeNode, tr)
		treeNodes = append(treeNodes, treeNode)
	}
	return treeNodes, nil
}
