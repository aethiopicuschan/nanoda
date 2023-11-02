package nanoda

/*
利用可能なデバイスの情報。
あくまでVOICEVOX COREライブラリが対応しているかどうかであることに注意すること。
*/
type SupportedDevices struct {
	Cpu  bool `json:"cpu"`
	Cuda bool `json:"cuda"`
	Dml  bool `json:"dml"`
}
