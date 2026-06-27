package main

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// web framework route'lari (embedded "web" lib) to'g'ri ishlashini tekshiradi.
func TestWebFramework(t *testing.T) {
	app := `
import "web";
a = new web.App();
a.get("/", func(req) { return {"ok": true}; });
a.get("/users/:id", func(req) { return {"id": req["params"]["id"]}; });
a.run();
`
	dir := t.TempDir()
	file := filepath.Join(dir, "index.bug")
	if err := os.WriteFile(file, []byte(app), 0644); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		path string
		want string
	}{
		{"/", `"ok":true`},
		{"/users/42", `"id":"42"`},
		{"/yoq", `Not Found`},
	}
	for _, c := range cases {
		req := httptest.NewRequest("GET", c.path, nil)
		res := processJob(Job{File: file, Request: req, BaseDir: dir + "/"})
		if !strings.Contains(res.Body, c.want) {
			t.Errorf("path %s: body=%q, %q kutilgan", c.path, res.Body, c.want)
		}
	}
}

// processJob web yordamchilari (status/cookie/header) to'g'ri ishlashini tekshiradi.
func TestProcessJobWebHelpers(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "index.bug")
	code := `
status(201);
header("Content-Type", "application/json");
cookie("session", "abc");
print("{\"ok\":true}");
`
	if err := os.WriteFile(file, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/?x=1", nil)
	res := processJob(Job{File: file, Request: req, BaseDir: dir + "/"})

	if res.Status != 201 {
		t.Errorf("status = %d, kutilgan 201", res.Status)
	}
	if !strings.Contains(res.Body, "ok") {
		t.Errorf("body = %q, 'ok' kutilgan", res.Body)
	}
	var hasCT, hasCookie bool
	for _, h := range res.Headers {
		if h.Key == "Content-Type" && h.Value == "application/json" {
			hasCT = true
		}
		if h.Key == "Set-Cookie" && strings.Contains(h.Value, "session=abc") {
			hasCookie = true
		}
	}
	if !hasCT {
		t.Error("Content-Type header yo'q")
	}
	if !hasCookie {
		t.Error("Set-Cookie header yo'q")
	}
}

// web framework + db (sqlite) integratsiyasi: handler req["db"] orqali so'rov qiladi.
func TestWebFrameworkWithDB(t *testing.T) {
	dir := t.TempDir()
	dbfile := filepath.Join(dir, "app.db")
	app := `
import "web";
import "db";
conn = new db.DB("` + dbfile + `", "sqlite");
conn.exec("create table users(id integer primary key, name text)", []);
conn.exec("insert into users(name) values(?)", ["Ali"]);
conn.exec("insert into users(name) values(?)", ["Vali"]);
a = new web.App();
a.setDB(conn);
a.get("/users", func(req) {
    return {"users": req["db"].query("select id, name from users order by id", [])};
});
a.run();
`
	file := filepath.Join(dir, "index.bug")
	if err := os.WriteFile(file, []byte(app), 0644); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/users", nil)
	res := processJob(Job{File: file, Request: req, BaseDir: dir + "/"})
	for _, want := range []string{"Ali", "Vali", `"id":1`} {
		if !strings.Contains(res.Body, want) {
			t.Errorf("body=%q, %q kutilgan", res.Body, want)
		}
	}
}

// panic recover qilinib 500 ga aylanishini tekshiradi (worker o'lmasligi kerak).
func TestProcessJobRecoversPanic(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "boom.bug")
	if err := os.WriteFile(file, []byte(`throw("boom");`), 0644); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/", nil)
	res := processJob(Job{File: file, Request: req, BaseDir: dir + "/"})
	if !strings.Contains(res.Body, "boom") {
		t.Errorf("body = %q, 'boom' kutilgan", res.Body)
	}
}
