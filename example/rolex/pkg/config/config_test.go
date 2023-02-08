package config

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

type AudioGeneral struct {
	// 声音报警使能开关
	Enable int `json:"Enable"`
}

func TestViperConcurrent(t *testing.T) {
	opt := &Option{
		Path: ".",
		Name: "custom",
		Type: "json",
	}

	c := New(opt)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		count := 100
		var wg1 sync.WaitGroup

		for i := 0; i < count; i++ {
			wg1.Add(1)
			go func(index int) {
				defer wg1.Done()

				intn := rand.Intn(5)
				time.Sleep(time.Millisecond * time.Duration(intn))

				a := AudioGeneral{Enable: index}
				c.Set("audiogeneral", a)
			}(i)
		}

		wg1.Wait()
	}()

	go func() {
		defer wg.Done()

		count := 1000
		var wg2 sync.WaitGroup

		for i := 0; i < count; i++ {
			wg2.Add(1)

			go func(index int) {
				defer wg2.Done()

				intn := rand.Intn(5)
				time.Sleep(time.Millisecond * time.Duration(intn))

				c.IsSet("audiogeneral")

				var a AudioGeneral
				if err := c.UnmarshalKey("audiogeneral", &a); err != nil {
					t.Fatal(err)
				}
			}(i)
		}

		wg2.Wait()

	}()

	wg.Wait()
}

type StreamEnables struct {
	Enable bool `json:"Enable" mapstructure:"Enable"`
}
type StreamInfoCfg struct {
	Main   StreamEnables `json:"Main" mapstructure:"Main"`
	Extra1 StreamEnables `json:"Extra1" mapstructure:"Extra1"`
	Extra2 StreamEnables `json:"Extra2" mapstructure:"Extra2"`
}

func TestReadNotExistKey(t *testing.T) {
	opt := &Option{
		Path: ".",
		Name: "custom",
		Type: "json",
	}

	c := New(opt)

	var cfg []StreamInfoCfg
	if err := c.UnmarshalKey("streamInfo", &cfg); err != nil {
		t.Error(err)
	}

	t.Logf("cfg: %+v", cfg)
}

func TestSetNotExistKey(t *testing.T) {
	opt := &Option{
		Path: ".",
		Name: "custom",
		Type: "json",
	}

	c := New(opt)

	cfg := []StreamInfoCfg{
		{
			Main:   StreamEnables{Enable: true},
			Extra1: StreamEnables{Enable: true},
			Extra2: StreamEnables{Enable: true},
		},
		{
			Main:   StreamEnables{Enable: false},
			Extra1: StreamEnables{Enable: false},
			Extra2: StreamEnables{Enable: false},
		},
	}

	if err := c.Set("streamInfo", cfg); err != nil {
		t.Error(err)
	}
}
