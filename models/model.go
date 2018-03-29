package models

import (
	"encoding/json"
	"reflect"
	"risk-ext/config"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	colls   = make(map[string]*mgo.Collection)
	EmptyId bson.ObjectId //空objectId
)

type Model struct {
}

type Redis struct {
}

func initColl(coll interface{}) string {
	key := _disName(coll)
	if colls[key] == nil {
		colls[key] = config.Mongo.C(key)
	}
	return key
}

func (this *Redis) Map(key, field string, result interface{}) (err error) {

	data, err := config.Redis.HGet(key, field).Bytes()
	if err != nil {
		return
	}
	err = json.Unmarshal(data, result)
	return
}

func (this *Redis) Save(key, field string, result interface{}) (err error) {
	data, err := json.Marshal(result)
	if err != nil {
		return
	}
	strData := string(data)
	err = config.Redis.HSet(key, field, []byte(strData)).Err()
	return
}

func (this *Model) Collection(coll interface{}) (c *mgo.Collection) {
	key := initColl(coll)
	c = colls[key]
	return
}

func (this *Model) One(coll interface{}) {
	key := initColl(coll)
	c := colls[key]
	v := reflect.ValueOf(coll).Elem()

	if &v != nil {
		fieldName := strings.ToUpper(key[0:1]) + key[1:len(key)-1] + "Id"
		idobj := v.FieldByName(fieldName)
		if idobj.IsValid() {
			c.FindId(idobj.Interface()).One(coll)
		}

	}
}

func _disName(coll interface{}) string {
	collName := reflect.TypeOf(coll).String()
	collName = strings.ToLower(collName)
	strs := strings.Split(collName, ".")
	key := strs[len(strs)-1]
	return key
}
