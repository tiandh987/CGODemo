package cron

import (
	"encoding/json"
	"github.com/tiandh987/CGODemo/example/rolex/ptzV2/dsd"
	"testing"
)

func TestNew(t *testing.T) {

	var movements []dsd.PtzAutoMovement
	if err := json.Unmarshal([]byte(data), &movements); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", movements)

	_, err := New(movements)
	if err != nil {
		t.Fatal(err)
	}
}

var data = `[
  {
    "ID": 1,
    "Enable": true,
    "AutoHoming": {
      "Time": 30,
      "Enable": true
    },
    "Function": 4,
    "LinearScanID": 3,
    "TourID": 1,
    "PresetID": 1,
    "RegionScanID": 1,
    "PatternID": 1,
    "RunningFunction": 1,
    "Schedule": {
      "WeekDay": [
        {
          "Section": [
            {
              "TimeStr": [
                "11:14:48",
                "12:38:25"
              ],
              "TimeSec": [
                63582,
                80674
              ]
            },
            {
              "TimeStr": [
                "15:11:26",
                "19:35:29"
              ],
              "TimeSec": [
                76108,
                39433
              ]
            }
          ]
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        }
      ]
    }
  },
{
    "ID": 2,
    "Enable": true,
    "AutoHoming": {
      "Time": 30,
      "Enable": true
    },
    "Function": 4,
    "LinearScanID": 3,
    "TourID": 1,
    "PresetID": 1,
    "RegionScanID": 1,
    "PatternID": 1,
    "RunningFunction": 1,
    "Schedule": {
      "WeekDay": [
        {
          "Section": [
            {
              "TimeStr": [
                "08:14:48",
                "09:38:25"
              ],
              "TimeSec": [
                63582,
                80674
              ]
            },
            {
              "TimeStr": [
                "20:11:26",
                "23:35:29"
              ],
              "TimeSec": [
                76108,
                39433
              ]
            }
          ]
        },
        {
          "Section": [
            {
              "TimeStr": [
                "08:14:48",
                "09:38:25"
              ],
              "TimeSec": [
                63582,
                80674
              ]
            },
			{
              "TimeStr": [
                "10:14:48",
                "15:38:25"
              ],
              "TimeSec": [
                63582,
                80674
              ]
            }
          ]
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        },
        {
          "Section": []
        }
      ]
    }
  }
]`
