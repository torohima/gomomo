package momo

type userLoginResponse struct {
	Error     float64
	Msg       string
	Balance   string
	AuthToken string
}

func (m *momoApi) UserLogin(pHash, key string) (*userLoginResponse, error) {
	action := "USER_LOGIN_MSG"
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
		"pass":      m.Password,
	}
	checksum, err := m.generateCheckSum(action, time, key)
	if err != nil {
		return nil, err
	}
	extra := map[string]interface{}{
		"checkSum":  checksum,
		"pHash":     pHash,
		"AAID":      "",
		"IDFA":      "",
		"TOKEN":     "",
		"SIMULATOR": "false",
		"SECUREID":  m.secureId(),
	}
	momoMsg := map[string]interface{}{
		"_class":  "mservice.backend.entity.msg.LoginMsg",
		"isSetup": true,
	}
	data["extra"] = extra
	data["momoMsg"] = momoMsg
	req, err := m.curlPost(m.Config.ApiAction[action], data, nil, action, "", m.Phone, "")
	if err != nil {
		return nil, err
	}
	if val, ok := req["result"]; ok {
		if val.(bool) == true {
			extra := req["extra"].(map[string]interface{})
			m.RequestEncryptKey = extra["REQUEST_ENCRYPT_KEY"].(string)
			m.AuthToken = extra["AUTH_TOKEN"].(string)
			balance := extra["BALANCE"].(string)
			return &userLoginResponse{
				Error:     0,
				Msg:       "success",
				Balance:   balance,
				AuthToken: m.AuthToken,
			}, nil
		}
	}
	return &userLoginResponse{
		Error: req["errorCode"].(float64),
		Msg:   req["errorDesc"].(string),
	}, nil
}
