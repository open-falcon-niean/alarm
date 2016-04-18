package exco

import (
	"fmt"
	"log"
	"time"

	"github.com/open-falcon/common/model"
	"github.com/open-falcon/common/utils"
)

// TODO:
// 1. 报警状态，放到外部存储(redis or mysql)
// 2. 垃圾状态清理
var (
	MainExcoStatus = NewExpressionCounterStatus()
	SubStatus      = NewExpressionCounterStatus()
)

// 复合表达式 检测
// return: do_alarm 是否继续进行报警动作
func CheckExcoAlarm(event *model.Event, isHigh bool) (do_alarm bool) {
	expression := event.Expression
	if !(expression != nil && expression.Id > 0) {
		do_alarm = true
		return
	}

	counter := utils.PK(event.Endpoint, expression.Metric, event.PushedTags)
	main_key := excoKey(expression.Id, counter)
	log.Printf("id:%d,counter:%s,key:%s", expression.Id, counter, main_key)

	// 设置 辅助表达式-指标 的状态
	if SubMap.Exist(main_key) {
		log.Println("-1: set sub status", main_key)
		SubStatus.Set(main_key, event.Status)
	}

	// 表达式 不是主表达式，继续之后的报警动作
	if !MainSubMap.Exist(main_key) {
		log.Println(main_key, "not main expression")
		do_alarm = true
		return
	}

	// 走到这里，说明，一定是主表达式了
	if event.Status == "PROBLEM" {
		log.Println("1")
		subs, _ := MainSubMap.Get(main_key)
		if !checkAllAlarm(subs) {
			// 辅助表达式 没有全部报警，则，此次报警不触发
			log.Printf("skip alarm: main %s, %s", main_key, utils.UnixTsFormat(time.Now().Unix()))
			do_alarm = false
			return
		}
		MainExcoStatus.Set(main_key, event.Status)
	} else {
		log.Println("2")
		old, found := MainExcoStatus.Get(main_key)
		MainExcoStatus.Set(main_key, event.Status)
		if !(found && old == "PROBLEM") {
			log.Println("2.1")
			// 收到恢复报警，但历史上没有发生过报警，则放弃本次恢复报警
			do_alarm = false
			return
		}
	}

	do_alarm = true
	return
}

func excoKey(id int, counter string) string {
	return fmt.Sprintf("%d_%s", id, counter)
}

// 判断，是否所有表达式，都处于报警状态
func checkAllAlarm(subs map[string]interface{}) bool {
	for key, _ := range subs {
		status, found := SubStatus.Get(key)
		if !found {
			// 状态没有保存，说明，一定不处于报警状态
			log.Println("not found: ", key)
			return false
		}
		if status != "PROBLEM" {
			log.Println("not problem: ", key)
			return false
		}
	}

	return true
}
