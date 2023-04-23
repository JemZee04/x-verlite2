package controllers

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"test/models"

	beego "github.com/beego/beego/v2/server/web"
)

// UsersController operations for Users
type UsersController struct {
	beego.Controller
}

// URLMapping ...
//func (c *UsersController) URLMapping() {
//	c.Mapping("Post", c.Post)
//	c.Mapping("GetOne", c.GetOne)
//	c.Mapping("GetAll", c.GetAll)
//	c.Mapping("Put", c.Put)
//	c.Mapping("Delete", c.Delete)
//}

// Post ...
// @Title Post
// @Description create Users
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 201 {int} models.Users
// @Failure 403 body is empty
// @router / [post]
func (c *UsersController) Post() {
	//var v models.Users
	x := map[string]string{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &x); err == nil {
		u := models.Users{
			Email: x["Email"], FirstName: x["FirstName"], LastName: x["LastName"], Password: x["Password"],
			Photo: x["Photo"], ThirdName: x["ThirdName"],
		}
		switch x["Type"] {
		case "Customer":
			cu := models.Customer{
				Citizenship: x["Citizenship"], Company: x["Company"], TypeOrganization: x["TypeOrganization"],
			}
			if _, err := models.AddCustomer(&u, &cu); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = cu
			} else {
				c.Data["json"] = err.Error()
			}
		case "Employee":
			em := models.Employees{JobTitle: x["JobTitle"]}
			if _, err := models.AddEmployee(&u, &em); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = em
				//c.Data["json"]
			} else {
				c.Data["json"] = err.Error()
			}
		case "Contractor":
			co := models.Contractor{
				Citizenship: x["Citizenship"], Company: x["Company"], Face: x["Face"],
				FieldActivity: x["FieldActivity"], Resume: x["Resume"],
			}
			if _, err := models.AddContractor(&u, &co); err == nil {
				c.Ctx.Output.SetStatus(201)
				c.Data["json"] = co
				//c.Data["json"]
			} else {
				c.Data["json"] = err.Error()
			}
		default:
			c.Data["json"] = err.Error()

		}

	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Users by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Users
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsersController) Get() {
	idStr := c.Ctx.Input.Param(":uid")
	uid, _ := strconv.Atoi(idStr)
	v, i, err := models.GetUsersById(uid)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		data := map[string]interface{}{
			"user": v,
			"type": i,
		}
		c.Data["json"] = data

	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Users
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Users
// @Failure 403
// @router / [get]
func (c *UsersController) GetAll() {
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

// Put ...
// @Title Put
// @Description update the Users
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Users	true		"body for Users content"
// @Success 200 {object} models.Users
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UsersController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Users{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateUsersById(&v); err == nil {
			c.Data["json"] = "OK"
		} else {
			c.Data["json"] = err.Error()
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Users
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsersController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeleteUsers(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// @Title Login
// @Description Logs user into the system
// @Param	email		query 	string	true		"The email for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 403 user not exist
// @router /login [get]
func (u *UsersController) Login() {
	email := u.GetString("email")
	password := u.GetString("password")

	v, i, err := models.Login(email, password)

	if err != nil {
		u.Data["json"] = err.Error()
	} else {
		temp := strings.Split(reflect.TypeOf(i).String(), ".")[1]
		data := map[string]interface{}{
			"user": v,
			temp:   i,
		}
		u.Data["json"] = data

	}
	u.ServeJSON()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UsersController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJSON()
}
