package momo

import (
	"api_bank/momo/momoencoder"
	"encoding/json"
)

type getHistoryData struct {
	CreatedAt   int
	Io          int
	ServiceId   string
	SourceName  string
	TargetName  string
	TotalAmount int
	TransId     int64
}

type getHistoryResponse struct {
	Error float64
	Msg   string
	Data  []getHistoryData
}

func (m *momoApi) GetHistory(startDate, endDate string, limit int) (*getHistoryResponse, error) {
	if m.AuthToken == "" {
		return &getHistoryResponse{
			Error: 1,
			Msg:   "Login failed",
		}, nil
	}
	action := "QUERY_TRAN_HIS_MSG"
	time := m.getTimeNow()
	data := map[string]interface{}{
		"requestId":   time,
		"startDate":   startDate,
		"endDate":     endDate,
		"offset":      0,
		"limit":       limit,
		"appCode":     m.Config.AppCode,
		"appVer":      m.Config.AppVer,
		"lang":        "vi",
		"deviceOS":    "IOS",
		"channel":     "APP",
		"buildNumber": 0,
		"appId":       "vn.momo.platform",
	}
	requestKeyRaw := m.randomString(32)
	requestKey, err := momoencoder.EncodeRSAMomoPubKey(requestKeyRaw, m.RequestEncryptKey)
	if err != nil {
		return nil, err
	}

	rawRq, err := m.encodeCurlPost(m.Config.ApiAction[action], data, nil, action, requestKey, requestKeyRaw, m.Phone, m.AuthToken)
	if err != nil {
		return nil, err
	}
	resp, err := momoencoder.DecodeMomoRq(rawRq, requestKeyRaw)
	if err != nil {
		return nil, err
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
		return nil, err
	}
	if v, ok := obj["momoMsg"]; ok {
		iData := v.([]interface{})
		var data []getHistoryData
		for i := range iData {
			bData, err := json.Marshal(iData[i].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			child := getHistoryData{}
			if err := json.Unmarshal(bData, &child); err != nil {
				return nil, err
			}
			data = append(data, child)
		}
		return &getHistoryResponse{
			Error: 0,
			Msg:   "success",
			Data:  data,
		}, nil
	} else {
		return &getHistoryResponse{
			Error: 1,
			Msg:   "Cannot get history",
			Data:  nil,
		}, nil
	}
}
