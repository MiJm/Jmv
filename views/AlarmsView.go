package views

import (
	"risk-ext/models"

	"github.com/kataras/iris"
)

type AlarmsView struct {
	Views
}

func (this *AlarmsView) Auth(ctx iris.Context) int {
	this.Views.Auth(ctx)
	var perms = PMS{
		"PUT":    MA{"CUSTOMER": A{1}},
		"GET":    MA{"CUSTOMER": A{1}},
		"POST":   MA{"CUSTOMER": A{1}},
		"DELETE": MA{"CUSTOMER": A{1}}}
	return this.CheckPerms(perms[ctx.Method()])
}

func (this *AlarmsView) Get(ctx iris.Context) (statuCode int, data M) {
	data = make(M)
	statuCode = 400
	alarmId := ctx.FormValue("alarmId")
	if alarmId == "" {
		data["code"] = 0
		data["error"] = "参数alarmId缺失"
		return
	}
	alarmInfo, err := new(models.Alarms).GetAlarmInfo(alarmId)
	if err != nil {
		data["code"] = 0
		data["error"] = err.Error()
		return
	}
	if alarmInfo.Alarm_read == 0 {
		alarmInfo.Alarm_read = 1
		alarmInfo.Update()
	}
	statuCode = 200
	data["code"] = 1
	data["alarmInfo"] = alarmInfo
	return
}

func (this *AlarmsView) Post(ctx iris.Context) (statuCode int, data M) {
	data = make(M)
	statuCode = 400
	deviceId := ctx.FormValue("deviceId")
	page := ctx.PostValueIntDefault("page", 1)
	pageSize := ctx.PostValueIntDefault("pageSize", 15)
	alarmList, count, err := new(models.Alarms).GetAlarmListByUserId(Session.Customer.UserId.Hex(), deviceId, page, pageSize)
	if err != nil {
		data["code"] = 0
		data["error"] = err.Error()
		return
	}
	unReadAlarmNum, err := new(models.Alarms).GetUnReadAlarmNums(deviceId)
	if err != nil {
		data["code"] = 0
		data["error"] = err.Error()
		return
	}
	statuCode = 200
	data["code"] = 1
	data["count"] = count
	data["unreadCount"] = unReadAlarmNum
	data["alarmList"] = alarmList
	return
}
