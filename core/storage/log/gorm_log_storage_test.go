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

func TestXxx(t *testing.T) {
	by1 := [32]byte{1}
	by2 := [32]byte{1}

	fmt.Println(by1 == by2)
}
