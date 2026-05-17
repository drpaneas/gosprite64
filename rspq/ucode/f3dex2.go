package ucode

import _ "embed"

//go:embed bin/rspboot.bin
var RSPBoot []byte

//go:embed bin/F3DEX2.bin
var F3DEX2Text []byte

//go:embed bin/F3DEX2_data.bin
var F3DEX2Data []byte
