package main

import (
	"bufio"
	"encoding/json"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Drive the real API through pipes so the UI test can also run without LAN
// or permission to open a listening socket.
func TestScreens(t *testing.T) {
	node := os.Getenv("NIJI_TEST_NODE")
	if node == "" {
		t.Skip("set NIJI_TEST_NODE to enable React DOM integration")
	}
	a := App{db: filepath.Join(t.TempDir(), "family.db"), pin: token(), unit: Unit{"es4-decimal-bridge-v1", 5, 3}, failures: map[string][]time.Time{}}
	s, _ := files.ReadFile("schema.sql")
	a.sql(string(s))
	cmd := exec.Command(node, "tests/screens.mjs")
	cmd.Dir = "../web"
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "NIJI_TEST_PIN="+a.pin)
	input, e := cmd.StdinPipe()
	if e != nil {
		t.Fatal(e)
	}
	output, e := cmd.StdoutPipe()
	if e != nil {
		t.Fatal(e)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer cmd.Process.Kill()
	ticket := ""
	scan := bufio.NewScanner(output)
	for scan.Scan() {
		req := httptest.NewRequest("POST", "/api", strings.NewReader(scan.Text()))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", ticket)
		w := httptest.NewRecorder()
		a.api(w, req)
		if cookies := w.Result().Cookies(); len(cookies) > 0 {
			ticket = cookies[0].Name + "=" + cookies[0].Value
		}
		if e := json.NewEncoder(input).Encode(Row{"status": w.Code, "body": w.Body.String()}); e != nil {
			t.Fatal(e)
		}
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		t.Fatal(err)
	}
}
