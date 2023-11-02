package nanoda

import "sort"

// スタイルのID
type StyleId uint32

// 話者のID
type SpeakerId string

// スタイル情報
type Style struct {
	Name string  `json:"name"`
	Id   StyleId `json:"id"`
}

// 内部用のメタ情報
type meta struct {
	Name        string    `json:"name"`
	Styles      []Style   `json:"styles"`
	Version     string    `json:"version"`
	SpeakerUuid SpeakerId `json:"speaker_uuid"`
}

// キャラクター単位のメタ情報
type Meta struct {
	Name      string
	Styles    []Style
	SpeakerId SpeakerId
}

func (v *Voicevox) getMetas() (m []meta) {
	for _, vvm := range v.vvms {
		m = append(m, vvm.metas...)
	}
	return
}

// VVMごとのメタ情報をキャラクター単位にマージする
func (v *Voicevox) sortMetas(ms []meta) []Meta {
	m := map[SpeakerId][]meta{}
	for _, meta := range ms {
		m[meta.SpeakerUuid] = append(m[meta.SpeakerUuid], meta)
	}
	metas := []Meta{}
	for _, v := range m {
		styles := []Style{}
		for _, meta := range v {
			styles = append(styles, meta.Styles...)
		}
		// スタイル情報をID順にソートする
		sort.Slice(styles, func(i, j int) bool {
			return styles[i].Id < styles[j].Id
		})
		metas = append(metas, Meta{
			Name:      v[0].Name,
			Styles:    styles,
			SpeakerId: v[0].SpeakerUuid,
		})
	}
	// メタ情報をスタイルのID順にソートする
	sort.Slice(metas, func(i, j int) bool {
		return metas[i].Styles[0].Id < metas[j].Styles[0].Id
	})
	return metas
}
