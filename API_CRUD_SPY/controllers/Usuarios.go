package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/astaxie/beego/orm"

	"github.com/sena_2824182/System-Parking-Yopal-BackEnd-CRUD/API_CRUD_SPY/models"

	"github.com/astaxie/beego"
)

// UsuariosController operations for Usuarios
type UsuariosController struct {
	beego.Controller
}

// URLMapping ...
func (c *UsuariosController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Usuarios
// @Param	body		body 	models.Usuarios	true		"body for Usuarios content"
// @Success 201 {int} models.Usuarios
// @Failure 403 body is empty
// @router / [post]
func (c *UsuariosController) Post() {
    // Estructura temporal para recibir datos
    var requestData struct {
        models.Usuarios
        Contrasena string `json:"contrasena"` // Contraseña en texto plano
    }
    
    // Parsear datos de entrada
    if err := json.Unmarshal(c.Ctx.Input.RequestBody, &requestData); err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "message": "Formato de datos incorrecto",
        }
        c.Ctx.Output.SetStatus(400)
        c.ServeJSON()
        return
    }

    // Hashear la contraseña con bcrypt
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(requestData.Contrasena), 
        bcrypt.DefaultCost,
    )
    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "message": "Error al procesar contraseña",
        }
        c.Ctx.Output.SetStatus(500)
        c.ServeJSON()
        return
    }

    // Crear registro de credenciales
    credencial := models.Credenciales{
        Contrasena: string(hashedPassword),
        Estado:     true,
    }

    // Asignar al usuario
    requestData.Usuarios.IdContrasenaFk = &credencial
    requestData.Usuarios.Estado = true // Activar usuario por defecto

    // Guardar en base de datos
    if _, err := models.AddUsuarios(&requestData.Usuarios); err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "message": "Error al crear usuario: " + err.Error(),
        }
        c.Ctx.Output.SetStatus(500)
    } else {
        c.Data["json"] = map[string]interface{}{
            "success": true,
            "message": "Usuario registrado exitosamente",
        }
        c.Ctx.Output.SetStatus(201)
    }
    c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Usuarios by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Usuarios
// @Failure 403 :id is empty
// @router /:id [get]
func (c *UsuariosController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUsuariosById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]interface{}{
			"success":  true,
			"status":  200,
			"message": "consulta realizada correctamente",
			"data":    v}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Usuarios
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Usuarios
// @Failure 403
// @router / [get]
func (c *UsuariosController) GetAll() {
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

	l, err := models.GetAllUsuarios(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = map[string]interface{}{
			"success":  true,
			"status":  200,
			"message": "consulta realizada correctamente",
			"data":    l}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Usuarios
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Usuarios	true		"body for Usuarios content"
// @Success 200 {object} models.Usuarios
// @Failure 403 :id is not int
// @router /:id [put]
func (c *UsuariosController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "ID inválido, debe ser un número entero",
		}
		c.ServeJSON()
		return
	}

	var v models.Usuarios
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "JSON inválido: " + err.Error(),
		}
		c.ServeJSON()
		return
	}
	v.Id = id

	if err := models.UpdateUsuariosById(&v); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al actualizar: " + err.Error(),
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"status":  200,
			"message": "Actualización realizada correctamente",
			"data":    v,
		}
	}
	c.ServeJSON()
}


// Delete ...
// @Title Delete
// @Description delete the Usuarios
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *UsuariosController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  400,
			"message": "ID inválido, debe ser un número entero",
		}
		c.ServeJSON()
		return
	}

	if err := models.DeleteUsuarios(id); err != nil {
		c.Data["json"] = map[string]interface{}{
			"success": false,
			"status":  500,
			"message": "Error al eliminar: " + err.Error(),
		}
	} else {
		c.Data["json"] = map[string]interface{}{
			"success": true,
			"status":  200,
			"message": "Se eliminó correctamente",
			"id":      id,
		}
	}
	c.ServeJSON()
}

// 3. Añadir nuevo endpoint para buscar por email
// @Title GetByEmail
// @Description Obtiene usuario por email
// @Param   email  query  string  true  "Email a buscar"
// @Success 200 {object} models.Usuarios
// @Failure 404 Not found
// @router /by-email [get]
func (c *UsuariosController) GetByEmail() {
    email := c.GetString("email")
    if email == "" {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "message": "Email es requerido",
        }
        c.Ctx.Output.SetStatus(400)
        c.ServeJSON()
        return
    }

    o := orm.NewOrm()
    var usuario models.Usuarios
    
    // Buscar usuario con relaciones
    err := o.QueryTable("usuarios").
        Filter("Email", email).
        RelatedSel("IdContrasenaFk").
        RelatedSel("IdRolesFk").
        One(&usuario)

    if err != nil {
        c.Data["json"] = map[string]interface{}{
            "success": false,
            "message": "Usuario no encontrado",
        }
        c.Ctx.Output.SetStatus(404)
    } else {
        c.Data["json"] = map[string]interface{}{
            "success": true,
            "data":    usuario,
        }
    }
    c.ServeJSON()
}