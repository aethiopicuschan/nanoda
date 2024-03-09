package model

// キャラクター単位のメタ情報
type Meta struct {
	Name    string
	Styles  []Style
	Version string
}

// スタイル単位のメタ情報
type Style struct {
	Name string
	Id   int
	Type string
}
