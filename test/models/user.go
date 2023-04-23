package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
	"reflect"
	"strings"
)

type Users struct {
	Id        int    `orm:"column(id);pk;auto"`
	FirstName string `orm:"column(first_name)"`
	LastName  string `orm:"column(last_name)"`
	ThirdName string `orm:"column(third_name);null"`
	Photo     string `orm:"column(photo);null"`
	Email     string `orm:"column(email)"`
	Password  string `orm:"column(password)"`
}

type Employees struct {
	//User *Users `orm:"reverse(one)" json:"-"`
	Id       int    `orm:"column(user_id);pk"`
	JobTitle string `orm:"column(job_title);null"`
}

type Customer struct {
	//User *Users `orm:"reverse(one)" json:"-"`
	Id               int    `orm:"column(user_id);pk"`
	Citizenship      string `orm:"column(citizenship);null"`
	Company          string `orm:"column(company);null"`
	TypeOrganization string `orm:"column(type_organization);null"`
}

type Contractor struct {
	//User *Users `orm:"reverse(one)" json:"-"`
	Id            int    `orm:"column(user_id);pk"`
	Citizenship   string `orm:"column(citizenship);null"`
	Company       string `orm:"column(company);null"`
	FieldActivity string `orm:"column(field_activity);null"`
	Face          string `orm:"column(face);null"`
	Resume        string `orm:"column(resume);null"`
}

func (t *Users) TableName() string {
	return "users"
}

func init() {
	orm.RegisterModel(new(Users))
	orm.RegisterModel(new(Customer))
	orm.RegisterModel(new(Contractor))
	orm.RegisterModel(new(Employees))

}

// AddUsers insert a new Users into database and returns
// last inserted Id on success.
func AddUsers(m *Users) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)

	return
}

func AddCustomer(m *Users, c *Customer) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	c.Id = int(id)
	id, err = o.Insert(c)
	return
}
func AddEmployee(m *Users, c *Employees) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	c.Id = int(id)
	id, err = o.Insert(c)
	return
}
func AddContractor(m *Users, c *Contractor) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	c.Id = int(id)
	id, err = o.Insert(c)
	return
}

// GetUsersById retrieves Users by Id. Returns error if
// Id doesn't exist
func GetUsersById(id int) (v *Users, i interface{}, err error) {
	o := orm.NewOrm()
	v = &Users{Id: id}
	em := &Employees{Id: id}
	cu := &Customer{Id: id}
	co := &Contractor{Id: id}
	err = o.Read(v)
	if err != nil {
		return nil, nil, err
	}
	if err = o.Read(em); err == nil {
		return v, em, nil
	}
	if err = o.Read(cu); err == nil {
		return v, cu, nil
	}
	if err = o.Read(co); err == nil {
		return v, co, nil
	}

	//if err = o.Read(v); err == nil {
	//	return v, nil
	//}
	return nil, nil, err
}

// GetAllUsers retrieves all Users matches certain condition. Returns empty list if
// no records exist
func GetAllUsers(
	query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64,
) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Users))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Users
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsers updates Users by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsersById(m *Users) (err error) {
	o := orm.NewOrm()
	v := Users{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsers deletes Users by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsers(id int) (err error) {
	o := orm.NewOrm()
	v := Users{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Users{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func Login(email string, password string) (v *Users, i interface{}, err error) {
	o := orm.NewOrm()
	//err = orm.NewOrm().QueryTable("users").Filter("email", email).Filter("password", password).One(v)
	err = o.Raw("SELECT * FROM users WHERE email = ?", email).QueryRow(&v)
	//v = &Users{Email: email, Password: password}
	if err == nil {
		if v.Password != password {
			return nil, nil, errors.New("Неверный пароль")
		}
		em := &Employees{Id: v.Id}
		cu := &Customer{Id: v.Id}
		co := &Contractor{Id: v.Id}
		//err = o.Read(v)
		if err != nil {
			return nil, nil, err
		}
		if err = o.Read(em); err == nil {
			return v, em, nil
		}
		if err = o.Read(cu); err == nil {
			return v, cu, nil
		}
		if err = o.Read(co); err == nil {
			return v, co, nil
		}

	} /*else{
		return nil, nil, err
	}*/
	return nil, nil, err
}
