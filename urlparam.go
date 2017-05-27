package urlparam

import (
	"errors"
	"net/url"
	"reflect"
	"strconv"
	"time"
	"regexp"
	"strings"
	"code.byted.org/caijing_backend/alpha/service"
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
		required := ft.Tag.Get("required")
		x := params.Get(tag)

		fv := val.Field(i)
		if len(x) == 0 && "true" == required {
			return errors.New("missing required param: " + tag)
		}
		reg,ok :=ft.Tag.Lookup("regex")
		if fv.Kind() == reflect.String && ok {
			code, err := CheckRegex(reg, x, tag)
			if err!=nil {
				return err
			}
			fv.SetString(code)
			continue
		}
		layout,ok :=ft.Tag.Lookup("layout")
		if fv.Kind() == reflect.Struct && ok {
			t, err := CheckLayout(layout, x, tag)
			if err != nil {
				return err
			}
			fv.Set(t)
			continue
		}
		strList,ok :=ft.Tag.Lookup("list")
		if fv.Kind() == reflect.String && ok{
			str, err := CheckList(strList, x, tag)
			if err != nil {
				return err
			}
			fv.SetString(str)
			continue
		}
		if fv.Kind() == reflect.Int || fv.Kind() == reflect.Int64 {
			intVal, err := CheckInt(x, tag)
			if err != nil {
				return err
			}
			fv.SetInt(intVal)
			continue
		}
		panic(errors.New("url marshal only support int int64 string. not support for: " + fv.Kind().String()))
	}
	return nil
}

func CheckRegex(reg, x, tag string) (val string, err error) {
	regex := regexp.MustCompile(reg)
	if !regex.MatchString(x) {
		return "", service.NewServiceError(service.StatusInvalidParam, tag+": " + x + " is invalid parameter.")
	}
	return x, nil
}

func CheckLayout(layout, x, tag string) (val reflect.Value, err error) {
	t, err := time.Parse(layout, x)
	if err != nil {
		return reflect.ValueOf(nil), service.NewServiceError(service.StatusInvalidParam, tag+": " + x + " is invalid parameter.")
	}
	if t.Format(layout) != x {
		return reflect.ValueOf(nil), service.NewServiceError(service.StatusInvalidParam, tag+": " + x + " is invalid parameter.")
	}
	return reflect.ValueOf(t), nil
}

func CheckInt(x, tag string) (val int64, err error) {
	if len(x) == 0 {
		return 0, nil
	}
	intVal, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return 0, service.NewServiceError(service.StatusInvalidParam, tag+": " + x + " is invalid parameter.")
	}
	return intVal, nil
}

func CheckList(strList, x, tag string) (val string, err error) {
	list := strings.Split(strList, ",")
	if !CheckInList(x, list) {
		return "", service.NewServiceError(service.StatusInvalidParam, tag+": " + x + " is invalid parameter.")
	}

	return x, nil
}

func CheckInList(target string, list []string) bool {
	for _, item := range list {
		if target == item {
			return true
		}
	}
	return false
}
