package datatable

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type RequestColumnSearch struct {
	Value string `json:"value"`
	Regex bool   `json:"regex"`
}

func (rc *RequestColumnSearch) FromKV(key string, vals []string) (err error) {
	if len(vals) == 0 {
		return nil
	}
	val := vals[0]

	switch key {
	case "[value]":
		rc.Value = val
	case "[regex]":
		if val == "true" {
			rc.Regex = true
		} else {
			rc.Regex = false
		}
	}
	return nil
}

type RequestColumn struct {
	Data       string               `json:"data"`
	Name       string               `json:"name"`
	Searchable bool                 `json:"searchable"`
	Orderable  bool                 `json:"orderable"`
	Search     *RequestColumnSearch `json:"search"`
}

func (rc *RequestColumn) FromKV(key string, vals []string) (err error) {
	if len(vals) == 0 {
		return nil
	}
	val := vals[0]

	switch key {
	case "[name]":
		rc.Name = val
	case "[data]":
		rc.Data = val
	case "[searchable]":
		if val == "true" {
			rc.Searchable = true
		} else {
			rc.Searchable = false
		}
	case "[orderable]":
		if val == "true" {
			rc.Orderable = true
		} else {
			rc.Orderable = false
		}
	default:
		if strings.HasPrefix(key, "[search]") {
			if rc.Search == nil {
				rc.Search = &RequestColumnSearch{}
			}
			if err := rc.Search.FromKV(strings.TrimPrefix(key, "[search]"), vals); err != nil {
				return errors.Wrapf(err, "invalid column parameter %s", key)
			}
		}
	}
	return nil
}

type RequestOrder struct {
	Column int64  `json:"column"`
	Dir    string `json:"dir"`
}

func (ro *RequestOrder) FromKV(key string, vals []string) (err error) {
	if len(vals) == 0 {
		return nil
	}
	val := vals[0]

	switch key {
	case "[column]":
		if ro.Column, err = strconv.ParseInt(val, 10, 64); err != nil {
			return errors.Wrapf(err, "invalid parameter %s: %s", key, val)
		}
	case "[dir]":
		ro.Dir = val
	}
	return nil
}

type Request struct {
	Columns map[int64]*RequestColumn `json:"columns"`
	Order   map[int64]*RequestOrder  `json:"order"`
	Start   int64                    `json:"start"`
	Length  int64                    `json:"length"`
	Draw    int64                    `json:"draw"`
	Search  *RequestColumnSearch     `json:"search"`
}

var rexpArr = regexp.MustCompile(`^([a-z]+)\[([0-9]+)\](.+)$`)

func (r *Request) FromRequest(req *http.Request) (err error) {
	r.Columns = map[int64]*RequestColumn{}
	r.Order = map[int64]*RequestOrder{}
	r.Start = 0
	r.Length = 0
	r.Search = &RequestColumnSearch{}
	for key, vals := range req.URL.Query() {
		if len(vals) < 1 {
			return errors.New(fmt.Sprintf("no value for key %s", key))
		}
		val := vals[0]
		parts := rexpArr.FindAllStringSubmatch(key, -1)
		if parts != nil && len(parts) > 0 {
			idx, err := strconv.ParseInt(parts[0][2], 10, 64)
			if err != nil {
				return errors.Wrapf(err, "invalid parameter %s", key)
			}
			switch parts[0][1] {
			case "columns":
				if _, ok := r.Columns[idx]; !ok {
					r.Columns[idx] = &RequestColumn{Search: &RequestColumnSearch{}}
				}
				if err := r.Columns[idx].FromKV(parts[0][3], vals); err != nil {
					return errors.Wrapf(err, "invalid column parameter %s", key)
				}
			case "order":
				if _, ok := r.Order[idx]; !ok {
					r.Order[idx] = &RequestOrder{}
				}
				if err := r.Order[idx].FromKV(parts[0][3], vals); err != nil {
					return errors.Wrapf(err, "invalid column parameter %s", key)
				}
			}
		} else {
			switch key {
			case "start":
				if r.Start, err = strconv.ParseInt(val, 10, 64); err != nil {
					return errors.Wrapf(err, "invalid parameter %s: %s", key, val)
				}
			case "length":
				if r.Length, err = strconv.ParseInt(val, 10, 64); err != nil {
					return errors.Wrapf(err, "invalid parameter %s: %s", key, val)
				}
			case "draw":
				if r.Draw, err = strconv.ParseInt(val, 10, 64); err != nil {
					return errors.Wrapf(err, "invalid parameter %s: %s", key, val)
				}
			default:
				if strings.HasPrefix(key, "search") {
					if err := r.Search.FromKV(strings.TrimPrefix(key, "search"), vals); err != nil {
						return errors.Wrapf(err, "invalid parameter %s: %s", key, val)
					}
				}
			}
		}
	}
	return nil
}
