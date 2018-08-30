package routers

import (
	"json2csv/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/json2csv",&controllers.CsvController{},"*:Json2Csv")
}
