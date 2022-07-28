package momo

func (m *momoApi) RegNewDevice() (*simpleResponse, error) {
	action := "SEND_OTP_MSG"
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
		"action":    "SEND",
		"rkey":      "12345678901234567890",
		"isVoice":   false,
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
		"cname":       "Vietnam",
		"ccode":       "084",
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
	req, err := m.curlPost(m.Config.ApiAction[action], data, nil, action, "", m.Phone, "")
	if err != nil {
		return nil, err
	}
	if val, ok := req["result"]; ok {
		if val.(bool) == true {
			return &simpleResponse{
				Error: 0,
				Msg:   "success",
			}, nil
		}
	}
	return &simpleResponse{
		Error: req["errorCode"].(float64),
		Msg:   req["errorDesc"].(string),
	}, nil
}
