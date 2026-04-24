//go:build !noos

package main

const (
	DefaultSFXRate           = 16000
	DefaultMusicRate         = 22050
	DefaultOutputRate        = 48000
	DefaultROMBudget         = 524288
	DefaultSFXResidentCap    = 32768
	DefaultMaxSFXInstances   = 4
	DefaultMaxMusicInstances = 1
	DefaultDecodeCostUsec    = 5

	BuildDirName    = "build"
	AudioBlobName   = "audio_v1.bin"
	AudioAuxName    = "audio_v1_aux.bin"
	AudioReportName = "audio_report.json"
	AudioEmbedName  = "audio_embed.go"
)
