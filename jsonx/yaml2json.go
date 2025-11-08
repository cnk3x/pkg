package jsonx

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"go.yaml.in/yaml/v3"
)

func translateStream(in io.Reader, out io.Writer) error {
	decoder := yaml.NewDecoder(in)
	for {
		var data any
		err := decoder.Decode(&data)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = transformData(&data)
		if err != nil {
			return err
		}
		output, err := Marshal(data)
		if err != nil {
			return err
		}
		data = nil
		_, err = out.Write(output)
		if err != nil {
			return err
		}
		_, err = io.WriteString(out, "\n")
		if err != nil {
			return err
		}
	}
}

func transformData(pIn *any) (err error) {
	switch in := (*pIn).(type) {
	case float64:
		if math.IsInf(in, 1) {
			*pIn = "+Inf"
		} else if math.IsInf(in, -1) {
			*pIn = "-Inf"
		} else if math.IsNaN(in) {
			*pIn = "NaN"
		}
		return nil
	case map[any]any:
		m := make(map[string]any, len(in))
		for k, v := range in {
			err = transformData(&v)
			if err != nil {
				return err
			}
			var sk string
			switch v := k.(type) {
			case string:
				sk = v
			case int:
				sk = strconv.Itoa(v)
			case bool:
				sk = strconv.FormatBool(v)
			case nil:
				sk = "null"
			case float64:
				f := v
				if math.IsInf(f, 1) {
					sk = "+Inf"
				} else if math.IsInf(f, -1) {
					sk = "-Inf"
				} else if math.IsNaN(f) {
					sk = "NaN"
				} else {
					sk = strconv.FormatFloat(f, 'f', -1, 64)
				}
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", v)
			}
			m[sk] = v
		}
		*pIn = m
	case map[string]any:
		m := make(map[string]any, len(in))
		for k, v := range in {
			err = transformData(&v)
			if err != nil {
				return err
			}
			m[k] = v
		}
		*pIn = m
	case []any:
		for i := len(in) - 1; i >= 0; i-- {
			err = transformData(&in[i])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
