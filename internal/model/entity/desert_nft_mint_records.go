// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// DesertNftMintRecords is the golang structure for table desert_nft_mint_records.
type DesertNftMintRecords struct {
	Id               int64       `json:"id"               ` //
	Uid              int64       `json:"uid"              ` //
	Address          string      `json:"address"          ` //
	ShardAssetId     int64       `json:"shardAssetId"     ` //
	ShardAssetAmount int64       `json:"shardAssetAmount" ` //
	MascotAssetId    int64       `json:"mascotAssetId"    ` //
	MascotNftId      int64       `json:"mascotNftId"      ` //
	MascotAmount     int64       `json:"mascotAmount"     ` //
	Status           int         `json:"status"           ` //
	Signer           string      `json:"signer"           ` //
	Sig              string      `json:"sig"              ` //
	CreatedAt        *gtime.Time `json:"createdAt"        ` //
	UpdatedAt        *gtime.Time `json:"updatedAt"        ` //
}
