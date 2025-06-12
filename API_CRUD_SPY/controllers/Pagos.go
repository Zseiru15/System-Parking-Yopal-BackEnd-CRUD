package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/sena_2824182/System-Parking-Yopal-BackEnd-CRUD/API_CRUD_SPY/models"

	"github.com/astaxie/beego"
)

// PagosController operations for Pagos
type PagosController struct {
	beego.Controller
}

// URLMapping ...
func (c *PagosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Pagos
// @Param	body		body 	models.Pagos	true		"body for Pagos content"
// @Success 201 {int} models.Pagos
// @Failure 403 body is empty
// @router / [post]
func (c *PagosController) Post() {
	var v models.Pagos
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// Si no se especifica, asignar Status = true (puedes cambiar a false si lo prefieres)
		v.Status = true

		if _, err := models.AddPagos(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{
				"succes":  true,
				"status":  201,
				"message": "creacion generada correctamente",
				"data":    v}
		} else {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"status":  500,
				"message": err.Error(),
			}
			c.Ctx.Output.SetStatus(500)
		}
	} else {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "Error al decodificar el cuerpo de la solicitud",
			"error":   err.Error(),
		}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Pagos by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Pagos
// @Failure 403 :id is empty
// @router /:id [get]
func (c *PagosController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPagosById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]interface{}{
			"succes":  true,
			"status":  200,
			"message": "consulta realizada correctamente",
			"data":    v}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Pagos
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Pagos
// @Failure 403
// @router / [get]
func (c *PagosController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Ctx.Output.SetStatus(400)
				c.Data["json"] = map[string]interface{}{
					"success": false,
					"status":  400,
					"message": "Error: invalid query key/value pair",
				}
				c.ServeJSON()
				return
			}
			query[kv[0]] = kv[1]
		}
	}

	l, err := models.GetAllPagos(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": err.Error(),
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"status":  200,
			"message": "Consulta realizada correctamente",
			"data":    l,
		}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Pagos
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Pagos	true		"body for Pagos content"
// @Success 200 {object} models.Pagos
// @Failure 403 :id is not int
// @router /:id [put]
func (c *PagosController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v := models.Pagos{Id: id}

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		// Aquí no es necesario validar strings, ya que Status es bool

		if err := models.UpdatePagosById(&v); err == nil {
			c.Data["json"] = map[string]interface{}{
				"succes":  true,
				"status":  200,
				"message": "actualizacion realizada correctamente",
				"data":    v}
		} else {
			c.Data["json"] = map[string]interface{}{
				"success": false,
				"status":  500,
				"message": err.Error(),
			}
			c.Ctx.Output.SetStatus(500)
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": err.Error(),
		}
		c.Ctx.Output.SetStatus(400)
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Pagos
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *PagosController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	if err := models.DeletePagos(id); err == nil {
		c.Data["json"] = "OK"
		c.Data["json"] = map[string]interface{}{
			"succes":                true,
			"status":                200,
			"message":               "se elimino correctamente",
			"dato eliminado con id": id}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
