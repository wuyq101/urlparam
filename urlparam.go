package urlparam

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
)

func Marshal(holder interface{}) url.Values {
	ret := make(url.Values)
	tp := reflect.TypeOf(holder)
	if tp.Kind() != reflect.Ptr {
		panic(errors.New("holder must be pointer"))
	}
	tp = reflect.TypeOf(holder).Elem()
	val := reflect.ValueOf(holder).Elem()
	for i := 0; i < val.NumField(); i++ {
		ft := tp.Field(i)
		tag := ft.Tag.Get("url")
		if len(tag) == 0 {
			continue
		}
		fv := val.Field(i)
		if fv.Kind() == reflect.String {
			ret.Set(tag, fv.String())
			continue
		}
		if fv.Kind() == reflect.Int || fv.Kind() == reflect.Int64 {
			ret.Set(tag, strconv.FormatInt(fv.Int(), 10))
			continue
		}
		panic(errors.New("url marshal only support int int64 string, not support for: " + fv.Kind().String()))
	}
	return ret
}

func Unmarshal(params url.Values, holder interface{}) error {
	tp := reflect.TypeOf(holder)
	if tp.Kind() != reflect.Ptr {
		panic(errors.New("holder must be pointer"))
	}
	tp = reflect.TypeOf(holder).Elem()
	val := reflect.ValueOf(holder).Elem()
	for i := 0; i < val.NumField(); i++ {
		ft := tp.Field(i)
		tag := ft.Tag.Get("url")
		if len(tag) == 0 {
			continue
		}
		required := ft.Tag.Get("require")
		fv := val.Field(i)
		x := params.Get(tag)
		if len(x) == 0 && "true" == required {
			return errors.New("missing required param: " + tag)
		}
		if fv.Kind() == reflect.String {
			fv.SetString(x)
			continue
		}
		if fv.Kind() == reflect.Int || fv.Kind() == reflect.Int64 {
			if len(x) == 0 {
				continue
			}
			intVal, err := strconv.ParseInt(x, 10, 64)
			if err != nil {
				return errors.New("param " + tag + " should be integer: " + err.Error())
			}
			fv.SetInt(intVal)
			continue
		}
		panic(errors.New("url marshal only support int int64 string. not support for: " + fv.Kind().String()))
	}
	return nil
}
