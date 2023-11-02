package nanoda

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"unsafe"
)

type voicevoxVoiceModel struct {
	ptr   uintptr
	id    string
	metas []meta
}

// vvmファイルを読み込み構築する
func (v *Voicevox) loadVVMs(modelPath string) (vvms []voicevoxVoiceModel, err error) {
	files, err := os.ReadDir(modelPath)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if filepath.Ext(file.Name()) != ".vvm" {
			continue
		}
		fp := filepath.Join(modelPath, file.Name())
		vvm := voicevoxVoiceModel{}
		if code := v.voicevoxVoiceModelNewFromPath(fp, uintptr(unsafe.Pointer(&vvm.ptr))); code != VOICEVOX_RESULT_OK {
			err = fmt.Errorf("%s: %s", fp, v.GetMessageFromResult(code))
			return
		}
		vvm.id = v.voicevoxVoiceModelId(vvm.ptr)
		raw := v.voicevoxVoiceModelGetMetasJson(vvm.ptr)
		if err = json.Unmarshal([]byte(raw), &vvm.metas); err != nil {
			return
		}
		vvms = append(vvms, vvm)
	}
	return
}

// vvmを破棄する
func (v *Voicevox) deleteVVMs() {
	for _, vvm := range v.vvms {
		v.voicevoxVoiceModelDelete(vvm.ptr)
	}
	v.vvms = nil
}

type findOptions struct {
	all         bool
	byStyleId   bool
	bySpeakerId bool
	styleId     StyleId
	speakerId   SpeakerId
}

func withFindAll() func(*findOptions) {
	return func(o *findOptions) {
		o.all = true
	}
}

func withFindByStyleId(styleId StyleId) func(*findOptions) {
	return func(o *findOptions) {
		o.byStyleId = true
		o.styleId = styleId
	}
}

func withFindBySpeakerId(speakerId SpeakerId) func(*findOptions) {
	return func(o *findOptions) {
		o.bySpeakerId = true
		o.speakerId = speakerId
	}
}

// 音声モデルを逆引きする
func (v *Voicevox) findVVMs(options ...func(*findOptions)) (vvms []voicevoxVoiceModel, err error) {
	opt := findOptions{}
	for _, o := range options {
		o(&opt)
	}
	for _, vvm := range v.vvms {
		if opt.all {
			vvms = append(vvms, vvm)
			continue
		}
		for _, meta := range vvm.metas {
			if opt.bySpeakerId && meta.SpeakerUuid == opt.speakerId {
				vvms = append(vvms, vvm)
				continue
			}
			for _, style := range meta.Styles {
				if opt.byStyleId && style.Id == opt.styleId {
					vvms = append(vvms, vvm)
					return
				}
			}
		}
	}
	if len(vvms) == 0 {
		err = fmt.Errorf("not found")
	}
	return
}

// 音声モデルを読み込む(内部用)
func (s *Synthesizer) load(vvms []voicevoxVoiceModel) (err error) {
	for _, vvm := range vvms {
		loaded := s.v.voicevoxSynthesizerIsLoadedVoiceModel(s.synthesizer, vvm.id)
		if !loaded {
			code := s.v.voicevoxSynthesizerLoadVoiceModel(s.synthesizer, vvm.ptr)
			if code != VOICEVOX_RESULT_OK {
				err = s.v.newError(code)
				return
			}
		}
	}
	return
}

// スタイルIDを元にして音声モデルを読み込む
func (s *Synthesizer) LoadModelsFromStyleId(styleId StyleId) (err error) {
	vvms, err := s.v.findVVMs(withFindByStyleId(styleId))
	if err != nil {
		return
	}
	err = s.load(vvms)
	return
}

// 話者IDを元にして音声モデルを読み込む
func (s *Synthesizer) LoadModelsFromSpeakerId(speakerId SpeakerId) (err error) {
	vvms, err := s.v.findVVMs(withFindBySpeakerId(speakerId))
	if err != nil {
		return
	}
	err = s.load(vvms)
	return
}

// すべての音声モデルを読み込む
func (s *Synthesizer) LoadAllModels() (err error) {
	vvms, err := s.v.findVVMs(withFindAll())
	if err != nil {
		return
	}
	err = s.load(vvms)
	return
}

// 音声モデルをアンロードする(内部用)
func (s *Synthesizer) unload(vvms []voicevoxVoiceModel) (err error) {
	for _, vvm := range vvms {
		loaded := s.v.voicevoxSynthesizerIsLoadedVoiceModel(s.synthesizer, vvm.id)
		if loaded {
			code := s.v.voicevoxSynthesizerUnloadVoiceModel(s.synthesizer, vvm.id)
			if code != VOICEVOX_RESULT_OK {
				err = s.v.newError(code)
				return
			}
		}
	}
	return
}

// スタイルIDを元にして音声モデルをアンロードする
func (s *Synthesizer) UnloadModelsFromStyleId(styleId StyleId) (err error) {
	vvms, err := s.v.findVVMs(withFindByStyleId(styleId))
	if err != nil {
		return
	}
	err = s.unload(vvms)
	return
}

// 話者IDを元にして音声モデルをアンロードする
func (s *Synthesizer) UnloadModelsFromSpeakerId(speakerId SpeakerId) (err error) {
	vvms, err := s.v.findVVMs(withFindBySpeakerId(speakerId))
	if err != nil {
		return
	}
	err = s.unload(vvms)
	return
}

// すべての音声モデルをアンロードする
func (s *Synthesizer) UnloadAllModels() (err error) {
	vvms, err := s.v.findVVMs(withFindAll())
	if err != nil {
		return
	}
	err = s.unload(vvms)
	return
}
