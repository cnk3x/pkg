package cmdx

import (
	"encoding/json"
	"time"

	"github.com/cnk3x/gopkg/jsonx"
)

type RestartOptions restartOptions

type restartOptions struct {
	Enable bool           `json:"-"`
	Delay  jsonx.Duration `json:"delay,omitempty"`
	Max    int            `json:"max,omitempty"`
}

func (p RestartOptions) MarshalJSON() ([]byte, error) {
	if !p.Enable {
		return json.Marshal(false)
	}
	if p.Delay == 0 && p.Max == 0 {
		return json.Marshal(p.Enable)
	}
	if p.Delay != 0 && p.Max != 0 {
		return json.Marshal(restartOptions(p))
	}
	if p.Delay == 0 {
		return json.Marshal(p.Max)
	}
	return json.Marshal(p.Delay)
}

func (p *RestartOptions) UnmarshalJSON(data []byte) (err error) {
	var v any
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch x := v.(type) {
	case bool:
		p.Enable = x
	case int:
		p.Enable = true
		p.Max = x
	case string:
		p.Enable = true
		d, e := time.ParseDuration(x)
		if err = e; e == nil {
			p.Delay = jsonx.Duration(d)
		}
	case map[any]any:
		p.Enable = true
		max, _ := x["max"].(float64)
		delay, _ := x["delay"].(string)
		p.Max = int(max)
		d, e := time.ParseDuration(delay)
		if err = e; e == nil {
			p.Delay = jsonx.Duration(d)
		}
	}
	return
}

func (p RestartOptions) ShouldRestart(t int) (time.Duration, bool) {
	if !p.Enable || (p.Max > 0 && t >= p.Max) {
		return 0, false
	}
	return max(p.Delay.Value(), time.Second), true //最低间隔1秒
}
