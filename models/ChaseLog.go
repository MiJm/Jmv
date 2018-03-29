package models

import (
	"gopkg.in/mgo.v2/bson"
)

type Log struct {
	LogId         bson.ObjectId `bson:"_id,omitempty" json:"log_id"`            //id
	LogDate       uint32        `bson:"log_date" json:"log_date"`               //日期
	LogItem       string        `bson:"log_item" json:"log_item"`               //查询的项目
	LogCompany    string        `bson:"log_company" json:"log_company"`         //客户名称
	LogCompanyId  string        `bson:"log_company_id" json:"log_company_id"`   //客户ID
	LogDetail     int16         `bson:"log_detail" json:"log_detail"`           //描述
	LogOperator   string        `bson:"log_operator" json:"log_operator"`       //操作人
	LogOperatorId string        `bson:"log_operator_id" json:"log_operator_id"` //操作人ID
}
