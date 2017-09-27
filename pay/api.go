package pay

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_1
const urlPreOrder = "https://api.mch.weixin.qq.com/pay/unifiedorder"

func (m *MerchantInfo) PreOrder(req PreOrderRequest) (*PreOrderResponse, error) {
	// check parameters
	if req.FeeType == "" {
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

func (m *MerchantInfo) QueryOrder(req QueryOrderRequest) (*QueryOrderResponse, error) {
	// check parameters

	resp := new(QueryOrderResponse)
	if e := m.postXML(urlPreOrder, req, resp, false); e != nil {
		return nil, e
	}

	return resp, nil
}

// need certification
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_3
const urlCloseOrder = "https://api.mch.weixin.qq.com/pay/closeorder"

func (m *MerchantInfo) CloseOrder(req CloseOrderRequest) error {
	// check parameters
	if e := m.postXML(urlPreOrder, req, nil, true); e != nil {
		return e
	}
	return nil
}

// need certification
// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_4
const urlRefundOrder = "https://api.mch.weixin.qq.com/secapi/pay/refund"

func (m *MerchantInfo) RefundOrder(req RefundRequest) (*RefundResponse, error) {
	// check parameters
	resp := new(RefundResponse)
	if e := m.postXML(urlPreOrder, req, resp, true); e != nil {
		return nil, e
	}
	return resp, nil
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_5
const urlQueryRefund = "https://api.mch.weixin.qq.com/pay/refundquery"
