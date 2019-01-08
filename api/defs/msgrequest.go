package defs

type AppUserLogin struct {
	Deviceid string `json:"deviceid"`
}

type AppChangeUserName struct {
	Accesstoken string `json:"accesstoken"`
	Nickname    string `json:"nick_name"`
}

type AppUnlockBook struct {
	Accesstoken string `json:"accesstoken"`
	Barcode     string `json:"barcode"`
	BookId      string `json:"bookId"`
}

type AppGetVideoUrl struct {
	Accesstoken string  `json:"accesstoken"`
	BookId      string  `json:"bookId"`
}

type UserLogin struct {
	Username string `json:"user_name"`
	Password string `json:"pwd"`
}

type UserChangePassword struct {
	Username    string `json:"user_name"`
	OldPassword string `json:"old_pwd"`
	NewPassword string `json:"new_pwd"`
}

type AddBook struct {
	BookName string `json:"book_name"`
	BookId   string `json:"book_id"`
}

type UpdateBook struct {
	Id       string `json:"id"`
	BookName string `json:"book_name"`
	BookId   string `json:"book_id"`
}

type FindBarcodeList struct {
	BookId string `json:"book_id"`
}

type ChangeBarcodeStatus struct {
	BarcodeId string `json:"barcode_id"`
}

type AddBarcode struct {
	AddNum string `json:"add_num"`
	BookId string `json:"book_id"`
}
