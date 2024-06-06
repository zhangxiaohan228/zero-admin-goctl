package upgrade

import "github.com/zhangxiaohan228/zero-admin-goctl/internal/cobrax"

// Cmd describes an upgrade command.
var Cmd = cobrax.NewCommand("upgrade", cobrax.WithRunE(upgrade))
