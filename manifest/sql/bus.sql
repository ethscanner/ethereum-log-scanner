DROP TABLE IF EXISTS "public"."desert_nft_mint_records";
CREATE TABLE "public"."desert_nft_mint_records" (
  id               BIGSERIAL PRIMARY KEY,
  "uid" int8 NOT NULL,
  "address" varchar(42) COLLATE "pg_catalog"."default" NOT NULL,
  "shard_asset_id" int8 NOT NULL,
  "shard_asset_amount" int8 NOT NULL,
  "mascot_asset_id" int8 NOT NULL,
  "mascot_nft_id" int8 NOT NULL,
  "mascot_amount" int8 NOT NULL,
  "status" int2,
  "signer" varchar(42) COLLATE "pg_catalog"."default" NOT NULL,
  "sig" varchar(132) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

ALTER TABLE "public"."desert_nft_mint_records" 
  OWNER TO "postgres";