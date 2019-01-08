package main

import (
	"Study/book_server/api/dbops"
	"Study/book_server/api/defs"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"uuid"
	"yy/ossops"
)

/*
Unity应用API
*/
func AppUserLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("APP用户登录")

	res, _ := ioutil.ReadAll(r.Body)

	ubody := &defs.AppUserLogin{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.AppUserLoginResponse{
			Code:        "0",
			Msg:         "JSON解析出错",
			AccessToken: "",
		})

		return
	}

	if err := dbops.AppUserLogin(ubody.Deviceid); err != nil {
		SendResponse(w, defs.AppUserLoginResponse{
			Code:        "0",
			Msg:         "服务器内部错误",
			AccessToken: "",
		})

		return
	}

	SendResponse(w, defs.AppUserLoginResponse{
		Code:        "1",
		Msg:         "登录成功",
		AccessToken: ubody.Deviceid,
	})
}

func AppChangeUserName(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("APP用户修改昵称")

	res, _ := ioutil.ReadAll(r.Body)

	ubody := &defs.AppChangeUserName{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.AppChangeUserNameResponse{
			Code: "0",
			Msg:  "JSON解析出错",
		})

		return
	}

	if err := dbops.AppUserChangeName(ubody.Accesstoken, ubody.Nickname); err != nil {
		SendResponse(w, defs.AppChangeUserNameResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	SendResponse(w, defs.AppChangeUserNameResponse{
		Code: "1",
		Msg:  "修改成功",
	})
}

func AppUserUnlocked(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	log.Println("APP用户解锁书籍")

	res, _ := ioutil.ReadAll(r.Body)

	ubody := &defs.AppUnlockBook{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.AppUnlockBookResponse{
			Code: "0",
			Msg:  "JSON解析出错",
		})

		return
	}
	resp, err := dbops.AppUnlockBook(ubody.Accesstoken, ubody.Barcode, ubody.BookId)

	if err != nil {
		SendResponse(w, defs.AppUnlockBookResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	SendResponse(w, resp)
}

func AppUploadVideo(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	//log.Println("steamserver的UploadHandler方法被调用")
	//log.Printf("%s\n", p.ByName("accesstoken&bookid"))

	log.Println("开始上传视频")
	strParams := p.ByName("accesstoken&bookid")

	par := strings.Split(strParams, "&")

	accesstoken := par[0]
	bookId := par[1]

	log.Printf("accesstoken: %s\n", accesstoken)
	log.Printf("bookId: %s\n", bookId)

	formFile, _, err := r.FormFile("video") // <from name="file">
	if err != nil {
		SendResponse(w, defs.AppUploadVideoResponse{
			Code: "0",
			Msg:  "获取文件错误",
		})
		return
	}
	defer formFile.Close()

	log.Println("获取上传的文件完成")

	uid := uuid.Rand()
	filename := uid.Hex()
	log.Printf("生成的VID：%s\n", filename)

	data, err := ioutil.ReadAll(formFile)
	if err != nil {
		log.Printf("Read file error: %v", err)
		SendResponse(w, defs.AppUploadVideoResponse{
			Code: "0",
			Msg:  "获取文件错误",
		})
		return
	}

	err = ioutil.WriteFile("./videos/"+filename, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		SendResponse(w, defs.AppUploadVideoResponse{
			Code: "0",
			Msg:  "写入文件错误",
		})
		return
	}

	log.Println("文件写入完成")
	log.Printf("文件名：%s", filename)
	ossfn := "videos/" + filename
	path := "./videos/" + filename
	bn := "yy-book-server"
	ret := ossops.UploadToOss(ossfn, path, bn)

	if !ret {
		SendResponse(w, defs.AppUploadVideoResponse{
			Code: "0",
			Msg:  "上传OSS失败",
		})
		return
	}
	log.Println("上传OSS完成")

	_ = os.Remove(path)

	err = dbops.AppUploadVideo(accesstoken, bookId, filename)
	if err != nil {
		SendResponse(w, defs.AppUploadVideoResponse{
			Code: "0",
			Msg:  "数据库出错",
		})
		return
	}
	log.Println("数据库更新完成")

	SendResponse(w, defs.AppUploadVideoResponse{
		Code: "1",
		Msg:  "文件上传成功",
	})

}

func AppGetVideoUrl(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)

	ubody := &defs.AppGetVideoUrl{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.AppGetVideoUrlResponse{
			Code: "0",
			Msg:  "JSON解析出错",
		})

		return
	}

	vid, err := dbops.GetVideo(ubody.Accesstoken, ubody.BookId)

	if err != nil {
		SendResponse(w, defs.AppGetVideoUrlResponse{
			Code: "0",
			Msg:  "数据库出错",
		})
	}

	log.Printf("数据库读取到的VID：%s\n", vid)

	url := "https://yy-book-server.oss-cn-hongkong.aliyuncs.com/videos/" + vid

	SendResponse(w, defs.AppGetVideoUrlResponse{
		Code: "1",
		Msg:  "获取成功",
		Url:  url,
	})
}

