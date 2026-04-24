package audiov1

import "testing"

func TestCommandRingPushPop(t *testing.T) {
	var ring CommandRing
	if !ring.Push(Command{Kind: CmdPlaySFX, ID: 1}) {
		t.Fatal("Push to empty ring returned false")
	}
	cmd, ok := ring.Pop()
	if !ok {
		t.Fatal("Pop from non-empty ring returned false")
	}
	if cmd.Kind != CmdPlaySFX || cmd.ID != 1 {
		t.Fatalf("Pop returned %+v, want CmdPlaySFX ID 1", cmd)
	}
}

func TestCommandRingFull(t *testing.T) {
	var ring CommandRing
	for i := 0; i < cmdRingSize-1; i++ {
		if !ring.Push(Command{Kind: CmdPlaySFX, ID: uint16(i)}) {
			t.Fatalf("Push %d failed", i)
		}
	}
	if ring.Push(Command{Kind: CmdPlaySFX, ID: 999}) {
		t.Fatal("Push to full ring returned true")
	}
}

func TestOverlap10SFXWith4MaxInstances(t *testing.T) {
	e := Engine{SFXGain: GainFull, MusicGain: GainFull}
	for i := 0; i < 10; i++ {
		e.Seq++
		idx := e.AllocateSFXVoice(AssetID(7), DefaultMaxInst)
		e.Voices[idx] = Voice{State: VoicePlaying, AssetID: 7, Class: ClassSFX, StartSeq: e.Seq}
	}

	activeCount := 0
	for i := FirstSFXVoice; i <= MaxSFXVoices; i++ {
		if e.Voices[i].State == VoicePlaying && e.Voices[i].AssetID == 7 {
			activeCount++
			if e.Voices[i].StartSeq < 7 {
				t.Fatalf("active voice has StartSeq %d, expected only the 4 most recent", e.Voices[i].StartSeq)
			}
		}
	}
	if activeCount != DefaultMaxInst {
		t.Fatalf("active instances = %d, want %d", activeCount, DefaultMaxInst)
	}
}

func TestDrainCommandsInitializesSFXSource(t *testing.T) {
	data := []byte{0x10, 0x11, 0, 0, 0, 0, 0, 0, 0}
	aux := make([]byte, CodebookInts*2)
	e := Engine{
		SFXGain:   GainFull,
		MusicGain: GainFull,
		Manifest:  []AssetEntry{{ID: 0, Class: ClassSFX, Flags: FlagResident, Rate: 16000, AudibleFrames: 16, EncodedFrames: 16, DataBytes: uint32(len(data)), AuxBytes: uint32(len(aux)), MaxInstances: DefaultMaxInst}},
		Data:      data,
		Aux:       aux,
	}
	e.Ring.Push(Command{Kind: CmdPlaySFX, ID: 0})
	e.DrainCommands()

	active := -1
	for i := FirstSFXVoice; i <= MaxSFXVoices; i++ {
		if e.Voices[i].State == VoicePlaying {
			active = i
			break
		}
	}
	if active < 0 {
		t.Fatal("no active SFX voice")
	}
	dst := make([]int16, 16)
	n, ended := e.Sources[active].Fill(dst)
	if n == 0 || !ended {
		t.Fatalf("source Fill returned n=%d ended=%v, want final samples", n, ended)
	}
}

func TestDrainCommandsPerClassManifestIndex(t *testing.T) {
	e := Engine{
		SFXGain:   GainFull,
		MusicGain: GainFull,
		Manifest: []AssetEntry{
			{ID: 0, Class: ClassSFX, MaxInstances: DefaultMaxInst},
			{ID: 0, Class: ClassMusic, Flags: FlagLoop, MaxInstances: 1},
		},
	}
	e.Ring.Push(Command{Kind: CmdPlayMusic, ID: 0})
	e.DrainCommands()

	if e.Voices[MusicVoiceIndex].ManifestIndex != 1 {
		t.Fatalf("music ManifestIndex = %d, want 1", e.Voices[MusicVoiceIndex].ManifestIndex)
	}
}

func TestDrainCommandsStopMusicRequestsRamp(t *testing.T) {
	e := Engine{SFXGain: GainFull, MusicGain: GainFull}
	entry := &AssetEntry{ID: 0, Class: ClassMusic, Flags: FlagLoop, AudibleFrames: 128, EncodedFrames: 128, LoopLen: 128}
	InitSource(&e.Sources[MusicVoiceIndex], entry, make([]byte, 8*BlockBytes), &Codebook{}, State{})
	e.Voices[MusicVoiceIndex] = Voice{State: VoicePlaying, AssetID: 0, Class: ClassMusic, ManifestIndex: 0}

	e.Ring.Push(Command{Kind: CmdStopMusic})
	e.DrainCommands()

	if e.Voices[MusicVoiceIndex].State != VoicePlaying {
		t.Fatal("music voice stopped immediately instead of ramping")
	}
	if !e.Sources[MusicVoiceIndex].stopping {
		t.Fatal("music source did not receive stop ramp")
	}
}
