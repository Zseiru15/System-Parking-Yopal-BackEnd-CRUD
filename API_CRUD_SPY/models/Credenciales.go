package models

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/astaxie/beego/orm"
)

type Credenciales struct {
	Id                  int       `orm:"column(Id_Credenciales);pk;auto"`
	Contrasena          string    `orm:"column(Contrasena)"`
	Estado              bool      `orm:"column(Estado)"`
	FechaRegistro       time.Time `orm:"column(Fecha_Registro);type(timestamp with time zone);auto_now_add"`
	FechaModificacion   time.Time `orm:"column(Fecha_Modificacion);type(timestamp with time zone);auto_now"`
}

func (t *Credenciales) TableName() string {
	return "Credenciales"
}

func init() {
	orm.RegisterModel(new(Credenciales))
}

func HashContrasena(contrasena string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(contrasena), bcrypt.DefaultCost)
	return string(hash), err
}

func (c *Credenciales) VerificarContrasena(contrasena string) error {
	return bcrypt.CompareHashAndPassword([]byte(c.Contrasena), []byte(contrasena))
}

func AddCredenciales(m *Credenciales) (id int64, err error) {
	o := orm.NewOrm()

	if m.Contrasena != "" {
		hashed, err := HashContrasena(m.Contrasena)
		if err != nil {
			return 0, err
		}
		m.Contrasena = hashed
	}

	id, err = o.Insert(m)
	return
}

func GetCredencialesById(id int) (v *Credenciales, err error) {
	o := orm.NewOrm()
	v = &Credenciales{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAllCredenciales(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {

	o := orm.NewOrm()
	qs := o.QueryTable(new(Credenciales)).RelatedSel()
	// query k=v
	for k, v := range query {
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, v == "true" || v == "1")
		} else {
			qs = qs.Filter(k, v)
		}
	}

	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("invalid order, must be asc or desc")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(order) == 1 {
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("invalid order, must be asc or desc")
				}
				sortFields = append(sortFields, orderby)
			}
		} else {
			return nil, errors.New("sortby and order size mismatch")
		}
	}

	var l []Credenciales
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
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

func UpdateCredencialesById(m *Credenciales) error {
	o := orm.NewOrm()
	v := Credenciales{Id: m.Id}
	if err := o.Read(&v); err == nil {
		if m.Contrasena != "" {
			hashed, err := HashContrasena(m.Contrasena)
			if err != nil {
				return err
			}
			m.Contrasena = hashed
		}
		_, err := o.Update(m)
		return err
	}
	return errors.New("credencial no encontrada")
}

func DeleteCredenciales(id int) error {
	o := orm.NewOrm()
	if num, err := o.Delete(&Credenciales{Id: id}); err == nil && num > 0 {
		return nil
	} else {
		return err
	}
}
