package cmdx

import (
	"cmp"
	"context"
	"encoding/json"
	"time"
)

type restartConfig struct {
	Type  string        `json:"type,omitempty"` //none, alway, unless-stopped, on-failure
	Delay time.Duration `json:"delay,omitempty"`
	Max   int           `json:"max,omitempty"`
}

type RestartConfig restartConfig

func (p RestartConfig) MarshalJSON() ([]byte, error) {
	if p.Type == "" || p.Type == "none" || ((p.Delay == 0 || p.Delay == time.Second*5) && (p.Max == 0 || p.Max == 10)) {
		return json.Marshal(p.Type)
	}
	return json.Marshal(restartConfig(p))
}

func (p *RestartConfig) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 0 {
		return
	}
	if data[0] == '"' {
		return json.Unmarshal(data, &p.Type)
	}

	r := (*restartConfig)(p)
	if err = json.Unmarshal(data, &r); err != nil {
		return
	}
	*p = RestartConfig(*r)
	return
}

func (p RestartConfig) CheckWait(ctx context.Context, stop_ctx context.Context, count int, err error) (restart bool) {
	restart = func() bool {
		switch p.Type {
		case "always":
			return true
		case "unless-stopped":
			return stop_ctx.Err() == nil
		case "on-failure":
			return stop_ctx.Err() == nil && err != nil
		default:
			return false
		}
	}()

	if !restart {
		return
	}

	delay := max(cmp.Or(p.Delay, time.Second*5), time.Second/2)

	select {
	case <-ctx.Done():
		return false //退出了
	case <-stop_ctx.Done():
		return true
	case <-time.After(delay): //最低1s, 默认5s
		return true
	}
}
