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

func (t Timestamp) Year() uint64 {
	return uint64((t-ts_min)/ts_ticks_per_year) + 1
}

type TimeOfDay uint8

const (
	TimeOfDay_NA TimeOfDay = iota

	TimeOfDay_Night
	TimeOfDay_Dawn
	TimeOfDay_Morning
	TimeOfDay_Afternoon
	TimeOfDay_Dusk

	timeOfDay_max
)

var timeOfDayNames [timeOfDay_max]string = [...]string{
	"N/A",

	"night",
	"dawn",
	"morning",
	"afternoon",
	"dusk",
}

func (tod TimeOfDay) String() string {
	return timeOfDayNames[tod]
}

func (t Timestamp) TimeOfDay() TimeOfDay {
	if t == 0 {
		return TimeOfDay_NA
	}
	return TimeOfDay((t.Tick()-1)*5/uint64(ts_ticks_per_day) + 1)
}

type Season uint8

const (
	Season_NA Season = iota

	Season_TheThaw
	Season_EarlySpring
	Season_MidSpring
	Season_LateSpring

	Season_TheBurn
	Season_EarlySummer
	Season_MidSummer
	Season_LateSummer

	Season_TheFall
	Season_EarlyAutumn
	Season_MidAutumn
	Season_LateAutumn

	Season_TheFreeze
	Season_EarlyWinter
	Season_MidWinter
	Season_LateWinter
	Season_YearsEnd

	season_max
)

var seasonNames [season_max]string = [...]string{
	"N/A",

	"the thaw",
	"early spring",
	"midspring",
	"late spring",

	"the burn",
	"early summer",
	"midsummer",
	"late summer",

	"the fall",
	"early autumn",
	"midautumn",
	"late autumn",

	"the freeze",
	"early winter",
	"midwinter",
	"late winter",
	"year's end",
}

func (s Season) String() string {
	return seasonNames[s]
}

func (t Timestamp) Season() Season {
	if t == 0 {
		return Season_NA
	}
	d := t.Day()
	switch {
	case d < 2:
		return Season_TheThaw
	case d < 53+2:
		return Season_EarlySpring
	case d < 53*2+2:
		return Season_MidSpring
	case d < 53*3+2:
		return Season_LateSpring
	case d < 53*3+3:
		return Season_TheBurn
	case d < 53*4+3:
		return Season_EarlySummer
	case d < 53*5+3:
		return Season_MidSummer
	case d < 53*6+3:
		return Season_LateSummer
	case d < 53*6+4:
		return Season_TheFall
	case d < 53*7+4:
		return Season_EarlyAutumn
	case d < 53*8+4:
		return Season_MidAutumn
	case d < 53*9+4:
		return Season_LateAutumn
	case d < 53*9+5:
		return Season_TheFreeze
	case d < 53*10+5:
		return Season_EarlyWinter
	case d < 53*11+5:
		return Season_MidWinter
	case d < 53*12+5:
		return Season_LateWinter
	default:
		return Season_YearsEnd
	}
}
