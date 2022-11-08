package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/beego/beego/v2/client/orm"
)

type Good struct {
	Id     int     `orm:"column(id);auto" description:"商品表主键" json:"id"`
	Name   string  `orm:"column(name);size(45);null" description:"商品名称" json:"name"`
	Cover  string  `orm:"column(cover);size(45);null" description:"商品封面图" json:"cover"`
	Image1 string  `orm:"column(image1);size(45);null" description:"商品详情图1" json:"image_1"`
	Image2 string  `orm:"column(image2);size(45);null" description:"商品详情图2" json:"image_2"`
	Price  string `orm:"column(price);null" description:"商品价格" json:"price"`
	Intro  string  `orm:"column(intro);size(300);null" description:"商品描述" json:"intro"`
	Stock  int     `orm:"column(stock);null" description:"商品库存" json:"stock"`
	TypeId int     `orm:"column(type_id);null" description:"商品类型" json:"typeId"`
}

func (t *Good) TableName() string {
	return "good"
}

func init() {
	orm.RegisterModel(new(Good))
}

// AddGood insert a new Good into database and returns
// last inserted Id on success.
func AddGood(m *Good) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetGoodById retrieves Good by Id. Returns error if
// Id doesn't exist
func GetGoodById(id int) (v *Good, err error) {
	o := orm.NewOrm()
	v = &Good{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllGood retrieves all Good matches certain condition. Returns empty list if
// no records exist
func GetAllGood(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Good))
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

	var l []Good
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

// UpdateGood updates Good by Id and returns error if
// the record to be updated doesn't exist
func UpdateGoodById(m *Good) (err error) {
	o := orm.NewOrm()
	v := Good{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteGood deletes Good by Id and returns error if
// the record to be deleted doesn't exist
func DeleteGood(id int) (err error) {
	o := orm.NewOrm()
	v := Good{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Good{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetGoodByName(name string) (v *Good, err error) {
	o := orm.NewOrm()
	v = &Good{Name: name}
	if err = o.Read(v, "name"); err == nil {
		return v, nil
	}
	return nil, err
}
