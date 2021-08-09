package util

import "time"

// TimeParse is a function for formating from type data string to type data time.Time
// you can pick one the layout, format time layout :
// 1. "2006-01-02T15:04:05.000Z"
// 2. "2006-01-02 15:04:05"
// 3. "2006-01-02"
// 4. "02-01-2006"
// 5. "02 Jan 2006"
func TimeParse(s string, tl int) time.Time {

	// format template
	const (
		timeLayout1 = "2006-01-02T15:04:05.000Z"
		timeLayout2 = "2006-01-02 15:04:05"
		timeLayout3 = "2006-01-02"
		timeLayout4 = "02-01-2006 15:04:05"
		timeLayout5 = "02 Jan 2006"
	)

	// declare variable
	var timeLayout string

	// mapping format template
	switch tl {
	case 1:
		timeLayout = timeLayout1
	case 2:
		timeLayout = timeLayout2
	case 3:
		timeLayout = timeLayout3
	case 4:
		timeLayout = timeLayout4
	}

	// parsing time layout
	pTime, err := time.Parse(timeLayout, s)
	if err != nil {
		return pTime
	}

	// send result
	return pTime

}

// TimeFormat is a function for formating from type data time.Time to type data string
// you can pick one the layout, format time layout :
// 1. "2006-01-02 15:04:05"
// 2. "2006-01-02"
// 3. "02 January 2006"
func TimeFormat(t time.Time, tf int) string {

	// format template
	const (
		timeFormat1 = "2006-01-02 15:04:05"
		timeFormat2 = "2006-01-02"
		timeFormat3 = "02 January 2006"
	)

	// declare variable
	var timeFormat string

	// mapping formate template
	switch tf {
	case 1:
		timeFormat = timeFormat1
	case 2:
		timeFormat = timeFormat2
	case 3:
		timeFormat = timeFormat3
	}

	// formating time layout and send result
	return t.Format(timeFormat)

}
