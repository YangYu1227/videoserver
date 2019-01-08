package defs

type NormalResponse struct {
	Code string `json:"code"` //0：代表出错；1：代表成功；2：代表认证失败
	Msg  string `json:"message"`
}

type LoginResponse struct {
	Code string `json:"code"`
	Msg  Login  `json:"message"`
}

type BookManagerListResponse struct {
	Code string          `json:"code"` //0：代表出错；1：代表成功；2：代表认证失败
	Msg  BookManagerList `json:"message"`
}

type FindBarcodeListResponse struct {
	Code string      `json:"code"` //0：代表出错；1：代表成功；2：代表认证失败
	Msg  FindBarcode `json:"message"`
}

type BarcodeCountResponse struct {
	Code string       `json:"code"` //0：代表出错；1：代表成功；2：代表认证失败
	Msg  BarcodeCount `json:"message"`
}

//=======================================================================

type AppUserLoginResponse struct {
	Code        string `json:"code"` //0：代表出错；1：代表成功；
	Msg         string `json:"message"`
	AccessToken string `json:"accesstoken"`
}

type AppChangeUserNameResponse struct {
	Code string `json:"code"` //0：代表出错；1：代表成功；
	Msg  string `json:"message"`
}

type AppUnlockBookResponse struct {
	Code string `json:"code"` //0：代表出错；1：代表成功；
	Msg  string `json:"message"`
	// 1:激活码有效，解锁成功; 2:该激活码最多只能激活三台设备
	// 3:该激活码不存在或已失效; 4:该设备已激活过
	// 5:该激活码不属于这本书
	Status int `json:"status"`
}

type AppUploadVideoResponse struct {
	Code string `json:"code"` //0：代表出错；1：代表成功；
	Msg  string `json:"message"`
}

type AppGetVideoUrlResponse struct {
	Code string `json:"code"` //0：代表出错；1：代表成功；
	Msg  string `json:"message"`
	Url  string `json:"url"`
}
