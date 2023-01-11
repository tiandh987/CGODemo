package blp

type Blp struct {
	ar AudioRepo
}

func New(ar AudioRepo) *Blp {
	return &Blp{ar: ar}
}

func (b *Blp) Play() int {
	return b.ar.Play()
}

type AudioRepo interface {
	Play() int
}
