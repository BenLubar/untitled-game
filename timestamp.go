package main

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

func (t Timestamp) Season() string {
	if t == 0 {
		return "N/A"
	}
	d := t.Day()
	switch {
	case d < 2:
		return "the thaw"
	case d < 53+2:
		return "early spring"
	case d < 53*2+2:
		return "midspring"
	case d < 53*3+2:
		return "late spring"
	case d < 53*3+3:
		return "the burn"
	case d < 53*4+3:
		return "early summer"
	case d < 53*5+3:
		return "midsummer"
	case d < 53*6+3:
		return "late summer"
	case d < 53*6+4:
		return "the fall"
	case d < 53*7+4:
		return "early autumn"
	case d < 53*8+4:
		return "midautumn"
	case d < 53*9+4:
		return "late autumn"
	case d < 53*9+5:
		return "the freeze"
	case d < 53*10+5:
		return "early winter"
	case d < 53*11+5:
		return "midwinter"
	case d < 53*12+5:
		return "late winter"
	default:
		return "year's end"

	}
}

func (t Timestamp) Year() uint64 {
	return uint64((t-ts_min)/ts_ticks_per_year) + 1
}
