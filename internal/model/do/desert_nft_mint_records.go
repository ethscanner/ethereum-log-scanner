// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// DesertNftMintRecords is the golang structure of table desert_nft_mint_records for DAO operations like Where/Data.
type DesertNftMintRecords struct {
	g.Meta           `orm:"table:desert_nft_mint_records, do:true"`
	Id               interface{} //
	Uid              interface{} //
	Address          interface{} //
	ShardAssetId     interface{} //
	ShardAssetAmount interface{} //
	MascotAssetId    interface{} //
	MascotNftId      interface{} //
	MascotAmount     interface{} //
	Status           interface{} //
	Signer           interface{} //
	Sig              interface{} //
	CreatedAt        *gtime.Time //
	UpdatedAt        *gtime.Time //
}
