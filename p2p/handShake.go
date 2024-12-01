package p2p

type HandShakerFunc func(any) error

func NOPHandShakeFunc(any) error { return nil }
