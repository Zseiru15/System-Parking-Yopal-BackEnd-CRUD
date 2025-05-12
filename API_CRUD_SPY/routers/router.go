// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/sena_2824182/System-Parking-Yopal-BackEnd-CRUD/API_CRUD_SPY/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",

		beego.NSNamespace("/parqueaderos_promociones",
			beego.NSInclude(
				&controllers.EstacionamientosPromocionesController{},
			),
		),

		beego.NSNamespace("/credenciales",
			beego.NSInclude(
				&controllers.CredencialesController{},
			),
		),

		beego.NSNamespace("/promociones",
			beego.NSInclude(
				&controllers.PromocionesController{},
			),
		),

		beego.NSNamespace("/parqueaderos",
			beego.NSInclude(
				&controllers.EstacionamientosController{},
			),
		),

		beego.NSNamespace("/comentarios",
			beego.NSInclude(
				&controllers.ComentariosController{},
			),
		),

		beego.NSNamespace("/roles",
			beego.NSInclude(
				&controllers.RolesController{},
			),
		),

		beego.NSNamespace("/pagos",
			beego.NSInclude(
				&controllers.PagosController{},
			),
		),

		beego.NSNamespace("/usuarios",
			beego.NSInclude(
				&controllers.UsuariosController{},
			),
		),

		beego.NSNamespace("/ranking",
			beego.NSInclude(
				&controllers.RankingController{},
			),
		),

		beego.NSNamespace("/slots",
			beego.NSInclude(
				&controllers.SlotsController{},
			),
		),

		beego.NSNamespace("/vehiculos",
			beego.NSInclude(
				&controllers.VehiculosController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
