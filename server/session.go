package main

import (
	"encoding/json"
	"fmt"
)

func (a *App) session(d Input, owner string, now int64) (Row, int, string) {
	if owner == "" || owner == "parent" {
		return nil, 403, "もういちど学習者をえらんでね"
	}
	if d.ID != owner {
		return nil, 409, "学習者をえらびなおしてね"
	}
	if d.Action == "home" {
		return Row{"elements": a.sql("SELECT coalesce(sum(reward),0) FROM run WHERE owner=" + q(owner)), "mastery": a.score(owner)}, 200, ""
	}
	active := " WHERE owner=" + q(owner) + " AND unit=" + q(a.unit.ID) + " AND ended IS NULL"
	id := a.sql("SELECT id FROM run" + active)
	if id == "" {
		if d.Action != "play" {
			return nil, 409, "ホームから開きなおしてね"
		}
		id = token()
		b, _ := json.Marshal(generate(a.unit.Rounds))
		a.sql(fmt.Sprintf("INSERT INTO run(id,owner,unit,started,targets,since) VALUES(%s,%s,%s,%d,%s,%d)", q(id), q(owner), q(a.unit.ID), now, q(string(b)), now))
	}
	row := a.rows("SELECT json_object('targets',json(targets),'step',step,'hint',hint,'since',since) FROM run" + active)[0]
	step := int(row["step"].(float64))
	target := int(row["targets"].([]any)[step].(float64))
	state := Row{"run": id, "step": step, "target": target, "hint": row["hint"]}
	if d.Action == "play" {
		return state, 200, ""
	}
	if d.Run != id || d.Step != step {
		return nil, 409, "島を開きなおしてね"
	}
	if d.Action == "hint" {
		a.sql("UPDATE run SET hint=1" + active)
		state["hint"] = 1
		return state, 200, ""
	}
	value, ok := evaluate(d.Parts, target)
	if value < 0 {
		return nil, 400, "パーツを確認してね"
	}
	ms := now - int64(row["since"].(float64))
	if ms > 180000 {
		ms = 180000
	}
	update := fmt.Sprintf("since=%d", now)
	if ok {
		step++
		update += ",step=step+1,hint=0"
		if step == a.unit.Rounds {
			update += fmt.Sprintf(",ended=%d,reward=%d", now, a.unit.Reward)
			state["reward"] = a.unit.Reward
		}
	}
	flag := 0
	if ok {
		flag = 1
	}
	a.sql(fmt.Sprintf("BEGIN IMMEDIATE; INSERT INTO answer VALUES(%s,%d,%d,%d,%d,%.0f,%d,%d); UPDATE run SET %s %s; COMMIT;", q(id), d.Step, target, value, flag, row["hint"], ms, now, update, active))
	state["correct"] = ok
	return state, 200, ""
}
