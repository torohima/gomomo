package momo

import (
	"api_bank/momo/momoencoder"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type simpleResponse struct {
	Error float64
	Msg   string
}

type momoApiCfg struct {
	AppVer    int
	AppCode   string
	Device    string
	ApiAction map[string]string
}

type momoApi struct {
	Phone             string
	Password          string
	AuthToken         string
	RequestEncryptKey string
	Imei              string
	Config            momoApiCfg
}

func NewMomo(phone string, password string) *momoApi {
	apiAction := map[string]string{
		"QUERY_TRAN_HIS_MSG": "https://api.momo.vn/sync/transhis/browse",                      // check lịch sử giao dịch
		"M2MU_CONFIRM":       "https://owa.momo.vn/api/sync",                                  // chuyển tiền cho user momo
		"M2MU_INIT":          "https://owa.momo.vn/api/sync",                                  // tạo phiên chuyển tiền user momo
		"USER_LOGIN_MSG":     "https://owa.momo.vn/public/login",                              // đăng nhập lấy token
		"CHECK_USER_BE_MSG":  "https://api.momo.vn/backend/auth-app/public/CHECK_USER_BE_MSG", // get dữ liệu user
		"SEND_OTP_MSG":       "https://api.momo.vn/backend/otp-app/public/SEND_OTP_MSG",       // gửi otp đăng nhập thiết bị mới
		"REG_DEVICE_MSG":     "https://api.momo.vn/backend/otp-app/public/REG_DEVICE_MSG",     // xác thực otp thiết bị mới}
		"DETAIL_TRANS":       "https://api.momo.vn/sync/transhis/details",
	}
	return &momoApi{
		Phone:    phone,
		Password: password,
		Imei:     "2414F000-9986-4414-A29A-E5059307CDA4",
		Config: momoApiCfg{
			ApiAction: apiAction,
			Device:    "IPhone",
			AppVer:    31171,
			AppCode:   "3.1.17",
		},
	}
}

//public function generateCheckSum($msgType, $time, $key)
//{
//$l = $time . '000000';
//$f = $this->phone . $l . $msgType . ($time / 1e12) . "E12";
//return @openssl_encrypt($f, 'AES-256-CBC',  substr($key, 0, 32), 0, '');
//}

func (m *momoApi) generateCheckSum(msgType, time, key string) (string, error) {
	l := time + "000000"
	s, err := strconv.ParseFloat(time, 64)
	if err != nil {
		return "", err
	}
	s = s / 1e12
	f := m.Phone + l + msgType + fmt.Sprintf("%g", s) + "E12"
	res, err := momoencoder.EncodeMomoRq(f, key[:32])
	if err != nil {
		return "", err
	}
	return res, nil
}

func (m *momoApi) getTimeNow() string {
	p0, p1 := microTime()
	return fmt.Sprintf("%d", int64(p0*1000.0)+p1*1000)
}

func (m *momoApi) randomString(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
func (m *momoApi) secureId() string {
	n := 17
	var letterRunes = []rune("0123456789abcde")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (m *momoApi) curlPost(apiUrl string, data map[string]interface{}, header map[string]interface{}, msgType, requestKey, phone, auth string) (map[string]interface{}, error) {
	targetURL, err := url.Parse(apiUrl)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	json_data, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(json_data))
	if err != nil {
		return nil, err
	}
	headers := map[string]interface{}{
		"Accept":          "application/json",
		"msgtype":         msgType,
		"Content-Length":  fmt.Sprintf("%d", len(data)),
		"accept-encoding": "gzip, deflate",
		"Content-Type":    "application/json",
		"Host":            strings.Split(targetURL.Host, ":")[0],
	}
	if auth == "" {
		headers["authorization"] = "Bearer"
	} else {
		headers["authorization"] = "Bearer " + auth
	}
	if requestKey != "" {
		headers["requestkey"] = requestKey
	}
	if phone != "" {
		headers["userid"] = phone
	}
	if header != nil {
		for k, v := range header {
			headers[k] = v
		}
	}

	for k, v := range headers {
		req.Header.Add(k, v.(string))
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	var res map[string]interface{}
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	if err := json.NewDecoder(reader).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (m *momoApi) encodeCurlPost(apiUrl string, _data map[string]interface{}, header map[string]interface{}, msgType, requestKey, requestKeyRaw, phone, auth string) (string, error) {
	rawData, err := json.Marshal(_data)
	if err != nil {
		return "", err
	}
	data, err := momoencoder.EncodeMomoRq(string(rawData), requestKeyRaw)
	if err != nil {
		return "", err
	}
	targetURL, err := url.Parse(apiUrl)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return "", err
	}
	headers := map[string]interface{}{
		"Accept":          "application/json",
		"msgtype":         msgType,
		"Content-Length":  fmt.Sprintf("%d", len(data)),
		"accept-encoding": "gzip, deflate",
		"Content-Type":    "application/json",
		"Host":            strings.Split(targetURL.Host, ":")[0],
	}
	if auth == "" {
		headers["authorization"] = "Bearer"
	} else {
		headers["authorization"] = "Bearer " + auth
	}
	if requestKey != "" {
		headers["requestkey"] = requestKey
	}
	if phone != "" {
		headers["userid"] = phone
	}
	if header != nil {
		for k, v := range header {
			headers[k] = v
		}
	}

	for k, v := range headers {
		req.Header.Add(k, v.(string))
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}
	buf := new(strings.Builder)
	_, err = io.Copy(buf, reader)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
