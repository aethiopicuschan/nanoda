package nanoda

import (
	"unsafe"

	"github.com/aethiopicuschan/nanoda/internal/strings"
	"github.com/google/uuid"
)

// ユーザ辞書
// voicevox.NewUserDictで作成する
type UserDict struct {
	v        *Voicevox
	userDict uintptr
}

// ユーザ辞書を作成する
func (v *Voicevox) NewUserDict() (ud *UserDict) {
	ud = &UserDict{
		v:        v,
		userDict: v.voicevoxUserDictNew(),
	}
	return
}

// 単語を追加する
func (ud *UserDict) AddWord(word Word) (id string, err error) {
	iw := word.toinner()
	idPtr := make([]byte, 16)
	code := ud.v.voicevoxUserDictAddWord(ud.userDict, uintptr(unsafe.Pointer(&iw)), *(*uintptr)(unsafe.Pointer(&idPtr)))
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
		return
	}
	_uuid, err := uuid.FromBytes(idPtr)
	if err != nil {
		return
	}
	id = _uuid.String()
	return
}

// 単語を更新する
func (ud *UserDict) UpdateWord(id string, word Word) (err error) {
	iw := word.toinner()
	_uuid, err := uuid.Parse(id)
	if err != nil {
		return
	}
	ptr, err := _uuid.MarshalBinary()
	if err != nil {
		return
	}
	code := ud.v.voicevoxUserDictUpdateWord(ud.userDict, *(*uintptr)(unsafe.Pointer(&ptr)), uintptr(unsafe.Pointer(&iw)))
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// 単語を削除する
func (ud *UserDict) RemoveWord(id string) (err error) {
	_uuid, err := uuid.Parse(id)
	if err != nil {
		return
	}
	ptr, err := _uuid.MarshalBinary()
	if err != nil {
		return
	}
	code := ud.v.voicevoxUserDictRemoveWord(ud.userDict, *(*uintptr)(unsafe.Pointer(&ptr)))
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// ユーザ辞書を閉じる
func (ud *UserDict) Close() {
	ud.v.voicevoxUserDictDelete(ud.userDict)
}

// 他のユーザ辞書をインポートする
func (ud *UserDict) Import(other *UserDict) (err error) {
	code := ud.v.voicevoxUserDictImport(ud.userDict, other.userDict)
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// 指定されたパスからユーザ辞書を読み込む
func (ud *UserDict) Load(path string) (err error) {
	code := ud.v.voicevoxUserDictLoad(ud.userDict, path)
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// ユーザ辞書を指定されたパスに保存する
// 形式はJSONで、文字列として得たい場合はToJsonを使うこと
func (ud *UserDict) Save(path string) (err error) {
	code := ud.v.voicevoxUserDictSave(ud.userDict, path)
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// ユーザ辞書を使用する
func (ud *UserDict) Use() (err error) {
	code := ud.v.voicevoxOpenJtalkRcUseUserDict(ud.v.openJtalkRc, *(*uintptr)(unsafe.Pointer(&ud.userDict)))
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
	}
	return
}

// ユーザ辞書をJSONとして出力する
func (ud *UserDict) ToJson() (j string, err error) {
	var ptr *byte
	code := ud.v.voicevoxUserDictToJson(ud.userDict, uintptr(unsafe.Pointer(&ptr)))
	if code != VOICEVOX_RESULT_OK {
		err = ud.v.newError(code)
		return
	}
	defer ud.v.voicevoxJsonFree(uintptr(unsafe.Pointer(ptr)))
	j = strings.GoString(ptr)
	return
}
