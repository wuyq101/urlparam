package urlparam

import (
	"net/url"
	"testing"
	"time"
)

type TestObj struct {
	Code  string `url:"code" regex:"^[a-zA-Z0-9 .]+$" required:"true"`
	Date  time.Time `url:"date" layout:"2006-01-02" required:"true"`
	Datetime time.Time `url:"datetime" layout:"200601021504" required:"true"`
	Number  int `url:"number" required:"true"`
	Parameter  string `url:"parameter" list:"param1,param2,param3" required:"true"`
}

func TestMarshal(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2017-05-26")
	datetime, _ := time.Parse("200601021504","201705261503")
	obj := &TestObj{
		Code: "123",
		Date: date,
		Datetime:  datetime,
		Number: 454,
		Parameter: "param1",
	}
	params := Marshal(obj)
	t.Logf("params %v", params)
	t.Logf("query %s", params.Encode())
}

func TestUnmarshal(t *testing.T) {
	params := make(url.Values)
	params.Set("code", "123")
	params.Set("date", "2017-05-25")
	params.Set("datetime", "201705261504")
	params.Set("number", "16983")
	params.Set("parameter", "param1")
	var obj TestObj
	err := Unmarshal(params, &obj)
	if err != nil {
		t.Fatal("unmarshal error: %v", err)
	}
	t.Logf("result %+v", obj)
	t.Logf("params %+v", params)
}
