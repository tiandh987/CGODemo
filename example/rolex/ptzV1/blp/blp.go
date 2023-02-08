package blp

type Blp struct {
	//Serial SerialRepo
	Ptz ptzRepo
}

func New(ptz ptzRepo) *Blp {
	return &Blp{
		Ptz: ptz,
	}
}
