package pay

import (
	"encoding/xml"
	"log"
	"net/http"
	"reflect"
)

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_7
// 注意：同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知。
// 推荐的做法是，当收到通知进行处理时，首先检查对应业务数据的状态，判断该通知是否已经处理过，
// 如果没有处理过再进行处理，如果处理过直接返回结果成功。
// 在对业务数据进行状态检查和处理之前，要采用数据锁进行并发控制，以避免函数重入造成的数据混乱。
// 特别提醒：商户系统对于支付结果通知的内容一定要做签名验证,并校验返回的订单金额是否与商户侧的订单金额一致，
// 防止数据泄漏导致出现“假通知”，造成资金损失。
type PayNoticeHander interface {
	HandlePayNotice(*PayNotice) error
}

func (m *MerchantInfo) PayNotice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fields, e := parseToFields(r.Body)
	if e != nil {
		noticeFailed(w)
		return
	}
	if verbose {
		log.Println("got pay notice: ", fields)
		defer func() {
			if e != nil {
				log.Println("handle pay notice error: ", e)
			}
		}()
	}

	if e = checkReturnCode(fields); e != nil {
		noticeFailed(w)
		return
	}
	if e = m.checkAppID(fields); e != nil {
		noticeFailed(w)
		return
	}
	if e = m.checkSign(fields); e != nil {
		noticeFailed(w)
		return
	}

	notice := new(PayNotice)
	if e = checkResultCode(fields); e != nil {
		notice.TradeState = PayError
	} else {
		notice.TradeState = PaySuccess
	}

	if e = composeStruct(fields, reflect.ValueOf(notice)); e != nil {
		noticeFailed(w)
		return
	}
	if verbose {
		log.Printf("PayNotice: %+v", notice)
	}

	if e = m.HandlePayNotice(notice); e != nil {
		noticeFailed(w)
		return
	}
	noticeSuccess(w)
}

// https://pay.weixin.qq.com/wiki/doc/api/jsapi.php?chapter=9_16&index=9
// 注意：同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知。
// 推荐的做法是，当收到通知进行处理时，首先检查对应业务数据的状态，判断该通知是否已经处理过，
// 如果没有处理过再进行处理，如果处理过直接返回结果成功。
// 在对业务数据进行状态检查和处理之前，要采用数据锁进行并发控制，以避免函数重入造成的数据混乱。
type RefundNoticeHandler interface {
	HandleRefundNotice(*RefundNotice) error
}

func (m *MerchantInfo) RefundNotice(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fields, e := parseToFields(r.Body)
	if e != nil {
		noticeFailed(w)
		return
	}
	if verbose {
		log.Println("got refund notice: ", fields)
		defer func() {
			if e != nil {
				log.Println("handle refund notice error: ", e)
			}
		}()
	}

	if e = checkReturnCode(fields); e != nil {
		noticeFailed(w)
		return
	}
	if e = m.checkAppID(fields); e != nil {
		noticeFailed(w)
		return
	}

	data, e := decodeNoticeMessage(fields["req_info"], m.PaymentKey)
	if e != nil {
		noticeFailed(w)
		return
	}

	notice := new(RefundNotice)
	if e = xml.Unmarshal(data, notice); e != nil {
		noticeFailed(w)
		return
	}

	if e = m.HandleRefundNotice(notice); e != nil {
		noticeFailed(w)
		return
	}
	noticeSuccess(w)
}

func noticeFailed(w http.ResponseWriter) {
	r := NoticeResult{
		ReturnCode: CData{"FAIL"},
	}
	data, e := xml.Marshal(r)
	if e != nil {
		return
	}
	if verbose {
		log.Println("handle notice failed: ", string(data))
	}
	w.Write(data)
}

func noticeSuccess(w http.ResponseWriter) {
	r := NoticeResult{
		ReturnCode:    CData{"SUCCESS"},
		ReturnMessage: CData{"OK"},
	}
	data, e := xml.Marshal(r)
	if e != nil {
		return
	}
	if verbose {
		log.Println("handle notice success: ", string(data))
	}
	w.Write(data)
}
