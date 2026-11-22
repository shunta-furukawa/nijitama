package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func token() string {
	b := make([]byte, 24)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func digest(key, id, pin string) string {
	h := hmac.New(sha256.New, []byte(key))
	fmt.Fprint(h, id, "/", pin)
	return hex.EncodeToString(h.Sum(nil))
}
func (a *App) authenticate(w http.ResponseWriter, d Input, now int64) bool {
	ticket := func(id string) {
		t := token()
		a.sql(fmt.Sprintf("INSERT INTO login VALUES(%s,%s,%d)", q(t), q(id), now+8*3600000))
		http.SetCookie(w, &http.Cookie{Name: "ticket", Value: t, Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode})
	}
	key := d.Action
	if d.Action == "select" {
		key += "/" + d.ID
	}
	recent := []time.Time{}
	for _, t := range a.failures[key] {
		if time.Since(t) < time.Minute {
			recent = append(recent, t)
		}
	}
	if len(recent) >= 8 {
		http.Error(w, "少し待ってから試してね", 429)
		return false
	}
	want, got, id := a.pin, d.Pin, "parent"
	if d.Action == "select" {
		id = d.ID
		want = a.sql("SELECT pin FROM learner WHERE id=" + q(id))
		got = digest(a.pin, id, d.Pin)
	}
	if subtle.ConstantTimeCompare([]byte(want), []byte(got)) != 1 {
		a.failures[key] = append(recent, time.Now())
		http.Error(w, "番号を確認してね", 403)
		return false
	}
	delete(a.failures, key)
	ticket(id)
	return true
}
