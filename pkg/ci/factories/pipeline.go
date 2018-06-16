package factories

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"errors"
	"fmt"
	"github.com/hidevopsio/hicicd/pkg/ci"
	"github.com/hidevopsio/hicicd/pkg/ci/impl"
)

type PipelineFactory struct{}

const (
	JavaPipelineType    = "java"
	JavaWarPipelineType = "java-war"
	NodeJsPipelineType  = "nodejs"
	GitbookPipelineType = "gitbook"
	JavaCqrsType        = "java-cqrs"
	NodeJsDist          = "nodejs-dist"
)

func (pf *PipelineFactory) New(pipelineType string) (ci.PipelineInterface, error) {
	log.Debug("pf.NewPipeline()")
	switch pipelineType {
	case JavaPipelineType:
		return new(impl.JavaPipeline), nil
	case JavaWarPipelineType:
		return new(impl.JavaWarPipeline), nil
	case NodeJsPipelineType:
		return new(impl.NodeJsPipeline), nil
	case JavaCqrsType:
		return new(impl.JavaCqrsPipeline), nil
	case NodeJsDist:
		return new(impl.NodeJsDistPipeline), nil
	default:
		return nil, errors.New(fmt.Sprintf("pipeline type %d not recognized\n", pipelineType))
	}
	return nil, nil
}
