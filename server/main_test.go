package main

import (
	"encoding/json"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLearning(t *testing.T) {
	for i := 0; i < 50; i++ {
		seen := map[int]bool{}
		for _, x := range generate(5) {
			if x < 2 || x > 199 || seen[x] {
				t.Fatal(x)
			}
			seen[x] = true
			if _, ok := evaluate([]int{x / 100, x / 10 % 10, x % 10}, x); !ok {
				t.Fatal(x)
			}
		}
	}
	for _, p := range [][]int{{0, 1}, {0, 10, 0}, {0, -1, 2}, {0, 0, 0}} {
		if _, ok := evaluate(p, 12); ok {
			t.Fatal(p)
		}
	}
	if mastery([]Row{{"ok": 1.0, "hint": 1.0}, {"ok": 0.0}, {"ok": 1.0}}) != 50 {
		t.Fatal("mastery")
	}
}
func TestFamily(t *testing.T) {
	a := App{db: filepath.Join(t.TempDir(), "test.db"), pin: token(), unit: Unit{"bridge", 5, 3}, failures: map[string][]time.Time{}}
	s, _ := files.ReadFile("schema.sql")
	a.sql(string(s))
	cookie := ""
	call := func(status int, d Row) any {
		t.Helper()
		b, _ := json.Marshal(d)
		req := httptest.NewRequest("POST", "/api", strings.NewReader(string(b)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", cookie)
		w := httptest.NewRecorder()
		a.api(w, req)
		if w.Code != status {
			t.Fatalf("%v: %d %s", d["action"], w.Code, w.Body.String())
		}
		if c := w.Result().Cookies(); len(c) > 0 {
			cookie = c[0].Name + "=" + c[0].Value
		}
		var v any
		json.Unmarshal(w.Body.Bytes(), &v)
		return v
	}
	call(403, Row{"action": "report"})
	call(200, Row{"action": "parent", "pin": a.pin})
	parent := cookie
	call(400, Row{"action": "create", "grade": 0, "name": "未登録", "pin": "1234"})
	for _, name := range []string{"ソラ", "ミオ"} {
		call(200, Row{"action": "create", "name": name, "grade": 2, "pin": "8273"})
	}
	list := call(200, Row{"action": "profiles"}).([]any)
	id := list[0].(map[string]any)["id"]
	call(403, Row{"action": "select", "id": id, "pin": "9999"})
	call(200, Row{"action": "select", "id": id, "pin": "8273"})
	call(409, Row{"action": "home", "id": list[1].(map[string]any)["id"]})
	call(403, Row{"action": "report"})
	for i := 0; i < 5; i++ {
		r := call(200, Row{"action": "play", "id": id}).(map[string]any)
		again := call(200, Row{"action": "play", "id": id}).(map[string]any)
		if again["run"] != r["run"] || again["target"] != r["target"] || again["step"] != float64(i) {
			t.Fatal("resume")
		}
		target := int(r["target"].(float64))
		d := Row{"action": "hint", "id": id, "run": r["run"], "step": i}
		call(200, d)
		d["action"] = "answer"
		d["parts"] = []int{0, 0, 0}
		call(200, d)
		d["parts"] = []int{target / 100, target / 10 % 10, target % 10}
		call(200, d)
		call(409, d)
	}
	other := list[1].(map[string]any)["id"]
	call(200, Row{"action": "select", "id": other, "pin": "8273"})
	if call(200, Row{"action": "home", "id": other}).(map[string]any)["elements"] != "0" {
		t.Fatal("learner isolation")
	}
	cookie = parent
	results := call(200, Row{"action": "report"}).([]any)
	if len(results) != 2 || results[0].(map[string]any)["reward"] != 3.0 {
		t.Fatal(results)
	}
	reopened := App{db: a.db}
	reopened.sql(string(s))
	if reopened.sql("SELECT count(*) FROM answer") != "10" {
		t.Fatal("lost history")
	}
}

func TestBackupRestore(t *testing.T) {
	cli, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("sqlite3 CLI required")
	}
	dir := t.TempDir()
	a := App{db: filepath.Join(dir, "original.db")}
	s, _ := files.ReadFile("schema.sql")
	a.sql(string(s))
	a.sql("INSERT INTO learner VALUES('one','テスト',4,0,'test-hash')")
	cmd := exec.Command(cli, a.db, ".backup 'copy.db'")
	cmd.Dir = dir
	if b, e := cmd.CombinedOutput(); e != nil {
		t.Fatalf("%s: %v", b, e)
	}
	restored := App{db: filepath.Join(dir, "copy.db")}
	if restored.sql("PRAGMA integrity_check") != "ok" || restored.sql("SELECT name FROM learner") != "テスト" {
		t.Fatal("restore")
	}
}
