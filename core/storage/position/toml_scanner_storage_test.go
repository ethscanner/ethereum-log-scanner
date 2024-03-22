package position

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestTomlStorageUpdate(t *testing.T) {
	ctx := gctx.GetInitCtx()
	s := NewTomlScannerStorage()
	s.InsertUint64(ctx, "hello", 1)
	config, _ := s.GetConfig(ctx)
	fmt.Println(config)
}
