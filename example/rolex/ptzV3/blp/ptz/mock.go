package ptz

import (
	"github.com/tiandh987/CGODemo/example/rolex/pkg/log"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV3/dsd"
	"time"
)

type MockAbility struct {
}

var _ AbilityRepo = (*MockAbility)(nil)

func NewMockAbility() *MockAbility {
	return &MockAbility{}
}

func (m MockAbility) Version() (string, error) {
	return "v1.2.3", nil
}

func (m MockAbility) Model() (string, error) {
	return "m1.2.3", nil
}

func (m MockAbility) Restart() error {
	log.Info("ptz is ready to restart")

	return nil
}

func (m MockAbility) Stop() error {
	log.Info("ptz stop")

	return nil
}

func (m MockAbility) Up(speed Speed) error {
	log.Infof("ptz up (speed: %d)", speed)

	return nil
}

func (m MockAbility) Down(speed Speed) error {
	log.Infof("ptz down (speed: %d)", speed)

	return nil
}

func (m MockAbility) Left(speed Speed) error {
	log.Infof("ptz left (speed: %d)", speed)

	return nil
}

func (m MockAbility) Right(speed Speed) error {
	log.Infof("ptz right (speed: %d)", speed)

	return nil
}

func (m MockAbility) LeftUp(speed Speed) error {
	//TODO implement me
	panic("implement me")
}

func (m MockAbility) RightUp(speed Speed) error {
	log.Infof("ptz right up (speed: %d)", speed)

	return nil
}

func (m MockAbility) LeftDown(speed Speed) error {
	log.Infof("ptz left down (speed: %d)", speed)

	return nil
}

func (m MockAbility) RightDown(speed Speed) error {
	log.Infof("ptz right down (speed: %d)", speed)

	return nil
}

func (m MockAbility) ZoomAdd() error {
	log.Infof("ptz zoom add")

	return nil
}

func (m MockAbility) ZoomSub() error {
	log.Infof("ptz zoom sub")

	return nil
}

var i = 0

func (m MockAbility) Position() (*dsd.Position, error) {
	//log.Infof("ptz get position(%d)", i)

	position := dsd.NewPosition()

	if i == 10 {
		position.Pan = 50
	}

	if i == 20 {
		position.Pan = 150
		i = 0
		return &position, nil
	}

	i++

	return &position, nil
}

func (m MockAbility) Goto(position *dsd.Position) error {
	log.Infof("ptz goto position %+v", position)

	time.Sleep(time.Second * 3)
	return nil
}
