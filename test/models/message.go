package models

import "github.com/beego/beego/v2/client/orm"

type Messages struct {
	Id     int64  `orm:"column(id);pk;auto" json:"id"`
	ChatId int    `orm:"column(chat_id)" json:"chatId"`
	UserId int    `orm:"column(user_id)" json:"userId"`
	Text   string `orm:"column(text)" json:"text"`
	//IfRead int    `orm:"column(if_read);null" json:"ifRead"`
}

func init() {
	orm.RegisterModel(new(Messages))
}
