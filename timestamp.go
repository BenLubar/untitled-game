package main

import (
	"fmt"
)

// Timestamp 0 is the "nil time".
type Timestamp uint64

const (
	ts_ticks_per_day  Timestamp = 65535
	ts_days_per_year  Timestamp = 641
	ts_ticks_per_year Timestamp = ts_ticks_per_day * ts_days_per_year
	ts_max_years      Timestamp = 439125228929
	ts_min            Timestamp = 1
	ts_max            Timestamp = ^Timestamp(0)

	// and now we validate these constants (at compile time).
	_ts_0 = ts_max - ts_min + 1
	_ts_1 = ts_ticks_per_day * ts_days_per_year * ts_max_years
	_ts_2 = _ts_0 ^ _ts_1
)

var (
	// here comes the magic!
	_ [int64(_ts_2)]struct{}
	_ [-int64(_ts_2)]struct{}
)

func (t Timestamp) Tick() uint64 {
	return uint64((t-ts_min)%ts_ticks_per_day) + 1
}

func (t Timestamp) Day() uint64 {
	return uint64((t-ts_min)/ts_ticks_per_day%ts_days_per_year) + 1
}

func (t Timestamp) Year() uint64 {
	return uint64((t-ts_min)/ts_ticks_per_year) + 1
}

func (t Timestamp) String() string {
	if t == 0 {
		return "N/A"
	}

	return fmt.Sprintf("%d-%d-%d", t.Year(), t.Day(), t.Tick())
}
