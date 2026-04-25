package audioengine

import "testing"

func BenchmarkRegistryLoadCached(b *testing.B) {
	reg := NewRegistry()
	data := make([]byte, 50992)
	for i := range data {
		data[i] = byte(i)
	}
	reg.StorePCM(1, data)
	loader := func(string) ([]byte, error) { return nil, nil }

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reg.Load(1, loader)
	}
}

func BenchmarkRegistryLoadFirstTime(b *testing.B) {
	rawData := make([]byte, 50992)
	for i := range rawData {
		rawData[i] = byte(i)
	}
	loader := func(string) ([]byte, error) {
		return append([]byte(nil), rawData...), nil
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reg := NewRegistry()
		reg.RegisterFile(1, "sfx_test.raw")
		reg.Load(1, loader)
	}
}

func BenchmarkPlaybackStateActivateAndRelease(b *testing.B) {
	b.ReportAllocs()
	ps := NewPlaybackState()
	for i := 0; i < b.N; i++ {
		ps.Activate(1, false)
		ps.Release(1)
	}
}

func BenchmarkPlaybackStateActivate8th(b *testing.B) {
	ps := NewPlaybackState()
	for j := 0; j < 7; j++ {
		ps.Activate(j, false)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ps.Activate(7, false)
		ps.Release(7)
	}
}

func BenchmarkResolveSFX(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ResolveSFX("paddle_computer")
	}
}

func BenchmarkRegistryHas(b *testing.B) {
	reg := NewRegistry()
	reg.StorePCM(1, []byte{0x01})
	reg.RegisterFile(2, "sfx_test.raw")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reg.Has(1)
		reg.Has(2)
		reg.Has(99)
	}
}

func BenchmarkStorePCMCopy(b *testing.B) {
	reg := NewRegistry()
	data := make([]byte, 50992)
	for i := range data {
		data[i] = byte(i)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		reg.StorePCM(1, data)
	}
}
