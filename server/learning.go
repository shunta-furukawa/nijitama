package main

import (
	"crypto/rand"
	"math/big"
)

func generate(n int) []int {
	out := []int{}
	seen := map[int]bool{}
	for len(out) < n {
		v, e := rand.Int(rand.Reader, big.NewInt(198))
		if e != nil {
			panic(e)
		}
		x := int(v.Int64()) + 2
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}
func evaluate(p []int, target int) (int, bool) {
	if len(p) != 3 {
		return -1, false
	}
	for _, x := range p {
		if x < 0 || x > 9 {
			return -1, false
		}
	}
	v := 100*p[0] + 10*p[1] + p[2]
	return v, v == target
}
func mastery(rows []Row) int {
	if len(rows) == 0 {
		return 0
	}
	sum := 0
	for _, r := range rows {
		if r["ok"] == float64(1) {
			sum += 100
			if r["hint"] == float64(1) {
				sum -= 50
			}
		}
	}
	return sum / len(rows)
}
