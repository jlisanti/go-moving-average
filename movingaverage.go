package movingaverage

import (
	"fmt"
	"time"
)

// specify length of moving average in hours
type MovingAverage struct {
	Length        float64
	Value         []float64
	Time          []time.Time
	TimeDiffs     []float64
	TimeValues    []string
	ValueSum      float64
	AverageValue  float64
	Averages      []float64
	Populated     bool
	NumValues     int
	SlopeNegative bool
	Slope         float64
	SlopeChange   int
	Intercept     float64
}

func NewMovingAverage(length float64) *MovingAverage {
	ma := MovingAverage{Length: length, NumValues: 0, ValueSum: 0.0, Populated: false}
	return &ma
}

func UpdateValue(ma *MovingAverage, newValue float64, newTime time.Time) {

	// Determine if the moving average already contains the correct length of time
	if len(ma.Time) > 1 {

		t1 := ma.Time[0]
		t2 := newTime
		timeDiff := t2.Sub(t1)

		// Need check to determine if time gap has exceed moving average
		//  if yes should reset average

		if timeDiff.Hours() > ma.Length {
			// Reset moving average
			fmt.Println("data seperation greater than average window")
			reset(ma)
			addValue(ma, newValue, newTime)
		}

		if !ma.Populated {
			addValue(ma, newValue, newTime)

			if timeDiff.Hours() >= ma.Length {
				ma.Populated = true
			}

		} else {

			delete := 0
			for i, _ := range ma.Value {
				timeDiff2 := newTime.Sub(ma.Time[i])
				if timeDiff2.Hours() > ma.Length {
					delete += 1
					ma.ValueSum -= ma.Value[i]
				} else {
					break
				}
			}

			ma.ValueSum += newValue

			ma.Value = append(ma.Value, newValue)
			ma.Time = append(ma.Time, newTime)
			ma.Averages = append(ma.Averages, newValue)
			ma.TimeValues = append(ma.TimeValues, newTime.String())
			ma.TimeDiffs = append(ma.TimeDiffs, timeDiff.Hours())

			if delete != 0 {
				ma.Value = append(ma.Value[:delete], ma.Value[delete+1:]...)
				ma.Time = append(ma.Time[:delete], ma.Time[delete+1:]...)
				ma.Averages = append(ma.Averages[:delete], ma.Averages[delete+1:]...)
				ma.TimeValues = append(ma.TimeValues[:delete], ma.TimeValues[delete+1:]...)
				ma.TimeDiffs = append(ma.TimeDiffs[:delete], ma.TimeDiffs[delete+1:]...)
			}

			ma.AverageValue = ma.ValueSum / float64(len(ma.Value))
		}

	} else {
		addValue(ma, newValue, newTime)
	}
}

func addValue(ma *MovingAverage, newValue float64, newTime time.Time) {
	ma.Value = append(ma.Value, newValue)
	ma.Time = append(ma.Time, newTime)
	ma.ValueSum += newValue
	ma.AverageValue = ma.ValueSum / float64(len(ma.Value))
	ma.Averages = append(ma.Averages, ma.ValueSum/float64(len(ma.Value)))
	ma.TimeValues = append(ma.TimeValues, newTime.String())
	ma.TimeDiffs = append(ma.TimeDiffs, 0.0)
}

func reset(ma *MovingAverage) {
	ma.Value = nil
	ma.Time = nil
	ma.ValueSum = 0.0
	ma.AverageValue = 0.0
	ma.Averages = nil
	ma.TimeValues = nil
	ma.TimeDiffs = nil
	ma.Populated = false
}