//==========================================================================================

/*
后台页面API
*/

func BackgroundLogin(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	res, _ := ioutil.ReadAll(r.Body)
	//log.Printf("数据包：%s", res)
	ubody := &defs.UserLogin{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		SendResponse(w, defs.LoginResponse{
			Code: "0",
			Msg: defs.Login{
				Success:   false,
				SessionId: "",
				Msg:       "JSON解析错误",
			}})
		return
	}

	//log.Printf("%s", ubody.Username)
	pwd, err := dbops.GetBackgroundUserPassword(ubody.Username)
	//log.Printf("Login pwd: %s", pwd)
	//log.Printf("Login body pwd: %s", ubody.Password)
	if err != nil || len(pwd) == 0 || pwd != ubody.Password {
		SendResponse(w, defs.LoginResponse{
			Code: "0",
			Msg: defs.Login{
				Success:   false,
				SessionId: "",
				Msg:       "账号密码错误",
			}})
		return
	}

	id := GenerateNewSessionId(ubody.Username)
	si := &defs.Login{Success: true, SessionId: id, Msg: "登录成功"}
	log.Printf("创建SessionId完成：%s\n", id)

	SendResponse(w, defs.LoginResponse{
		Code: "1",
		Msg:  *si,
	})
}

func ChangePassword(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})
		return
	}

	//log.Println("后台用户修改密码")
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.UserChangePassword{}
	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("%s", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})

		return
	}

	pwd, err := dbops.GetBackgroundUserPassword(ubody.Username)

	if err != nil || len(pwd) == 0 || pwd != ubody.OldPassword {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "原始密码错误",
		})

		return
	}

	err = dbops.ChangeBackgroundUserPassword(ubody.Username, ubody.NewPassword)

	if err != nil {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部出错",
		})

		return
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1", //代表成功
		Msg:  "更改密码成功",
	})
}

func BookManagerList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.BookManagerListResponse{
			Code: "2", //代表认证失败
			Msg: defs.BookManagerList{
				Books: nil,
				Msg:   "用户认证失败",
			},
		})

		return
	}

	log.Println("拉去book列表")

	vs, err := dbops.ListBook()
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		SendResponse(w, defs.BookManagerListResponse{
			Code: "0",
			Msg: defs.BookManagerList{
				Books: nil,
				Msg:   "服务器内部出错",
			},
		})
		return
	}

	vsi := &defs.BookManagerList{Books: vs, Msg: "获取成功"}

	SendResponse(w, defs.BookManagerListResponse{
		Code: "1",
		Msg:  *vsi,
	})
}

func AddBook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})
		return
	}

	log.Println("添加书籍")
	res, _ := ioutil.ReadAll(r.Body)
	ubody := &defs.AddBook{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})
		return
	}
	//log.Printf("获取到的POST表单数据：%v", ubody)
	str, err := dbops.AddBook(ubody.BookName, ubody.BookId)
	if err != nil {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  str,
		})
		return
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1",
		Msg:  "添加成功",
	})
}

