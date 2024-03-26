package log

import (
	"fmt"
	"testing"

	"github.com/ethscanner/ethereum-log-scanner/core/scanner"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func TestQueryLogs(t *testing.T) {
	ctx := gctx.GetInitCtx()
	d := NewGormLogStorage()
	var State int = 0
	query := scanner.LogQuery{
		State:   &State,
		OrderBy: "id",
		Desc:    true,
	}
	logs, err := d.QueryLogs(ctx, query)
	if err != nil {
		t.Fatalf("error %e", err)
	}
	fmt.Println(logs)
}
