package blp

type Blp struct {
	Serial SerialRepo
	Ptz    ptzRepo
}

func New(serial SerialRepo, ptz ptzRepo) *Blp {
	return &Blp{
		Serial: serial,
		Ptz:    ptz,
	}
}
