package solver

import (
	"errors"
	"fmt"
	"log"
	"os"
)

// TimeData is an element of a TimeSeries
type TimeData struct {
	time  float64
	state []float64
}

// NewTimeData returns a TimeData created from the given t and X
func NewTimeData(t float64, X []float64) TimeData {
	data := TimeData{t, X}
	return data
}

// GetTime returns the value of time
func (data TimeData) GetTime() float64 {
	return data.time
}

// GetState returns the value of X (the state vector)
func (data TimeData) GetState() []float64 {
	return data.state
}

func (data TimeData) String() string {
	return fmt.Sprintf("t: %.4f, v: %v", data.time, data.state)
}

// Clone returns a deep copy of this TimeData
func (data TimeData) Clone() TimeData {
	state := make([]float64, len(data.state))
	for i := 0; i < len(data.state); i++ {
		state[i] = data.state[i]
	}
	return NewTimeData(data.time, state)
}

// TimeSeries defines an array of states (to keep trace of history)
type TimeSeries []TimeData

// Append adds a new data in the time series
func (series *TimeSeries) Append(data TimeData) {
	*series = append(*series, data)
}

// Clear reset the series to zero (no element)
func (series *TimeSeries) Clear() {
	*series = make(TimeSeries, 0)
}

func (series TimeSeries) String() string {
	s := ""
	for i := 0; i < len(series); i++ {
		s += fmt.Sprintf("%s\n", series[i].String())
	}
	return s
}

// Clone returns a deep copy of this TimeSeries
func (series TimeSeries) Clone() TimeSeries {
	clone := make(TimeSeries, len(series))
	for i := 0; i < len(series); i++ {
		clone[i] = series[i].Clone()
	}
	return clone
}

// ToCSV saves the TimeSeries in a file whose path is filepath. The header of
// the CSV file is "t;x0;x1;x2; ..."
func (series TimeSeries) ToCSV(filepath string) error {
	// Generate the default list of names for the X composantes
	names := make([]string, len(series[0].state))
	for j := 0; j < len(series[0].state); j++ {
		names[j] = fmt.Sprintf("x%d", j)
	}

	return series.ToCSVwithNames(filepath, names)
}

// ToCSVwithNames saves the TimeSeries in a file whose path is filepath, and
// with a header generated from the given list of names. The names are the names
// of the X components, then the header is "t;names[0];names[1]; ..."
func (series TimeSeries) ToCSVwithNames(filepath string, names []string) error {
	if len(series) == 0 {
		return errors.New("the timeseries has no data")
	}
	if len(names) != len(series[0].state) {
		return fmt.Errorf("the names (%v) does not match with the state dimension", names)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Printf("Creating the data file %s containing the time series %v", filepath, names)
	line := "t"
	for j := 0; j < len(names); j++ {
		line += fmt.Sprintf(";%s", names[j])
	}
	line += "\n"
	file.WriteString(line)
	for i := 0; i < len(series); i++ {
		data := series[i]
		line = fmt.Sprintf("%.4f", data.time)
		for j := 0; j < len(data.state); j++ {
			line += fmt.Sprintf(";%.12f", data.state[j])
		}
		line += "\n"
		file.WriteString(line)
	}
	return nil
}
