package defs

type Login struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
	Msg       string `json:"login_msg"`
}

type BookManagerList struct {
	Books []*Book `json:"books"`
	Msg   string  `json:"books_msg"`
}

type FindBarcode struct {
	BarcodeIds []*BarcodeId `json:"barcodeIds"`
	Msg        string       `json:"barcode_msg"`
}

type Book struct {
	Id           int    `json:"id" db:"id"`
	BookName     string `json:"book_name" db:"book_name"`
	BookId       string `json:"book_id" db:"book_id"`
	BarcodeCount int    `json:"barcode_count" db:"barcode_count"`
}

type BarcodeId struct {
	Id         string `json:"id" db:"id"`
	BarcodeId  string `json:"barcode_id" db:"bar_code_id"`
	BookId     string `json:"book_id" db:"book_id"`
	AllCount   int    `json:"all_count" db:"all_count"`
	UseCount   int    `json:"use_count" db:"use_count"`
	Status     string `json:"status" db:"status"`
	CreateTime string `json:"create_time"`
}

type BarcodeCount struct {
	BarcodeAllCount int    `json:"all_count"`
	Putaway         int    `json:"putaway_count"`
	NoPut           int    `json:"noput_count"`
	BarcodeUseCount int    `json:"use_count"`
	Msg             string `json:"barcode_count_msg"`
}

type SimpleSession struct {
	Username string //login name
	TTL      int64
}

//
//
//type UserSession struct {
//	Username  string `json:"user_name"`
//	SessionId string `json:"session_id"`
//}
