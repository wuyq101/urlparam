package urlparam

import (
	"net/url"
	"testing"
)

type TestObj struct {
	Id    int    `url:"id"`
	AppId int64  `url:"app_id"`
	Name  string `url:"name"`
}

func TestMarshal(t *testing.T) {
	obj := &TestObj{
		Id:    1,
		AppId: 1001,
		Name:  "test app",
	}
	params := Marshal(obj)
	t.Logf("params %v", params)
	t.Logf("query %s", params.Encode())
}

func TestUnmarshal(t *testing.T) {
	params := make(url.Values)
	params.Set("id", "1")
	params.Set("app_id", "1001")
	params.Set("name", "test app")
	var obj TestObj
	err := Unmarshal(params, &obj)
	if err != nil {
		t.Fatal("unmarshal error: %v", err)
	}
	t.Logf("result %+v", obj)
}
