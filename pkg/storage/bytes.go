package storage

type ByteSize uint64

const (
	B  ByteSize = 1
	KB          = B << 10
	MB          = KB << 10
	GB          = MB << 10
	TB          = GB << 10
)

func (b ByteSize) Bytes() uint64 {
	return uint64(b)
}

func (b ByteSize) KBytes() float64 {
	return b.converterPref(KB)
}

func (b ByteSize) MBytes() float64 {
	return b.converterPref(MB)
}

func (b ByteSize) GBytes() float64 {
	return b.converterPref(GB)
}

func (b ByteSize) TBytes() float64 {
	return b.converterPref(TB)
}

func (b ByteSize) converterPref(targetSize ByteSize) float64 {
	return float64(b/targetSize) + float64(b%targetSize)/float64(targetSize)
}
