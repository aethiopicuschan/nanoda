package constant

// 処理結果を示す結果コード
type ResultCode int32

const (
	OK                               ResultCode = 0  // 成功
	NOT_LOADED_OPENJTALK_DICT_ERROR  ResultCode = 1  // open_jtalk辞書ファイルが読み込まれていない
	GET_SUPPORTED_DEVICES_ERROR      ResultCode = 3  // サポートされているデバイス情報取得に失敗した
	GPU_SUPPORT_ERROR                ResultCode = 4  // GPUモードがサポートされていない
	STYLE_NOT_FOUND_ERROR            ResultCode = 6  // スタイルIDに対するスタイルが見つからなかった
	MODEL_NOT_FOUND_ERROR            ResultCode = 7  // 音声モデルIDに対する音声モデルが見つからなかった
	INFERENCE_ERROR                  ResultCode = 8  // 推論に失敗した
	EXTRACT_FULL_CONTEXT_LABEL_ERROR ResultCode = 11 // コンテキストラベル出力に失敗した
	INVALID_UTF8_INPUT_ERROR         ResultCode = 12 // 無効なutf8文字列が入力された
	PARSE_KANA_ERROR                 ResultCode = 13 // AquesTalk風記法のテキストの解析に失敗した
	INVALID_AUDIO_QUERY_ERROR        ResultCode = 14 // 無効なAudioQuery
	INVALID_ACCENT_PHRASE_ERROR      ResultCode = 15 // 無効なAccentPhrase
	OPEN_ZIP_FILE_ERROR              ResultCode = 16 // ZIPファイルを開くことに失敗した
	READ_ZIP_ENTRY_ERROR             ResultCode = 17 // ZIP内のファイルが読めなかった
	MODEL_ALREADY_LOADED_ERROR       ResultCode = 18 // すでに読み込まれている音声モデルを読み込もうとした
	STYLE_ALREADY_LOADED_ERROR       ResultCode = 26 // すでに読み込まれているスタイルを読み込もうとした
	INVALID_MODEL_DATA_ERROR         ResultCode = 27 // 無効なモデルデータ
	LOAD_USER_DICT_ERROR             ResultCode = 20 // ユーザー辞書を読み込めなかった
	SAVE_USER_DICT_ERROR             ResultCode = 21 // ユーザー辞書を書き込めなかった
	USER_DICT_WORD_NOT_FOUND_ERROR   ResultCode = 22 // ユーザー辞書に単語が見つからなかった
	USE_USER_DICT_ERROR              ResultCode = 23 // OpenJTalkのユーザー辞書の設定に失敗した
	INVALID_USER_DICT_WORD_ERROR     ResultCode = 24 // ユーザー辞書の単語のバリデーションに失敗した
	INVALID_UUID_ERROR               ResultCode = 25 // UUIDの変換に失敗した
)
