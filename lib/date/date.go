package date

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"time"
	"unicode"

	"github.com/araddon/dateparse"
)

// Date generic date struc
type Date struct {
	Base time.Time
}

// Now return current date
func Now() *Date {
	return &Date{
		Base: time.Now(),
	}
}

// Midnight return midnight of given date
func (d *Date) Midnight() *Date {
	d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), 0, 0, 0, 0, d.Base.Location())
	return d
}

// Calculate calculates relative date to given date
func (d *Date) Calculate(expr string) (*Date, error) {
	expr = strings.ToLower(expr)
	fields := strings.Fields(expr)
	if len(fields) == 0 {
		return d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	if strings.Contains(expr, "midnight") {
		d.Midnight()
	}

	if fields[0] == "tomorrow" {
		d.Base = d.Base.AddDate(0, 0, 1)
		return d, nil
	} else if fields[0] == "yesterday" {
		d.Base = d.Base.AddDate(0, 0, -1)
		return d, nil
	} else if fields[0] == "today" {
		return d, nil
	}
	if len(fields) < 2 {
		return d, fmt.Errorf("unable to parse date expression:%s", expr)
	}
	var i int
	if fields[0] == "next" {
		i = 1
	} else if fields[0] == "last" {
		i = 2
	} else {
		var err error
		i, err = strconv.Atoi(fields[0])
		if err != nil {
			return d, fmt.Errorf("unable to parse date expression:%s", expr)
		}
		if len(fields) > 2 {
			if fields[2] == "after" {
				if i < 0 {
					i = i * -1
				}
			}
			if fields[2] == "before" {
				if i > 0 {
					i = i * -1
				}
			}
		}
	}

	if strings.HasPrefix(fields[1], "year") {
		d.Base = d.Base.AddDate(i, 0, 0)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), 1, 1, 0, 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "month") {
		d.Base = d.Base.AddDate(0, i, 0)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), 0, 0, 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "day") {
		d.Base = d.Base.AddDate(0, 0, i)
		if strings.Contains(expr, "start") {
			d.Midnight()
		}
	} else if strings.HasPrefix(fields[1], "week") {
		d.Base = d.Base.AddDate(0, 0, i*7)

		if strings.Contains(expr, "start") {
			// Roll back to Monday:
			if wd := d.Base.Weekday(); wd == time.Sunday {
				d.Base = d.Base.AddDate(0, 0, -6)
			} else {
				d.Base = d.Base.AddDate(0, 0, -int(wd)+1)
			}
			d.Midnight()
		}

	} else if strings.HasPrefix(fields[1], "hour") {
		d.Base = d.Base.Add(time.Duration(i) * time.Hour)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), 0, 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "minute") {
		d.Base = d.Base.Add(time.Duration(i) * time.Minute)
		if strings.Contains(expr, "start") {
			d.Base = time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), d.Base.Minute(), 0, 0, d.Base.Location())
		}
	} else if strings.HasPrefix(fields[1], "second") {
		d.Base = d.Base.Add(time.Duration(i) * time.Second)
	}

	return d, nil

}

// DiffUnix add int64 to given date then return timestamp
func (d *Date) DiffUnix(t int64) time.Duration {
	return time.Duration(d.Base.Unix()-t) * time.Second
}

// DiffDate add date to given date return timestamp
func (d *Date) DiffDate(t Date) time.Duration {
	return time.Duration(d.Base.Unix()-t.Unix()) * time.Second
}

// DiffExpr add expr to date return timestamp
func (d *Date) DiffExpr(expr string) (time.Duration, error) {
	t := time.Date(d.Base.Year(), d.Base.Month(), d.Base.Day(), d.Base.Hour(), d.Base.Minute(), d.Base.Second(), d.Base.Nanosecond(), d.Base.Location())
	_, err := d.Calculate(expr)
	if err != nil {
		return time.Duration(0), err
	}
	return d.DiffTime(t), nil
}

// DiffTime add given time date return timestamp
func (d *Date) DiffTime(t time.Time) time.Duration {
	return time.Duration(d.Base.Unix()-t.Unix()) * time.Second
}

// Format formats given date
func (d *Date) Format(expr string) string {
	return d.Base.Format(expr)
}

// FormatS format given date as strftime syntax
func (d *Date) FormatS(f string) string {
	var (
		buf bytes.Buffer
		s   scanner.Scanner
	)
	var t = &d.Base

	if d == nil {
		now := time.Now()
		t = &now
	}

	s.Init(strings.NewReader(f))
	s.IsIdentRune = func(ch rune, i int) bool {
		return (ch == '%' && i <= 1) || (unicode.IsLetter(ch) && i == 1)
	}

	// Honor all white space characters.
	s.Whitespace = 0

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if len(txt) < 2 || !strings.HasPrefix(txt, "%") {
			buf.WriteString(txt)

			continue
		}

		buf.WriteString(formats.Apply(*t, txt[1:]))
	}

	return buf.String()
}

