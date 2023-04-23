package controllers

import (
	"encoding/json"
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"strings"
	"test/models"
)

type ChatController struct {
	beego.Controller
}

//type WebSocketController struct {
//	beego.
//}

func (c *ChatController) Post() {
	var ch models.Chats
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ch); err == nil {
		if _, err := models.AddChat(&ch); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = ch
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

func (c *ChatController) Get() {
	idStr := c.Ctx.Input.Param(":uid")
	uid, _ := strconv.Atoi(idStr)
	v, err := models.GetChatById(uid)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {

		c.Data["json"] = v

	}
	c.ServeJSON()
}

func (c *ChatController) GetAll() {
	//idStr := c.Ctx.Input.Param(":uid")
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	l, err := models.GetAllUsers(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = l
	}
	c.ServeJSON()
}
