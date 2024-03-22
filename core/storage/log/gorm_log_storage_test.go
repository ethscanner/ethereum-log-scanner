package log

import (
	"fmt"
	"testing"

	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestQueryLogs(t *testing.T) {
	ctx := gctx.GetInitCtx()
	d := NewGormLogStorage()
	query := scanner.LogQuery{}
	logs, err := d.QueryLogs(ctx, query)
	if err != nil {
		t.Fatalf("error %e", err)
	}
	fmt.Println(logs)
}
