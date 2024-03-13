package function

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/jmhobbs/velvetcache.org/hack/opengraph-images/internal/opengraph"
)

func init() {
	functions.HTTP("GenerateOpengraphImage", opengraph.Handler)
}
