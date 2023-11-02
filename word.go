package nanoda

import "github.com/aethiopicuschan/nanoda/internal/strings"

// 内部用の単語
type innerWord struct {
	surface       uintptr
	pronunciation uintptr
	accentType    uintptr
	wordType      int32
	priority      uint32
}

// 単語の種類
type WordType int32

const (
	VOICEVOX_USER_DICT_WORD_TYPE_PROPER_NOUN WordType = iota // 固有名詞
	VOICEVOX_USER_DICT_WORD_TYPE_COMMON_NOUN                 // 一般名詞
	VOICEVOX_USER_DICT_WORD_TYPE_VERB                        // 動詞
	VOICEVOX_USER_DICT_WORD_TYPE_ADJECTIVE                   // 形容詞
	VOICEVOX_USER_DICT_WORD_TYPE_SUFFIX                      // 接尾辞
)

// 単語
type Word struct {
	Surface       string   // 表記
	Pronunciation string   // 読み
	AccentType    uint64   // アクセント型(音が下がる場所を指す)
	WordType      WordType // 単語の種類
	Priority      uint32   // 優先度(0〜10までの整数)
}

func (w *Word) toinner() innerWord {
	return innerWord{
		surface:       strings.CString(w.Surface),
		pronunciation: strings.CString(w.Pronunciation),
		accentType:    uintptr(w.AccentType),
		wordType:      int32(w.WordType),
		priority:      w.Priority,
	}
}

// NewWordを呼ぶ際にアクセント型を指定する
func WithAccentType(at uint64) func(*Word) {
	return func(w *Word) {
		w.AccentType = at
	}
}

// NewWordを呼ぶ際に単語の種類を指定する
func WithWordType(wt WordType) func(*Word) {
	return func(w *Word) {
		w.WordType = wt
	}
}

// NewWordを呼ぶ際に単語の種類を指定する
func WithPriority(p uint32) func(*Word) {
	return func(w *Word) {
		w.Priority = p
	}
}

// 単語を作成する
// 辞書に登録するには別途AddWordを呼び出す必要がある
func NewWord(surface, pronunciation string, options ...func(*Word)) (w Word) {
	w = Word{
		Surface:       surface,
		Pronunciation: pronunciation,
		AccentType:    0,
		WordType:      0,
		Priority:      5,
	}
	for _, o := range options {
		o(&w)
	}

	return
}