func UpdateBook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})
		return
	}

	//log.Println("更新书籍")
	res, _ := ioutil.ReadAll(r.Body)
	//log.Printf("91 获取到的POST表单数据：%v", res)
	ubody := &defs.UpdateBook{}

	if err := json.Unmarshal(res, ubody); err != nil {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})

		return
	}
	//log.Printf("99 获取到的POST表单数据：%v", ubody)
	if err := dbops.UpdateBook(ubody.BookName, ubody.BookId, ubody.Id); err != nil {
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部出错",
		})
		return
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1",
		Msg:  "更新书籍成功",
	})
}

func FindBarcodeList(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.FindBarcodeListResponse{
			Code: "2", //代表认证失败
			Msg: defs.FindBarcode{
				BarcodeIds: nil,
				Msg:        "用户认证失败",
			},
		})
		return
	}

	log.Println("查找解锁码列表")
	res, _ := ioutil.ReadAll(r.Body)
	//log.Printf("91 获取到的POST表单数据：%v", res)
	ubody := &defs.FindBarcodeList{}

	if err := json.Unmarshal(res, ubody); err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.FindBarcodeListResponse{
			Code: "0",
			Msg: defs.FindBarcode{
				BarcodeIds: nil,
				Msg:        "JSON解析错误",
			},
		})

		return
	}

	//log.Printf("获取到的POST表单数据：%v", ubody)
	//log.Printf("BookId: %s", ubody.BookId)
	//log.Printf("BarcodeId: %s", ubody.BarcodeId)
	//log.Printf("PageNumber: %d", ubody.PageNumber)
	//log.Printf("Status: %s", ubody.Status)

	barcode, err := dbops.ListBarcodeId(ubody.BookId)
	if err != nil {
		log.Printf("Error in ListAllvideos: %s", err)
		SendResponse(w, defs.FindBarcodeListResponse{
			Code: "0",
			Msg: defs.FindBarcode{
				BarcodeIds: nil,
				Msg:        "服务器内部出错",
			},
		})
		return
	}

	barcodes := &defs.FindBarcode{BarcodeIds: barcode, Msg: "获取成功"}
	SendResponse(w, defs.FindBarcodeListResponse{
		Code: "1",
		Msg:  *barcodes,
	})
}

func ChangeBarcodeStatus(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})

		return
	}

	//log.Println("更改Barcode为上架")

	res, _ := ioutil.ReadAll(r.Body)
	ubady := &defs.ChangeBarcodeStatus{}
	if err := json.Unmarshal(res, ubady); err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})

		return
	}

	//log.Printf("BarcodeId：%s", ubady.BarcodeId)

	err := dbops.ChangeBarcodeIdStatus(ubady.BarcodeId)

	if err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部出错",
		})
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1",
		Msg:  "更改成功",
	})
}

func AddBarcode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})
		return
	}

	log.Println("添加Barcode")

	res, _ := ioutil.ReadAll(r.Body)

	ubady := &defs.AddBarcode{}
	if err := json.Unmarshal(res, ubady); err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})

		return
	}

	//log.Printf("AddNum: %s, BookId: %s", ubady.AddNum, ubady.BookId)

	intNum, err := strconv.Atoi(ubady.AddNum)
	if err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	err = dbops.AddBarcode(intNum, ubady.BookId)
	if err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1",
		Msg:  "添加解锁码成功",
	})
}

func DeleteBarcode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.NormalResponse{
			Code: "2", //代表认证失败
			Msg:  "用户认证失败",
		})
		return
	}

	log.Println("删除Barcode")

	res, _ := ioutil.ReadAll(r.Body)
	ubady := &defs.ChangeBarcodeStatus{}
	if err := json.Unmarshal(res, ubady); err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "JSON解析错误",
		})
		return
	}

	err := dbops.DeleteBarcode(ubady.BarcodeId)

	if err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	SendResponse(w, defs.NormalResponse{
		Code: "1",
		Msg:  "删除解锁码成功",
	})
}

