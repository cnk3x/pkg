package jsonx

import (
	"encoding/json"
	"time"
)

type Duration time.Duration

const (
	Nanosecond  Duration = 1
	Microsecond          = 1000 * Nanosecond
	Millisecond          = 1000 * Microsecond
	Second               = 1000 * Millisecond
	Minute               = 60 * Second
	Hour                 = 60 * Minute
)

func (d Duration) Value() time.Duration         { return time.Duration(d) }
func (d Duration) MarshalJSON() ([]byte, error) { return json.Marshal(time.Duration(d).String()) }
func (d *Duration) UnmarshalJSON(data []byte) (err error) {
	var r string
	if err = json.Unmarshal(data, &r); err == nil {
		_d, e := time.ParseDuration(r)
		if err = e; err == nil {
			*d = Duration(_d)
		}
	}
	return
}
