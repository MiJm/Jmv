package models

import (
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Reports struct {
	Model          `bson:"-" json:"-"` //model基类
	ReportId       bson.ObjectId       `bson:"_id,omitempty" json:"report_id"`           //id
	ReportType     uint8               `bson:"report_type" json:"report_type"`           //报表类型 0=追车，1=电话 2=违章
	ReportPlate    string              `bson:"report_plate" json:"report_plate"`         //报表车牌号
	ReportDataFrom uint8               `bson:"report_data_from" json:"report_data_from"` //报表源数据来源0=内部数据库 1=外部导入 @ReportType=0 仅追车有效
	ReportStatus   uint8               `bson:"report_status" json:"report_status"`       //报表状态0=查询中 1=查询成功 2=查询失败
	ReportShares   map[string]Shares   `bson:"report_shares" json:"report_shares"`       //报表分享列表  share_id作为键值
	ReportCreateAt int64               `bson:"report_createat" json:"report_createat"`   //报表查询时间
	ReportDoneAt   int64               `bson:"report_doneat" json:"report_doneat"`       //报表查询完成时间
	ReportDeleteAt int64               `bson:"report_deleteat" json:"report_deleteat"`   //报表删除时间 0=未删除
}

type Shares struct {
	ShareId       string `bson:"share_id" json:"share_id"`             //id bson.ObjectId hex
	ShareMobile   string `bson:"share_mobile" json:"share_mobile"`     //报表分享所属手机
	ShareFname    string `bson:"share_fname" json:"share_fname"`       //分享人姓名
	ShareViews    uint32 `bson:"share_views" json:"share_views"`       //报表分享人查看次数
	ShareCreateAt int64  `bson:"share_createat" json:"share_createat"` //报表分享时间
}

func (this *Reports) RemoveShare(shareId string) {
	if this.ReportShares != nil {
		delete(this.ReportShares, shareId)
	}
	this.Update()
}

func (this *Reports) Delete() error {
	if this.ReportId == EmptyId {
		return errors.New("id be empty")
	}
	return this.Collection(this).RemoveId(this.ReportId)
}

func (this *Reports) Insert() {
	this.ReportId = bson.NewObjectId()
	this.ReportCreateAt = time.Now().Unix()
	this.ReportPlate = "沪A123456"
	fmt.Println(this.Collection(this).Insert(*this))
}

func (this *Reports) Update() {
	if this.ReportId != EmptyId && this.ReportShares != nil {
		update := bson.M{"$set": *this}
		this.Collection(this).UpdateId(this.ReportId, update)
	}
}

func (this *Reports) List(query interface{}) (rs []*Reports) {
	if this.ReportId == EmptyId {
		this.Collection(this).Find(query).All(&rs)
	}
	return
}
