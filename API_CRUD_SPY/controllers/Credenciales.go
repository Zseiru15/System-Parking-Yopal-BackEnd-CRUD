package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/sena_2824182/System-Parking-Yopal-BackEnd-CRUD/API_CRUD_SPY/models"
)

type CredencialesController struct {
	beego.Controller
}

func (c *CredencialesController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @router / [post]
func (c *CredencialesController) Post() {
	var v models.Credenciales
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if _, err := models.AddCredenciales(&v); err == nil {
			c.Ctx.Output.SetStatus(201)
			c.Data["json"] = map[string]interface{}{
				"success": true,
				"status":  201,
				"message": "creación generada correctamente",
				"data":    v}
		} else {
			c.CustomAbort(500, err.Error())
		}
	} else {
		c.CustomAbort(400, err.Error())
	}
	c.ServeJSON()
}

// @router /:id [get]
func (c *CredencialesController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(400, "ID inválido")
		return
	}
	v, err := models.GetCredencialesById(id)
	if err != nil {
		c.CustomAbort(404, err.Error())
		return
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"status":  200,
		"message": "consulta realizada correctamente",
		"data":    v}
	c.ServeJSON()
}

// @router / [get]
func (c *CredencialesController) GetAll() {
	var fields, sortby, order []string
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
				c.CustomAbort(400, "query inválida")
			}
			query[kv[0]] = kv[1]
		}
	}

	l, err := models.GetAllCredenciales(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.CustomAbort(500, err.Error())
	}
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"status":  200,
		"message": "consulta realizada correctamente",
		"data":    l}
	c.ServeJSON()
}

// @router /:id [put]
func (c *CredencialesController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(400, "ID inválido")
	}
	v := models.Credenciales{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		if err := models.UpdateCredencialesById(&v); err == nil {
			c.Data["json"] = map[string]interface{}{
				"success": true,
				"status":  200,
				"message": "actualización realizada correctamente",
				"data":    v}
		} else {
			c.CustomAbort(500, err.Error())
		}
	} else {
		c.CustomAbort(400, err.Error())
	}
	c.ServeJSON()
}

// @router /:id [delete]
func (c *CredencialesController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.CustomAbort(400, "ID inválido")
	}
	if err := models.DeleteCredenciales(id); err == nil {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"status":  200,
			"message": "se eliminó correctamente",
			"id":      id}
	} else {
		c.CustomAbort(500, err.Error())
	}
	c.ServeJSON()
}
