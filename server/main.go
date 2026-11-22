package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

//go:embed schema.sql
var files embed.FS

type Unit struct {
	ID             string
	Rounds, Reward int
}
type App struct {
	mu       sync.Mutex
	db, pin  string
	unit     Unit
	failures map[string][]time.Time
}
type Input struct {
	Action, ID, Name, Pin, Run string
	Grade, Icon, Step          int
	Parts                      []int
}
type Row map[string]any

func (a *App) rows(s string) []Row {
	out := []Row{}
	for _, l := range strings.Split(a.sql(s), "\n") {
		if l != "" {
			var r Row
			if json.Unmarshal([]byte(l), &r) != nil {
				panic("row")
			}
			out = append(out, r)
		}
	}
	return out
}
func (a *App) api(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if recover() != nil {
			http.Error(w, "保存できませんでした", 500)
		}
	}()
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != "POST" || r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "JSON POST required", 415)
		return
	}
	var d Input
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096))
	dec.DisallowUnknownFields()
	if dec.Decode(&d) != nil || dec.Decode(&struct{}{}) != io.EOF {
		http.Error(w, "入力を確認してね", 400)
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	now := time.Now().UnixMilli()
	owner := ""
	if c, e := r.Cookie("ticket"); e == nil {
		owner = a.sql(fmt.Sprintf("SELECT owner FROM login WHERE token=%s AND expires>%d", q(c.Value), now))
	}
	bad := func(code int, msg string) { http.Error(w, msg, code) }

	var out any = Row{"ok": true}
	switch d.Action {
	case "profiles":
		out = a.rows("SELECT json_object('id',id,'name',name,'icon',icon,'grade',grade) FROM learner ORDER BY rowid")
	case "parent", "select":
		if !a.authenticate(w, d, now) {
			return
		}
	case "logout":
		if c, e := r.Cookie("ticket"); e == nil {
			a.sql("DELETE FROM login WHERE token=" + q(c.Value))
		}
		http.SetCookie(w, &http.Cookie{Name: "ticket", Path: "/", MaxAge: -1})
	case "create", "report":
		if owner != "parent" {
			bad(403, "保護者の番号が必要です")
			return
		}
		if d.Action == "create" {
			d.Name = strings.TrimSpace(d.Name)
			if d.Name == "" || utf8.RuneCountInString(d.Name) > 16 || strings.ContainsRune(d.Name, 0) || d.Grade < 1 || d.Grade > 9 || d.Icon < 0 || d.Icon > 4 || len(d.Pin) != 4 || strings.Trim(d.Pin, "0123456789") != "" {
				bad(400, "名前・学年・4けたの番号を確認してください")
				return
			}
			id := token()
			a.sql(fmt.Sprintf("INSERT INTO learner VALUES(%s,%s,%d,%d,%s)", q(id), q(d.Name), d.Grade, d.Icon, q(digest(a.pin, id, d.Pin))))
		} else {
			out = a.rows("SELECT json_object('name',l.name,'grade',grade,'run',r.id,'started',started,'ended',ended,'step',r.step,'reward',reward,'answers',count(a.run),'correct',coalesce(sum(ok),0),'seconds',coalesce(sum(ms),0)/1000,'unit',unit,'owner',l.id) FROM learner l LEFT JOIN run r ON l.id=r.owner LEFT JOIN answer a ON r.id=a.run GROUP BY l.id,r.id ORDER BY started DESC")
			for _, r := range out.([]Row) {
				r["mastery"] = a.score(r["owner"].(string))
			}
		}
	case "home", "play", "hint", "answer":
		result, code, message := a.session(d, owner, now)
		if code != 200 {
			bad(code, message)
			return
		}
		out = result
	default:
		bad(404, "見つかりません")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
func (a *App) score(id string) int {
	return mastery(a.rows("SELECT json_object('ok',ok,'hint',answer.hint) FROM answer JOIN run ON run.id=answer.run WHERE owner=" + q(id) + " AND unit=" + q(a.unit.ID) + " ORDER BY answer.rowid DESC LIMIT 8"))
}
func env(k, v string) string {
	if s := os.Getenv(k); s != "" {
		return s
	}
	return v
}
func main() {
	dir := env("NIJI_DATA", "data")
	pin := os.Getenv("NIJI_PARENT_PIN")
	if len(pin) < 6 || strings.Trim(pin, "0123456789") != "" {
		log.Fatal("NIJI_PARENT_PIN に6けた以上の番号を設定してください")
	}
	if e := os.MkdirAll(dir, 0700); e != nil {
		log.Fatal(e)
	}
	a := App{db: filepath.Join(dir, "nijitama.db"), pin: pin, failures: map[string][]time.Time{}}
	var catalog struct{ Unit Unit }
	b, e := os.ReadFile("curriculum/islands.json")
	if e != nil || json.Unmarshal(b, &catalog) != nil || catalog.Unit.Rounds < 1 || catalog.Unit.Rounds > 198 {
		log.Fatal("curriculum invalid")
	}
	a.unit = catalog.Unit
	s, _ := files.ReadFile("schema.sql")
	if version := a.sql("PRAGMA user_version"); version == "0" {
		a.sql(string(s))
	} else if version != "1" {
		log.Fatal("unsupported database version ", version)
	}
	http.HandleFunc("/api", a.api)
	http.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	})
	webRoot := "web/standalone"
	if _, err := os.Stat("web/vendor"); err == nil {
		webRoot = "web"
	}
	if _, err := os.Stat("web/assets"); err == nil {
		webRoot = "web"
	}
	http.Handle("/", http.FileServer(http.Dir(env("NIJI_WEB", webRoot))))
	addr, err := listenAddress(os.Getenv("NIJI_LAN_IP"), env("NIJI_PORT", "8765"), env("NIJI_ADDR", "127.0.0.1:8080"))
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("待ち受けを開始できません（同じポートのサーバーが起動している場合は先に停止してください）: %v", err)
	}
	log.Printf("接続先 http://%s/", listener.Addr())
	if name := os.Getenv("NIJI_LOCAL_NAME"); name != "" && os.Getenv("NIJI_LAN_IP") != "" {
		_, port, _ := net.SplitHostPort(addr)
		log.Printf("家庭内Wi-Fi http://%s.local:%s/ （停止は Ctrl+C）", name, port)
	}
	srv := http.Server{Addr: addr, ReadTimeout: 10 * time.Second, WriteTimeout: 20 * time.Second, IdleTimeout: 60 * time.Second}
	log.Fatal(srv.Serve(listener))
}
