package pay

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
const urlPreOrder = "https://api.mch.weixin.qq.com/pay/unifiedorder"

func (m *MerchantInfo) PreOrder(req PreOrderRequest) (*PreOrderResponse, error) {
	// check parameters
	if req.FeeType <= 0 {
		req.FeeType = CNY
	}

	resp := new(PreOrderResponse)
	if e := m.postXML(urlPreOrder, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_2
const urlOrderQuery = "https://api.mch.weixin.qq.com/pay/orderquery"

func (m *MerchantInfo) QueryOrder(req PreOrderRequest) (*PreOrderResponse, error) {
	// check parameters
	if req.FeeType <= 0 {
		req.FeeType = CNY
	}

	resp := new(PreOrderResponse)
	if e := m.postXML(urlPreOrder, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}