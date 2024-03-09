package constant

type AccelerationMode int32

const (
	ACCELERATION_MODE_AUTO AccelerationMode = iota // 実行環境に合った適切なハードウェアアクセラレーションモードを選択する
	ACCELERATION_MODE_CPU                          // ハードウェアアクセラレーションモードを"CPU"に設定する
	ACCELERATION_MODE_GPU                          // ハードウェアアクセラレーションモードを"GPU"に設定する
)
