package controllers

import (
	"github.com/astaxie/beego"
	"github.com/yukithm/json2csv"

	// "github.com/wildducktheories/go-csv"
	"bytes"
	"encoding/json"
	_ "io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	simplejson "github.com/bitly/go-simplejson"
)

type CsvController struct {
	beego.Controller
}

var headerStyleTable = map[string]json2csv.KeyStyle{
	"jsonpointer": json2csv.JSONPointerStyle,
	"slash":       json2csv.SlashStyle,
	"dot":         json2csv.DotNotationStyle,
	"dot-bracket": json2csv.DotBracketStyle,
}

func (this *CsvController) Json2Csv() {
	beego.Info(999999999999)
	// beego.Info(this.Ctx.Input.RequestBody)
	var v interface{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &v)
	// beego.Info(v)
	results := []json2csv.KeyValue{}
	obj, err := json2obj(string(this.Ctx.Input.RequestBody))
	// beego.Info(obj)
	t := time.Now()
	str_t := t.Format("2006-01-02 15:04:05")
	filename := strings.Join(strings.Split(strings.Split(str_t, ":")[0], " "), "-") + ".csv"

	if err != nil {
		this.Data["json"] = &map[string]interface{}{"status": "5000", "msg": "解析失败"}
		this.ServeJSON()
		return
	}

	arraykeylist := getKeylist(obj)
	//
	resultlist := make([][]string, 0)

	//
	if ok, _ := PathExists(filename); ok {
		resultlist := handleResult(string(this.Ctx.Input.RequestBody), resultlist, arraykeylist)
		file, _ := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		beego.Info("文件已经存在")
		Writefile(file, resultlist)
	} else {
		resultlist = append(resultlist, arraykeylist)
		resultlist = handleResult(string(this.Ctx.Input.RequestBody), resultlist, arraykeylist)
		file, _ := os.Create(filename)
		Writefile(file, resultlist)

	}

	if err != nil {
		panic(err)
	}
	// headerStyle := headerStyleTable["slash"]
	// csv := json2csv.NewCSVWriter(os.Stdout)
	// csv := json2csv.NewCSVWriter(file)
	// csv.HeaderStyle = headerStyle
	// csv.Transpose =true

	//

	results, err = json2csv.JSON2CSV(obj)

	if err != nil {
		log.Fatal(err)
	}
	if len(results) == 0 {
		this.Data["json"] = &map[string]interface{}{"status": "5000", "msg": "解析失败"}
		this.ServeJSON()
		return
	}

	// if err := csv.WriteCSV(results); err != nil {
	// 		this.Data["json"] = &map[string]interface{}{"status":"5000","msg":err}
	// }

	// go func() {
	// 	c, err := io.Copy(file,os.Stdout)
	// 	beego.Info(c,err)
	// }()
	this.Data["json"] = &map[string]interface{}{"status": "200", "msg": "装换成功"}
	this.ServeJSON()

}

func json2obj(jsonstr string) (interface{}, error) {
	r := bytes.NewReader([]byte(jsonstr))
	d := json.NewDecoder(r)
	d.UseNumber()
	var obj interface{}
	if err := d.Decode(&obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func getKeylist(obj interface{}) []string {
	// fklist := make([]interface{},0)
	fvlist := make([]interface{}, 0)
	// sklist := make([]interface{},0)
	// svlist := make([]interface{},0)
	// tklist := make([]interface{},0)
	arraykeylist := make([]string, 0)
	arrayvalist := make([][]interface{}, 0)
Loop:
	for _, v := range obj.(map[string]interface{}) {
		if _, ok := v.(string); ok {
			beego.Info("11111")
			fvlist = append(fvlist, v)
		} else if _, ok := v.([]interface{}); ok {
			beego.Info("2222")
			flag := true
			for _, v1 := range v.([]interface{}) {
				for kk1, vv1 := range v1.(map[string]interface{}) {

					// beego.Info(kk1)
					// beego.Info(vv1)

					if _, ok := vv1.(map[string]interface{}); ok {
						templist := make([]interface{}, 0)
						for kkk1, vvv1 := range vv1.(map[string]interface{}) {
							// beego.Info(kk1+"/"+kkk1)

							templist = append(templist, vvv1)
							if flag {
								arraykeylist = append(arraykeylist, kk1+"/"+kkk1)
							}

							arrayvalist = append(arrayvalist, templist)
							// beego.Info(vvv1)
						}

					} else {
						if flag {
							arraykeylist = append(arraykeylist, kk1)
						}

					}
				}
				flag = false
				break Loop
			}
		} else {
			fvlist = append(fvlist, v)
		}

	}

	// beego.Info(fklist,fvlist)
	// beego.Info(sklist,svlist)
	// beego.Info(tklist)
	sort.Strings(arraykeylist)
	sort.Sort(sort.Reverse(sort.StringSlice(arraykeylist)))

	return arraykeylist
}

func handleResult(result string, resultlist [][]string, arraykeylist []string) [][]string {
	res, err := simplejson.NewJson([]byte(result))
	if err != nil {
		beego.Warn("%v\n", err)
		return [][]string{}
	}

	if rows, err := res.Get("alerts").Array(); err == nil {
		for _, row := range rows {
			templist := make([]string, 0)
			for _, v := range arraykeylist {
				if ok := strings.Contains(v, "/"); ok {
					a := strings.Split(v, "/")
					k1, k2 := a[0], a[1]
					if vv, ok := row.(map[string]interface{})[k1].(map[string]interface{})[k2].(string); ok {
						templist = append(templist, vv)
					} else {
						templist = append(templist, "")
					}

					// beego.Info(row.(map[string]interface{})[k1].(map[string]interface{})[k2])
				} else {
					if tt, ok := row.(map[string]interface{})[v].(string); ok {
						templist = append(templist, tt)
					} else {
						templist = append(templist, "")
					}

				}
			}
			resultlist = append(resultlist, templist)
			beego.Info(row.(map[string]interface{})[""])
		}
	}
	// beego.Info(resultlist)
	return resultlist
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Writefile(file *os.File, resultlist [][]string) {
	for _, rl := range resultlist {
		record := strings.Join(rl, ",")
		beego.Info(record)
		file.Write([]byte(record))
		file.Write([]byte("\n"))
	}
	defer file.Close()
}
