package audioquery

import (
	"encoding/json"
	"math"
	"testing"
)

const query = `{"accent_phrases":[{"moras":[{"text":"ア","vowel":"a"}],"accent":1}],"speedScale":1.0,"pitchScale":0.0,"intonationScale":1.0,"volumeScale":1.0,"prePhonemeLength":0.1,"postPhonemeLength":0.1,"outputSamplingRate":24000,"outputStereo":false,"kana":"ア'"}`

func ptr(v float64) *float64 { return &v }

func TestApplyZeroKeepsQuery(t *testing.T) {
	got, err := Params{}.Apply(query)
	if err != nil {
		t.Fatal(err)
	}
	if got != query {
		t.Errorf("got %s, want the query unchanged", got)
	}
}

func TestApplyOverridesOnlyGivenFields(t *testing.T) {
	p := Params{SpeedScale: ptr(1.5), PitchScale: ptr(-0.05), IntonationScale: ptr(0)}
	got, err := p.Apply(query)
	if err != nil {
		t.Fatal(err)
	}

	var before, after map[string]any
	if err := json.Unmarshal([]byte(query), &before); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal([]byte(got), &after); err != nil {
		t.Fatal(err)
	}

	want := map[string]float64{"speedScale": 1.5, "pitchScale": -0.05, "intonationScale": 0}
	for key, v := range want {
		if after[key] != v {
			t.Errorf("%s = %v, want %v", key, after[key], v)
		}
	}
	for key, v := range before {
		if _, overridden := want[key]; overridden {
			continue
		}
		b, _ := json.Marshal(v)
		a, _ := json.Marshal(after[key])
		if string(a) != string(b) {
			t.Errorf("%s = %s, want %s unchanged", key, a, b)
		}
	}
}

func TestApplyRejectsNonPositiveSpeedScale(t *testing.T) {
	for _, v := range []float64{0, -1, math.NaN()} {
		if _, err := (Params{SpeedScale: ptr(v)}).Apply(query); err == nil {
			t.Errorf("speedScale %v: want an error", v)
		}
	}
}

func TestApplyRejectsNaN(t *testing.T) {
	if _, err := (Params{PitchScale: ptr(math.NaN())}).Apply(query); err == nil {
		t.Error("want an error for NaN")
	}
}

func TestApplyRejectsInvalidJSON(t *testing.T) {
	if _, err := (Params{SpeedScale: ptr(1)}).Apply("not json"); err == nil {
		t.Error("want an error for invalid JSON")
	}
}