// Unix return timestamp of given date
func (d *Date) Unix() int64 {
	return d.Base.Unix()
}

// UnixNano return nano timestamp of given date
func (d *Date) UnixNano() int64 {
	return d.Base.UnixNano()
}

// FromString parse any string to Date
func FromString(expr string) (*Date, error) {
	t, err := dateparse.ParseLocal(expr)
	if err != nil {
		return nil, err
	}
	return &Date{
		Base: t,
	}, nil
}

// FromTime parse time to Date
func FromTime(t time.Time) *Date {
	return &Date{
		Base: t,
	}
}

// FomUnix parse timestamp to Date
func FromUnix(sec int64) *Date {
	t := time.Unix(sec, 0)
	return &Date{
		Base: t,
	}
}

func Parse(in interface{}) (*Date, error) {
	if v, ok := in.(int64); ok {
		return FromUnix(v), nil
	} else if v, ok := in.(time.Time); ok {
		return FromTime(v), nil
	} else if v, ok := in.(string); ok {
		return FromString(v)
	}
	return nil, fmt.Errorf("unrecognized date input")
}

func (f formatMap) Apply(t time.Time, txt string) string {
	fc, ok := f[txt]
	if !ok {
		return fmt.Sprintf("%%%s", txt)
	}

	return fc(t)
}

type formatMap map[string]func(time.Time) string

var formats = formatMap{
	"a": func(t time.Time) string { return t.Format("Mon") },
	"A": func(t time.Time) string { return t.Format("Monday") },
	"b": func(t time.Time) string { return t.Format("Jan") },
	"B": func(t time.Time) string { return t.Format("January") },
	"c": func(t time.Time) string { return t.Format(time.ANSIC) },
	"C": func(t time.Time) string { return t.Format("2006")[:2] },
	"d": func(t time.Time) string { return t.Format("02") },
	"D": func(t time.Time) string { return t.Format("01/02/06") },
	"e": func(t time.Time) string { return t.Format("_2") },
	"F": func(t time.Time) string { return t.Format("2006-01-02") },
	"g": func(t time.Time) string {
		y, _ := t.ISOWeek()
		return fmt.Sprintf("%d", y)[2:]
	},
	"G": func(t time.Time) string {
		y, _ := t.ISOWeek()
		return fmt.Sprintf("%d", y)
	},
	"h": func(t time.Time) string { return t.Format("Jan") },
	"H": func(t time.Time) string { return t.Format("15") },
	"I": func(t time.Time) string { return t.Format("03") },
	"j": func(t time.Time) string { return fmt.Sprintf("%03d", t.YearDay()) },
	"k": func(t time.Time) string { return fmt.Sprintf("%2d", t.Hour()) },
	"l": func(t time.Time) string { return fmt.Sprintf("%2s", t.Format("3")) },
	"m": func(t time.Time) string { return t.Format("01") },
	"M": func(t time.Time) string { return t.Format("04") },
	"n": func(t time.Time) string { return "\n" },
	"p": func(t time.Time) string { return t.Format("PM") },
	"P": func(t time.Time) string { return t.Format("pm") },
	"r": func(t time.Time) string { return t.Format("03:04:05 PM") },
	"R": func(t time.Time) string { return t.Format("15:04") },
	"s": func(t time.Time) string { return fmt.Sprintf("%d", t.Unix()) },
	"S": func(t time.Time) string { return t.Format("05") },
	"t": func(t time.Time) string { return "\t" },
	"T": func(t time.Time) string { return t.Format("15:04:05") },
	"u": func(t time.Time) string {
		d := t.Weekday()
		if d == 0 {
			d = 7
		}
		return fmt.Sprintf("%d", d)
	},
	// "U": func(t time.Time) string {
	// TODO
	// },
	"V": func(t time.Time) string {
		_, w := t.ISOWeek()
		return fmt.Sprintf("%02d", w)
	},
	"w": func(t time.Time) string {
		return fmt.Sprintf("%d", t.Weekday())
	},
	// "W": func(t time.Time) string {
	// TODO
	// },
	"x": func(t time.Time) string { return t.Format("01/02/2006") },
	"X": func(t time.Time) string { return t.Format("15:04:05") },
	"y": func(t time.Time) string { return t.Format("06") },
	"Y": func(t time.Time) string { return t.Format("2006") },
	"z": func(t time.Time) string { return t.Format("-0700") },
	"Z": func(t time.Time) string { return t.Format("MST") },
	"%": func(t time.Time) string { return "%" },
}
