package impl

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hicicd/pkg/ci"
)

type NodeJsDistPipeline struct {
	ci.Pipeline
}

func (p *NodeJsDistPipeline) Deploy() error {
	log.Debug("NodeJsDistPipeline.Deploy()")
	return nil
}
