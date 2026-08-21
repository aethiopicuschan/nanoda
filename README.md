# nanoda

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/nanoda.svg)](https://pkg.go.dev/github.com/aethiopicuschan/nanoda)
[![Go Report Card](https://goreportcard.com/badge/github.com/aethiopicuschan/nanoda?branch=main)](https://goreportcard.com/report/github.com/aethiopicuschan/nanoda)
[![CI](https://github.com/aethiopicuschan/nanoda/actions/workflows/ci.yml/badge.svg)](https://github.com/aethiopicuschan/nanoda/actions/workflows/ci.yml)

nanodaは[VOICEVOX CORE](https://github.com/VOICEVOX/voicevox_core)の動的ライブラリをGolangから叩くためのライブラリです。`cgo`ではなく[ebitengine/purego](https://github.com/ebitengine/purego/)を利用しているため、簡単に使用することが可能です。

## VOICEVOXについて

サポートするVOICEVOX COREのバージョンは `0.16` 系(`voicevox_get_version()` が `0.16.` から始まる全パッチバージョン)としており、開発・動作確認は `0.16.4` を元にしています。

`0.15` 系のサポートは別系統([`main`ブランチ](https://github.com/aethiopicuschan/nanoda)、モジュールパス `github.com/aethiopicuschan/nanoda`、バージョンサフィックス無し)にあります。この`v2`はAPIの破壊的変更のあるバージョンとしてモジュールパスに`v2`を含みます。

nanoda自体は[MITライセンス](/LICENSE)ですが、利用に際してはVOICEVOXやOpenJTalkの利用規約に則る必要があることに注意してください。`0.16`以降、ビルド済みのVOICEVOX CORE・VOICEVOX ONNX Runtime本体もMITライセンスになりましたが、音声モデル(VVMファイル)にはキャラクターごとに個別の利用規約があるので、必ず確認してください。

## 使い方

```sh
go get github.com/aethiopicuschan/nanoda/v2@latest
```

もっとも簡単な例は以下のようになります。

```go
v, _ := nanoda.New("voicevox_core.dll", "dict/open_jtalk_dic_utf_8-1.11", "models")
s, _ := v.NewSynthesizer()
s.LoadAllModels()
wav, _ := s.Tts("ずんだもんなのだ！", 3)
defer wav.Close()
f, _ := os.Create("output.wav")
defer f.Close()
io.Copy(f, wav)
```

その他 `examples` ディレクトリにサンプルコードを置いていますので、ご活用ください。

### 必要なファイルについて

動作には以下のものが必要です。

- コアライブラリ(`voicevox_core.dll`/`.so`/`.dylib`)
- ONNX Runtime(`voicevox_onnxruntime.dll`/`.so`/`.dylib`) — **コアライブラリと同じディレクトリに置いてください。** コアライブラリ内部がファイル名だけでOSの共有ライブラリ検索パスを使ってロードするため、`New`に渡すコアライブラリのパスと別ディレクトリに置くと(Windowsの既定のDLL検索順序ではサブディレクトリまでは辿らないため)見つかりません。`0.15`系には無かった、`0.16`からの要件です。
- OpenJTalkの辞書ディレクトリ
- 音声モデル(VVMファイル)

[VOICEVOX CORE](https://github.com/VOICEVOX/voicevox_core)のREADMEもしくは同梱のダウンローダーに従って用意してください。ただし、上述のバージョンに対応するものを利用するようにしてください。


## 開発方針

以下の理由からなるべくnanoda側で処理を受け持ったり抽象化したりして機能を提供することを目指しています。

- 使いやすさの向上
- メモリまわりの安全性要件の確保
- VOICEVOXとアプリケーション間の密結合を避け、APIの変更等に強くする

## テスト

TODOです。ありません。

## 対応状況(0.16.x / このv2ブランチ)

`internal/core/core_0_16_0` が内部的に利用している関数のリストであり、必ずしも一致する形で公開されているわけではありません。テキストからWAVを生成する最小の経路(`New` → `NewSynthesizer` → `LoadVoiceModel`/`LoadAllModels` → `Tts`)にまず対応した段階で、AudioQuery・AccentPhrase・ユーザー辞書・メタ情報取得はまだ未移植です(`0.15`系の実装には存在します)。

- [x] voicevox_error_result_to_message
- [x] voicevox_get_onnxruntime_lib_versioned_filename
- [x] voicevox_get_version
- [x] voicevox_onnxruntime_load_once
- [x] voicevox_open_jtalk_rc_new
- [x] voicevox_open_jtalk_rc_delete
- [x] voicevox_synthesizer_new
- [x] voicevox_synthesizer_delete
- [x] voicevox_synthesizer_load_voice_model
- [x] voicevox_synthesizer_tts
- [x] voicevox_voice_model_file_open
- [x] voicevox_voice_model_file_delete
- [x] voicevox_wav_free
- [ ] voicevox_json_free (メタ情報などJSONを返す系統に着手したら必要)
- [ ] voicevox_onnxruntime_create_supported_devices_json
- [ ] voicevox_synthesizer_create_accent_phrases(_from_kana)
- [ ] voicevox_synthesizer_create_audio_query(_from_kana)
- [ ] voicevox_synthesizer_create_metas_json / voicevox_voice_model_file_create_metas_json
- [ ] voicevox_synthesizer_get_onnxruntime
- [ ] voicevox_synthesizer_is_gpu_mode
- [ ] voicevox_synthesizer_is_loaded_voice_model
- [ ] voicevox_synthesizer_replace_mora_data / replace_mora_pitch / replace_phoneme_length
- [ ] voicevox_synthesizer_synthesis (AudioQueryからの合成)
- [ ] voicevox_synthesizer_tts_from_kana
- [ ] voicevox_synthesizer_unload_voice_model
- [ ] voicevox_user_dict_*(全般)
- [ ] 歌唱音声合成系(voicevox_synthesizer_create_sing_frame_*, voicevox_synthesizer_frame_synthesis)
