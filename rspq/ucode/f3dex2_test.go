package ucode

import "testing"

func TestF3DEX2Embedded(t *testing.T) {
	if len(RSPBoot) == 0 {
		t.Fatal("rspboot.bin is empty")
	}
	if len(RSPBoot) > 4096 {
		t.Fatalf("rspboot.bin too large: %d bytes", len(RSPBoot))
	}
	if len(F3DEX2Text) == 0 {
		t.Fatal("F3DEX2.bin is empty")
	}
	if len(F3DEX2Data) == 0 {
		t.Fatal("F3DEX2_data.bin is empty")
	}
	t.Logf("RSPBoot: %d bytes, F3DEX2 text: %d bytes, F3DEX2 data: %d bytes",
		len(RSPBoot), len(F3DEX2Text), len(F3DEX2Data))
}
