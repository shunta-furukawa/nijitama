package main

import "testing"

func TestListenAddress(t *testing.T) {
	for _, tc := range []struct{ ip, port, want string }{
		{"", "8765", "127.0.0.1:8080"}, {"192.168.1.2", "8765", "192.168.1.2:8765"}, {"0.0.0.0", "8765", ""}, {"8.8.8.8", "8765", ""}, {"::1", "8765", ""}, {"192.168.1.2", "65536", ""},
	} {
		got, e := listenAddress(tc.ip, tc.port, "127.0.0.1:8080")
		if got != tc.want || (e != nil) != (tc.want == "") {
			t.Fatalf("%v: %v %v", tc, got, e)
		}
	}
}
