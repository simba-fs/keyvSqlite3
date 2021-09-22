package keyvsqlite3

import "testing"
import "github.com/simba-fs/keyv"

var db *keyv.Keyv

type project struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type user struct {
	Name    string  `json:"name"`
	Mail    string  `json:"mail"`
	Project project `json:"project"`
}

func TestNew(t *testing.T) {
	d, err := keyv.New("sqlite3://database.sqlite", "")
	if err != nil {
		t.Error(err)
		panic(err)
	}
	db = d
}

func TestSetString(t *testing.T) {
	err := db.Set("package", "keyv")
	if err != nil {
		t.Error(err)
	}
}

func TestSetStruct(t *testing.T) {
	user := user{
		Name: "simba-fs",
		Mail: "simba-fs@example.com",
		Project: project{
			Name: "keyv",
			Url:  "pkg.go.dev/github.com/simba-fs/keyv",
		},
	}
	err := db.Set("user", user)
	if err != nil {
		t.Error(err)
	}
}

func TestGetString(t *testing.T) {
	packages := ""
	err := db.Get("package", &packages)
	if err != nil {
		t.Error(err)
	}
}

func TestGetStruct(t *testing.T) {
	user := user{}
	err := db.Get("user", &user)
	if err != nil {
		t.Error(err)
	}
}

func TestHas(t *testing.T) {
	ok := db.Has("user")
	if ok == false {
		t.Error("cannot get `user`")
	}

	ok = db.Has("notFound")
	if ok == true {
		t.Error("get `notFound`")
	}
}

func TestRemove(t *testing.T) {
	err := db.Remove("user")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestKeys(t *testing.T) {
	keys, err := db.Keys()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(keys)

	if len(keys) != 1 {
		t.Error("numbers of key is error")
	}
}

func TestClear(t *testing.T) {
	err := db.Clear()
	if err != nil {
		t.Error(err)
		return
	}

	keys, err := db.Keys()
	if err != nil {
		t.Error(err)
		return
	}

	if len(keys) != 0 {
		t.Error("keys: ", keys)
		t.Error("Not clear")
	}

}
