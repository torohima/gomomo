package momo

import (
	"api_bank/momo/momoencoder"
	"encoding/json"
)

type detailTxData struct {
	CreatedAt   int
	Io          int
	ServiceId   string
	SourceName  string
	TargetName  string
	TotalAmount int
	TransId     int64
	Comment     string
}

type detailTxResponse struct {
	Error float64
	Msg   string
	Data  *detailTxData
}

func (m *momoApi) DetailTrans(txId int64) (*detailTxResponse, error) {
	if m.AuthToken == "" {
		return &detailTxResponse{
			Error: 1,
			Msg:   "Login failed",
		}, nil
	}
	action := "DETAIL_TRANS"
	data := map[string]interface{}{
		"requestId":   m.getTimeNow(),
		"transId":     txId,
		"serviceId":   "transfer_p2p",
		"appCode":     m.Config.AppCode,
		"appVer":      m.Config.AppVer,
		"lang":        "vi",
		"deviceOS":    "IOS",
		"channel":     "APP",
		"buildNumber": 0,
		"appId":       "vn.momo.transactionhistory",
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
		iData := v.(map[string]interface{})
		var data detailTxData
		bData, err := json.Marshal(iData)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(bData, &data); err != nil {
			return nil, err
		}
		var cmtData map[string]interface{}
		if err := json.Unmarshal([]byte(iData["serviceData"].(string)), &cmtData); err != nil {
			return nil, err
		}
		data.Comment = cmtData["COMMENT_VALUE"].(string)
		return &detailTxResponse{
			Error: 0,
			Msg:   "success",
			Data:  &data,
		}, nil
	} else {
		return &detailTxResponse{
			Error: 1,
			Msg:   "Cannot get history",
			Data:  nil,
		}, nil
	}
}
