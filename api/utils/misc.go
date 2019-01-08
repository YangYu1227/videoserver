package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"uuid"
	"yy/config"
)

func singleBarcode() string {
	uid := uuid.Rand()
	uid2 := uid.Hex()
	//fmt.Printf("原始UUID：%s\n", uid2)

	uid2 = strings.Replace(uid2, "-", "", -1)

	//fmt.Printf("去掉“-”之后的UUID：%s\n", uid2)

	uid3 := uid2[0:20]

	//fmt.Printf("取前20个字符的UUID：%s\n", uid3)

	return uid3
}

func MultipleBarcode(num int) (barcodes []string) {
	for i := 0; i < num; i++ {
		barcode := singleBarcode()
		barcodes = append(barcodes, barcode)
	}
	return
}

func BarcodeId() string {
	s := time.Now().Format("20060102150405")
	s = "B" + s
	//fmt.Printf("time:%v\n", s)
	return s
}

func ByteToString(bs []uint8) string {
	var ba []byte
	for _, b := range bs {
		ba = append(ba, byte(b))
	}
	return string(ba)
}

func NewUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func GetCurrentTimestampSec() int {
	ts, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))
	return ts
}

func SendDeleteVideoRequest(id string) {
	addr := config.GetLbAddr() + ":9001"
	url := "http://" + addr + "/video-delete-record/" + id
	_, err := http.Get(url)
	if err != nil {
		log.Printf("Sending deleting video request error: %s", err)
	}
}