func ExportBarcode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	//if !ValidateUser(w, r) {
	//	SendResponse(w, defs.NormalResponse{
	//		Code: "2", //代表认证失败
	//		Msg:  "用户认证失败",
	//	})
	//	return
	//}

	log.Println("导出解锁码到Excel")

	res, _ := ioutil.ReadAll(r.Body)

	test2 := strings.Split(string(res), "&")

	code := test2[0]
	bId := test2[1]

	barcodeId := strings.Split(code, "=")[1]
	bookId := strings.Split(bId, "=")[1]

	log.Println(barcodeId)
	log.Println(bookId)

	var xlsxfile *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell

	xlsxfile = xlsx.NewFile()
	sheet, err := xlsxfile.AddSheet("图书解锁码表单")

	if err != nil {
		fmt.Printf("添加sheet出错：%v\n", err)
	}
	_ = sheet.SetColWidth(0, 0, 20)
	_ = sheet.SetColWidth(1, 1, 30)
	_ = sheet.SetColWidth(2, 2, 30)
	_ = sheet.SetColWidth(3, 3, 40)
	row = sheet.AddRow()
	cell = row.AddCell()
	cell.Value = "图书编号"
	cell = row.AddCell()
	cell.Value = "图书名称"
	cell = row.AddCell()
	cell.Value = "解锁码编号"
	cell = row.AddCell()
	cell.Value = "图书解锁码"

	bookName, err := dbops.GetBookName(bookId)
	if err != nil {
		log.Printf("数据库操作错误：%v\n", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})
		return
	}

	barcodes, err := dbops.ExportBarcode(barcodeId)

	if err != nil {
		log.Printf("数据库操作错误：%v\n", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})
		return
	}

	for _, v := range barcodes {
		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = bookId
		cell = row.AddCell()
		cell.Value = bookName
		cell = row.AddCell()
		cell.Value = barcodeId
		cell = row.AddCell()
		cell.Value = *v
	}

	fileName := "图书解锁码统计表单.xlsx"
	err = xlsxfile.Save("./" + fileName)
	if err != nil {
		log.Printf("存储xlsx文件出错：%v\n", err)
		SendResponse(w, defs.NormalResponse{
			Code: "0",
			Msg:  "服务器内部错误",
		})

		return
	}

	file, err := os.Open("./" + fileName)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}
	b, _ := ioutil.ReadAll(file)
	w.Header().Add("Content-Disposition", "attachment")
	//w.Header().Add("Content-Type", "application/vnd.ms-excel")
	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	_, _ = w.Write(b)
}

func BarcodeCount(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if !ValidateUser(w, r) {
		SendResponse(w, defs.BarcodeCountResponse{
			Code: "2", //代表认证失败
			Msg: defs.BarcodeCount{
				Msg: "用户认证失败",
			},
		})
		return
	}

	log.Println("获取Barcode数量数据")

	res, _ := ioutil.ReadAll(r.Body)
	ubady := &defs.FindBarcodeList{}
	if err := json.Unmarshal(res, ubady); err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.BarcodeCountResponse{
			Code: "0",
			Msg: defs.BarcodeCount{
				Msg: "JSON解析错误",
			},
		})

		return
	}

	barcodeCount, err := dbops.GetBarcodeCount(ubady.BookId)
	if err != nil {
		log.Printf("错误：%v", err)
		SendResponse(w, defs.BarcodeCountResponse{
			Code: "0",
			Msg: defs.BarcodeCount{
				Msg: "服务器内部错误",
			},
		})
		return
	}

	SendResponse(w, defs.BarcodeCountResponse{
		Code: "1",
		Msg:  barcodeCount,
	})

}
