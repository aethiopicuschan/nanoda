# nanoda

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen?style=flat-square)](/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/aethiopicuschan/nanoda.svg)](https://pkg.go.dev/github.com/aethiopicuschan/nanoda)
[![CI](https://github.com/aethiopicuschan/nanoda/actions/workflows/ci.yml/badge.svg)](https://github.com/aethiopicuschan/nanoda/actions/workflows/ci.yml)

nanodaは[VOICEVOX CORE](https://github.com/VOICEVOX/voicevox_core)の動的ライブラリをGolangから叩くためのライブラリです。`cgo`ではなく[ebitengine/purego](https://github.com/ebitengine/purego/)を利用しているため、簡単に使用することが可能です。

## VOICEVOXについて

サポートするVOICEVOX COREのバージョンは `0.15` としており、開発は `0.15.0-preview.13` を元にしています。

nanoda自体は[MITライセンス](/LICENSE)ですが、利用に際してはVOICEVOXやOpenJTalkの利用規約に則る必要があることに注意してください。

## 使い方

```sh
go get github.com/aethiopicuschan/nanoda@latest
```

もっとも簡単な例は以下のようになります。

```go
v, _ := nanoda.NewVoicevox("voicevox_core/libvoicevox_core.dylib", "voicevox_core/open_jtalk_dic_utf_8-1.11", "voicevox_core/model")
s, _ := v.NewSynthesizer()
s.LoadAllModels()
wav, _ := s.Tts("ずんだもんなのだ！", 3)
defer wav.Close()
f, _ := os.Create("output.wav")
defer f.Close()
io.Copy(f, wav)
```

その他 `examples` ディレクトリにサンプルコードを置いていますので、ご活用ください。

## 動作環境

```
GOARCH='arm64'
GOOS='darwin'
```

でのみ確認しています。

## 開発方針

以下の理由からなるべくnanoda側で処理を受け持ったり抽象化したりして機能を提供することを目指しています。

- 使いやすさの向上
- メモリまわりの安全性要件の確保
- VOICEVOXとアプリケーション間の密結合を避け、APIの変更等に強くする

## テスト

TODOです。ありません。

## 対応状況

以下は内部的に利用している関数のリストであり、必ずしも一致する形で公開されているわけではありません。

- [x] voicevox_create_supported_devices_json
- [x] voicevox_error_result_to_message
- [x] voicevox_get_version
- [x] voicevox_json_free
- [ ] voicevox_make_default_initialize_options
- [ ] voicevox_make_default_synthesis_options
- [ ] voicevox_make_default_tts_options
- [x] voicevox_open_jtalk_rc_delete
- [x] voicevox_open_jtalk_rc_new
- [x] voicevox_open_jtalk_rc_use_user_dict
- [x] voicevox_synthesizer_create_accent_phrases
- [x] voicevox_synthesizer_create_accent_phrases_from_kana
- [x] voicevox_synthesizer_create_audio_query
- [x] voicevox_synthesizer_create_audio_query_from_kana
- [x] voicevox_synthesizer_create_metas_json
- [x] voicevox_synthesizer_delete
- [x] voicevox_synthesizer_is_gpu_mode
- [x] voicevox_synthesizer_is_loaded_voice_model
- [x] voicevox_synthesizer_load_voice_model
- [x] voicevox_synthesizer_new_with_initialize
- [x] voicevox_synthesizer_replace_mora_data
- [x] voicevox_synthesizer_replace_mora_pitch
- [x] voicevox_synthesizer_replace_phoneme_length
- [x] voicevox_synthesizer_synthesis
- [x] voicevox_synthesizer_tts
- [x] voicevox_synthesizer_tts_from_kana
- [x] voicevox_synthesizer_unload_voice_model
- [x] voicevox_user_dict_add_word
- [x] voicevox_user_dict_delete
- [x] voicevox_user_dict_import
- [x] voicevox_user_dict_load
- [x] voicevox_user_dict_new
- [x] voicevox_user_dict_remove_word
- [x] voicevox_user_dict_save
- [x] voicevox_user_dict_to_json
- [x] voicevox_user_dict_update_word
- [ ] voicevox_user_dict_word_make
- [x] voicevox_voice_model_delete
- [x] voicevox_voice_model_get_metas_json
- [x] voicevox_voice_model_id
- [x] voicevox_voice_model_new_from_path
- [x] voicevox_wav_free
