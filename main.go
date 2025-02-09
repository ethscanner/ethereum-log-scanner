package main

import (
	_ "github.com/ethscanner/ethereum-log-scanner/internal/logic"

	"github.com/ethscanner/ethereum-log-scanner/internal/cmd"
	// _ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	cmd.Start(gctx.GetInitCtx())
}
