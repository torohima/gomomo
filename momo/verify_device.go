package momo

import (
	"api_bank/momo/momoencoder"
)

type verifyDeviceResponse struct {
	Error float64
	Msg   string
	Key   string
	PHash string
}

func (m *momoApi) VerifyDevice(code string) (*verifyDeviceResponse, error) {
	oHash := newSHA256(m.Phone + "12345678901234567890" + code)
	action := "REG_DEVICE_MSG"
	time := m.getTimeNow()
	data := map[string]interface{}{
		"user":      m.Phone,
		"msgType":   action,
		"cmdId":     time + "000000",
		"lang":      "vi",
		"channel":   "APP",
		"time":      time,
		"appVer":    m.Config.AppVer,
		"appCode":   m.Config.AppCode,
		"deviceOS":  "Ios",
		"result":    true,
		"errorCode": 0,
		"errorDesc": "",
	}
	extra := map[string]interface{}{
		"ohash":     oHash,
		"AAID":      "",
		"IDFA":      "",
		"TOKEN":     "",
		"SIMULATOR": "false",
		"SECUREID":  m.secureId(),
	}
	momoMsg := map[string]interface{}{
		"_class":      "mservice.backend.entity.msg.RegDeviceMsg",
		"number":      m.Phone,
		"imei":        m.Imei,
		"ccode":       "084",
		"cname":       "Vietnam",
		"device":      m.Config.Device,
		"firmware":    "19",
		"hardware":    "vbox86",
		"manufacture": "samsung",
		"csp":         "",
		"icc":         "",
		"mcc":         "",
		"device_os":   "Ios",
		"secure_id":   m.secureId(),
	}
	data["extra"] = extra
	data["momoMsg"] = momoMsg
	req, err := m.curlPost(m.Config.ApiAction[action], data, nil, action, "", "", "")
	if err != nil {
		return nil, err
	}
	if val, ok := req["result"]; ok {
		if val.(bool) == true {
			extra = req["extra"].(map[string]interface{})
			keySetup := extra["setupKey"].(string)
			key, err := momoencoder.DecodeMomoRq(keySetup, oHash[:32])
			key = key[:32]
			if err != nil {
				return nil, err
			}
			pHash, err := momoencoder.EncodeMomoRq(m.Imei+"|"+m.Password, key)
			if err != nil {
				return nil, err
			}
			return &verifyDeviceResponse{
				Error: 0,
				Msg:   "success",
				Key:   key,
				PHash: pHash,
			}, nil
		}
	}
	return &verifyDeviceResponse{
		Error: req["errorCode"].(float64),
		Msg:   req["errorDesc"].(string),
		Key:   "",
		PHash: "",
	}, nil
}
