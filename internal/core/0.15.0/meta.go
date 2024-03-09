package core_0_15_0

import "github.com/aethiopicuschan/nanoda/model"

type Meta struct {
	Name        string  `json:"name"`
	Styles      []Style `json:"styles"`
	SpeakerUuid string  `json:"speaker_uuid"`
	Version     string  `json:"version"`
}

type Style struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Type string `json:"type"`
}

func (m *Meta) ToCommon() (m2 model.Meta) {
	m2.Name = m.Name
	m2.Version = m.Version
	for _, s := range m.Styles {
		m2.Styles = append(m2.Styles, s.ToCommon())
	}
	return
}

func (s *Style) ToCommon() (s2 model.Style) {
	s2.Name = s.Name
	s2.Id = s.Id
	s2.Type = s.Type
	return
}
