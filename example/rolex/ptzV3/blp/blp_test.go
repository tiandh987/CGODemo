package blp

import "testing"

func TestRequest_Validate(t *testing.T) {
	req := Request{
		Trigger: PowerUpTrigger,
		Ability: 3,
		ID:      1,
		Speed:   1,
	}

	if err := req.Validate(); err != nil {
		t.Error(err)
	}

}
