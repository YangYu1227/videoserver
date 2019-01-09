package dbops

import (
	"Study/book_server/api/defs"
	"Study/book_server/api/utils"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"yy/ossops"
)

func AppUserLogin(deviceId string) error {

	// 查询数据库是否含有传入的deviceId

	stmtOut, err := dbConn.Prepare("SELECT deviceid FROM app_user WHERE deviceid=?")
	if err != nil {
		return err
	}

	var sqlDeviceId string

	rows, err := stmtOut.Query(deviceId)
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(&sqlDeviceId); err != nil {
			return err
		}
	}

	defer stmtOut.Close()

	log.Printf("数据库查询: %s\n", sqlDeviceId)
	log.Printf("客户端传入: %s\n", sqlDeviceId)
	if sqlDeviceId == deviceId {
		log.Println("相等")
		return nil
	} else {
		// 插入数据
		stmtIns, err := dbConn.Prepare("INSERT INTO app_user (deviceid, unlocked_books) VALUES (?, ?)")
		if err != nil {
			return err
		}

		_, err = stmtIns.Exec(deviceId, "")
		if err != nil {
			return err
		}

		defer stmtIns.Close()
		return nil
	}
}

func AppUserChangeName(accesstoken, nickname string) error {
	stmtIns, err := dbConn.Prepare("update app_user set nick_name=? where deviceid =?")
	if err != nil {
		log.Printf("数据库操作，更新书籍，错误：%v", err)
		return err
	}

	_, err = stmtIns.Exec(nickname, accesstoken)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func AppUnlockBook(accesstoken, barcode, bookId string) (defs.AppUnlockBookResponse, error) {
	// 查询该用户是否已经激活过这本书

	var res defs.AppUnlockBookResponse

	log.Println("查询数据库 用户已经解锁的书籍ID")
	// 查询数据库 用户已经解锁的书籍ID
	stmtOut, err := dbConn.Prepare("SELECT unlocked_books FROM app_user WHERE deviceid=?")
	if err != nil {
		log.Println(err)
		return res, err
	}

	var unlockedBooks string

	rows, err := stmtOut.Query(accesstoken)
	if err != nil {
		log.Println(err)
		return res, err
	}

	for rows.Next() {
		if err := rows.Scan(&unlockedBooks); err != nil {
			log.Println(err)
			return res, err
		}
	}

	defer stmtOut.Close()

	// 如果用户已经解锁过书籍
	if unlockedBooks != "" {
		if strings.Contains(unlockedBooks, bookId) {
			res.Code = "1"
			res.Msg = "该设备已激活过"
			res.Status = 4

			return res, nil
		}
	}

	log.Println("查询Barcode所属的Id和使用次数")
	// 查询Barcode所属的Id和使用次数
	stmtOut2, err := dbConn.Prepare("SELECT barcode_id, is_use FROM barcode WHERE bar_code=?")
	if err != nil {
		log.Println(err)
		return res, err
	}

	var barcodeId string
	var isUse int

	rows2, err := stmtOut2.Query(barcode)
	if err != nil {
		log.Println(err)
		return res, err
	}

	for rows2.Next() {
		if err := rows2.Scan(&barcodeId, &isUse); err != nil {
			log.Println(err)
			return res, err
		}
	}

	defer stmtOut.Close()

	// 如果该条码使用次数超过3次
	if isUse >= 3 {
		res.Code = "1"
		res.Msg = "该激活码最多只能激活三台设备"
		res.Status = 2

		return res, nil
	}
	log.Printf("BarcodeId:%s\n", barcodeId)

	// 查询使用总数和上架状态
	stmtOut3, err := dbConn.Prepare("SELECT use_count, status, book_id FROM book_bar_code WHERE bar_code_id=?")
	if err != nil {
		return res, err
	}

	var status int
	var useCount int
	var sqlBookId string

	rows3, err := stmtOut3.Query(barcodeId)
	if err != nil {
		return res, err
	}

	for rows3.Next() {
		if err := rows3.Scan(&useCount, &status, &sqlBookId); err != nil {
			return res, err
		}
	}

	defer stmtOut.Close()

	log.Printf("上架状态：%d\n", status)

	if sqlBookId != bookId {
		res.Code = "1"
		res.Msg = "该激活码不属于这本书"
		res.Status = 5

		return res, nil
	}

	// 如果该条码所属的订单号未上架
	if status != 1 {
		res.Code = "1"
		res.Msg = "该激活码不存在或已失效"
		res.Status = 3

		return res, nil
	}

	// 更新 Barcode 中 该条码的使用次数
	stmtIns, err := dbConn.Prepare("update barcode set is_use=? where bar_code =?")
	if err != nil {
		log.Printf("数据库操作，更新书籍，错误：%v", err)
		return res, err
	}

	isUse += 1

	_, err = stmtIns.Exec(isUse, barcode)
	if err != nil {
		return res, err
	}

	defer stmtIns.Close()

	if useCount == 0 {
		useCount += 1

		// 更新 Book_bar_code 中 该条码所属的订单号的使用次数
		stmtIns2, err := dbConn.Prepare("update book_bar_code set use_count=? where bar_code_id =?")
		if err != nil {
			log.Printf("数据库操作，更新书籍，错误：%v", err)
			return res, err
		}

		_, err = stmtIns2.Exec(useCount, barcodeId)
		if err != nil {
			return res, err
		}

		defer stmtIns.Close()
	}

	// 更新 app_user 中该用户已解锁书籍的ID列表
	stmtIns3, err := dbConn.Prepare("update app_user set unlocked_books=? where deviceid=?")
	if err != nil {
		log.Printf("数据库操作，更新书籍，错误：%v", err)
		return res, err
	}

	if unlockedBooks != "" {
		unlockedBooks += "," + bookId
	} else {
		unlockedBooks = bookId
	}

	_, err = stmtIns3.Exec(unlockedBooks, accesstoken)
	if err != nil {
		return res, err
	}

	defer stmtIns.Close()

	res.Code = "1"
	res.Msg = "激活码有效，解锁成功"
	res.Status = 1

	return res, nil
}

func AppUploadVideo(accesstoken, bookId, filename string) error {

	// 先查找数据库，看是否已经存在了同一本书的视频文件，如果存在，将删除该文件

	vid, err := GetVideo(accesstoken, bookId)

	if err != nil {
		return err
	}

	if vid != "" {

		err = deleteVideo(vid)
		if err != nil {
			return err
		}

		stmtIns, err := dbConn.Prepare("update app_upload_video set video_vid=? where user_id=? && book_id=?")
		if err != nil {
			log.Printf("数据库操作错误：%v", err)
			return err
		}

		_, err = stmtIns.Exec(filename, accesstoken, bookId)
		if err != nil {
			return err
		}

		defer stmtIns.Close()

		return nil
	} else {

		stmtIns, err := dbConn.Prepare("INSERT INTO app_upload_video (user_id, book_id, video_vid) VALUES (?, ?, ?)")
		if err != nil {
			return err
		}

		_, err = stmtIns.Exec(accesstoken, bookId, filename)
		if err != nil {
			return err
		}

		defer stmtIns.Close()

		return nil
	}
}

func GetVideo(accesstoken, bookId string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT video_vid FROM app_upload_video WHERE user_id=? && book_id=?")
	if err != nil {
		return "", err
	}

	rows, err := stmtOut.Query(accesstoken, bookId)
	if err != nil {
		return "", err
	}

	var videoVid string

	for rows.Next() {
		if err := rows.Scan(&videoVid); err != nil {
			if err == sql.ErrNoRows {
				return "", nil
			} else {
				return "", err
			}
		}
	}

	defer stmtOut.Close()

	return videoVid, nil
}

func deleteVideo(vid string) error {
	ossfn := "videos/" + vid
	bn := "njg"
	ok := ossops.DeleteObject(ossfn, bn)

	if !ok {
		log.Printf("Deleting video error, oss operation failed")
		return errors.New("deleting video error")
	}

	return nil
}

//=============================================================================================

func GetBackgroundUserPassword(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT password FROM book_user WHERE name = ?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	defer stmtOut.Close()

	return pwd, nil
}

func ChangeBackgroundUserPassword(loginName string, pwd string) error {
	log.Println("数据库操作，更新后台用户密码")
	stmtIns, err := dbConn.Prepare("update book_user set password=?, name=?")

	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmtIns.Exec(pwd, loginName)
	if err != nil {
		return err
	}

	return nil
}

func AddBook(bookName, bookId string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, book_name, book_id, barcode_count FROM book WHERE book_name=?")
	if err != nil {
		log.Printf("%s", err)
		return "服务器内部出错", err
	}

	rows, err := stmtOut.Query(bookName)
	if err != nil {
		log.Printf("%s", err)
		return "服务器内部出错", err
	}
	var book defs.Book
	for rows.Next() {

		if err := rows.Scan(&book.Id, &book.BookName, &book.BookId, &book.BarcodeCount); err != nil {
			return "服务器内部出错", err
		}
		log.Printf("数据：%v", book)
	}

	if book.BookName == bookName || book.BookId == bookId {
		return "该书已存在", nil
	}

	defer stmtOut.Close()

	stmtIns, err := dbConn.Prepare("INSERT INTO book (book_name, book_id, barcode_count) VALUES (?, ?, ?)")
	if err != nil {
		return "服务器内部出错", err
	}

	_, err = stmtIns.Exec(bookName, bookId, 0)
	if err != nil {
		return "服务器内部出错", err
	}

	defer stmtIns.Close()
	return "添加成功", nil
}

func UpdateBook(bookName, bookId, id string, ) error {
	log.Println("数据库操作，更新书籍")
	stmtIns, err := dbConn.Prepare("update book set book_name=?, book_id=? where id =?")
	if err != nil {
		log.Printf("数据库操作，更新书籍，错误：%v", err)
		return err
	}

	intId, err := strconv.Atoi(id)

	if err != nil {
		log.Printf("string转换int出错：%v", err)
		return err
	}

	_, err = stmtIns.Exec(bookName, bookId, intId)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

func ListBook() ([]*defs.Book, error) {
	stmtOut, err := dbConn.Query("SELECT id, book_name, book_id, barcode_count FROM book")

	var res []*defs.Book

	if err != nil {
		return res, err
	}
	for stmtOut.Next() {
		var book defs.Book
		if err := stmtOut.Scan(&book.Id, &book.BookName, &book.BookId, &book.BarcodeCount); err != nil {
			return res, err
		}
		log.Printf("数据：%v", book)
		res = append(res, &book)
	}

	if err = stmtOut.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut.Close()

	return res, nil
}

func ListBarcodeId(bookId string) ([]*defs.BarcodeId, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, bar_code_id, all_count, use_count, status, create_time FROM book_bar_code WHERE book_id=?")

	var res []*defs.BarcodeId

	if err != nil {
		return res, err
	}

	rows, err := stmtOut.Query(bookId)

	for rows.Next() {
		var barcodeId defs.BarcodeId
		var time []uint8
		if err := rows.Scan(&barcodeId.Id, &barcodeId.BarcodeId, &barcodeId.AllCount, &barcodeId.UseCount, &barcodeId.Status, &time); err != nil {
			return res, err
		}
		log.Printf("数据：%v", barcodeId)
		barcodeId.CreateTime = utils.ByteToString(time)
		res = append(res, &barcodeId)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut.Close()

	return res, nil
}

func ChangeBarcodeIdStatus(barcodeId string) error {
	stmtOut, err := dbConn.Prepare("update book_bar_code set status=? where bar_code_id =?")
	if err != nil {
		log.Printf("数据库操作，更新状态，错误：%v", err)
		return err
	}

	_, err = stmtOut.Exec(1, barcodeId)
	if err != nil {
		return err
	}

	defer stmtOut.Close()

	return nil
}

func DeleteBarcode(barcodeId string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM barcode where barcode_id=? LIMIT 5000")
	if err != nil {
		log.Printf("数据库删除barcode数据操作错误：%v", err)
		return err
	}

	_, err = stmtOut.Exec(barcodeId)
	if err != nil {
		return err
	}

	defer stmtOut.Close()

	stmtInput, err := dbConn.Prepare("SELECT book_id, all_count FROM book_bar_code WHERE bar_code_id=?")
	if err != nil {
		log.Printf("数据库获取book_bar_code数据操作错误：%v", err)
		return err
	}

	rows, err := stmtInput.Query(barcodeId)
	if err != nil {
		return err
	}
	var barcodeCount int
	var bookId string
	for rows.Next() {
		if err := rows.Scan(&bookId, &barcodeCount); err != nil {
			return err
		}
	}

	defer stmtInput.Close()

	stmtOut2, err := dbConn.Prepare("DELETE FROM book_bar_code where bar_code_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	_, err = stmtOut2.Exec(barcodeId)
	if err != nil {
		return err
	}

	defer stmtOut2.Close()

	stmtInput2, err := dbConn.Prepare("SELECT barcode_count FROM book WHERE book_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	rows2, err := stmtInput2.Query(bookId)
	if err != nil {
		return err
	}
	var count int
	for rows2.Next() {
		if err := rows2.Scan(&count); err != nil {
			return err
		}
	}

	defer stmtInput2.Close()

	count -= barcodeCount

	stmtInput3, err := dbConn.Prepare("update book set barcode_count=? where book_id =?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	_, err = stmtInput3.Exec(count, bookId)
	if err != nil {
		return err
	}

	defer stmtInput3.Close()

	return nil
}

func AddBarcode(addNum int, bookId string) error {
	// 生成一个BarcodeId
	barcodeId := utils.BarcodeId()

	// 生成给定数量的Barcode
	barcodes := utils.MultipleBarcode(addNum)

	batchSql := "INSERT INTO barcode (barcode_id, bar_code, is_use) VALUES"

	for i := 0; i < len(barcodes); i++ {
		batchSql += " ('" + barcodeId + "', '" + barcodes[i] + "', 0),"
	}

	sqlS := batchSql[0 : len(batchSql)-1]
	sqlS += ";"

	_, err = dbConn.Query(sqlS)
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	stmtInput, err := dbConn.Prepare("INSERT INTO book_bar_code (bar_code_id, book_id, all_count, use_count, status) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	_, err = stmtInput.Exec(barcodeId, bookId, addNum, 0, 0)
	if err != nil {
		return err
	}

	defer stmtInput.Close()

	stmtInput2, err := dbConn.Prepare("SELECT barcode_count FROM book WHERE book_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	rows, err := stmtInput2.Query(bookId)
	if err != nil {
		return err
	}
	var barcodeCount int
	for rows.Next() {
		if err := rows.Scan(&barcodeCount); err != nil {
			return err
		}
	}

	defer stmtInput2.Close()

	barcodeCount += addNum

	stmtInput3, err := dbConn.Prepare("update book set barcode_count=? where book_id =?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return err
	}

	_, err = stmtInput3.Exec(barcodeCount, bookId)
	if err != nil {
		return err
	}

	defer stmtInput3.Close()

	return nil
}

func ExportBarcode(barcodeId string) ([]*string, error) {
	stmtOut, err := dbConn.Prepare("SELECT bar_code FROM barcode WHERE barcode_id=?")

	var barcode []*string

	if err != nil {
		return nil, err
	}

	rows, err := stmtOut.Query(barcodeId)

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}

		barcode = append(barcode, &code)
	}

	if err = rows.Err(); err != nil {
		panic(err.Error())
	}

	defer stmtOut.Close()

	return barcode, nil
}

func GetBookName(bookId string) (string, error) {
	stmtInput2, err := dbConn.Prepare("SELECT book_name FROM book WHERE book_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return "", err
	}

	rows, err := stmtInput2.Query(bookId)
	if err != nil {
		return "", err
	}
	var bookName string
	for rows.Next() {
		if err := rows.Scan(&bookName); err != nil {
			return "", err
		}
	}

	return bookName, nil
}

func GetBookId(barcodeId string) (string, error) {
	stmtInput2, err := dbConn.Prepare("SELECT book_id FROM book_bar_code WHERE bar_code_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return "", err
	}

	rows, err := stmtInput2.Query(barcodeId)
	if err != nil {
		return "", err
	}
	var bookId string
	for rows.Next() {
		if err := rows.Scan(&bookId); err != nil {
			return "", err
		}
	}

	return bookId, nil
}

func GetBarcodeCount(bookId string) (defs.BarcodeCount, error) {
	log.Printf("前端传来的book_id：%s\n", bookId)

	var res defs.BarcodeCount

	barcodeAllCount := 0
	putawayCount := 0
	noPutCount := 0
	barcodeUseCount := 0

	stmtOut, err := dbConn.Prepare("SELECT all_count FROM book_bar_code WHERE book_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return res, err
	}

	rows, err := stmtOut.Query(bookId)
	if err != nil {
		return res, err
	}
	var eachCount int
	for rows.Next() {
		if err := rows.Scan(&eachCount); err != nil {
			return res, err
		}
		barcodeAllCount += eachCount
	}

	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut.Close()

	log.Printf("当前解锁码总数：%d\n", barcodeAllCount)

	stmtOut2, err := dbConn.Prepare("SELECT all_count FROM book_bar_code WHERE book_id=? && status=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return res, err
	}

	rows2, err := stmtOut2.Query(bookId, 1)
	if err != nil {
		return res, err
	}
	var eachCount2 int
	for rows2.Next() {
		if err := rows2.Scan(&eachCount2); err != nil {
			return res, err
		}
		putawayCount += eachCount2
	}

	if err = rows2.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut2.Close()

	log.Printf("当前上架的解锁码数量：%d\n", putawayCount)

	// 未上架的解锁码数量
	noPutCount = barcodeAllCount - putawayCount

	log.Printf("当前未上架的解锁码数量：%d\n", noPutCount)

	stmtOut3, err := dbConn.Prepare("SELECT use_count FROM book_bar_code WHERE book_id=?")
	if err != nil {
		log.Printf("数据库操作错误：%v", err)
		return res, err
	}

	rows3, err := stmtOut3.Query(bookId)
	if err != nil {
		return res, err
	}
	var useCount int
	for rows3.Next() {
		if err := rows3.Scan(&useCount); err != nil {
			return res, err
		}
		barcodeUseCount += useCount
	}

	if err = rows3.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	defer stmtOut3.Close()

	log.Printf("当前已激活的解锁码数量：%d\n", barcodeUseCount)

	res.BarcodeAllCount = barcodeAllCount
	res.Putaway = putawayCount
	res.NoPut = noPutCount
	res.BarcodeUseCount = barcodeUseCount
	res.Msg = "获取成功"

	return res, nil
}
