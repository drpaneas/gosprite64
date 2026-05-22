//go:build n64

package main

import (
	"embed"

	"github.com/clktmr/n64/drivers/cartfs"
	"github.com/drpaneas/gosprite64"
)

//go:embed assets/*
var embeddedAssets embed.FS

var assetFS = cartfs.Embed(embeddedAssets)

func init() {
	gosprite64.RegisterAssetFS(assetFS)
}
