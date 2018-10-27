package officialaccount

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/clbanning/mxj"
)

type Handler struct {
	Name    string
	Handler func(*Context) *ReplyMessage
}

type server struct {
	config   *Config
	handlers []Handler
}

type ReceiveMessage map[string]interface{}

type ReplyMessage struct {
	Reply bool
	Exit  bool
	Data  map[string]interface{}
}

type EncryptedXMLMessage struct {
	XMLName    struct{} `xml:"xml"`
	ToUserName string   `xml:"ToUserName"`
	Encrypt    string   `xml:"Encrypt"`
}

func (oa *OfficialAccount) Server() *server {
	if oa.server != nil {
		return oa.server
	}

	oa.server = &server{
		config: oa.config,
	}
	return oa.server
}

func (s *server) RegisterHandler(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

func (s *server) Handler(request *http.Request, response http.ResponseWriter) {
	timestamp, _ := s.query(request, "timestamp")
	nonce, _ := s.query(request, "nonce")

	var rawXMLMsgBytes []byte

	if encryptType, ok := s.query(request, "encrypt_type"); ok && encryptType == "aes" {
		// 安全模式
		signature, _ := s.query(request, "msg_signature")
		var encryptedXMLMessage EncryptedXMLMessage
		if err := xml.NewDecoder(request.Body).Decode(&encryptedXMLMessage); err == nil {
			if pass := s.validate(signature, timestamp, nonce, s.config.Token, encryptedXMLMessage.Encrypt); pass {
				_, rawXMLMsgBytes, _ = DecryptMsg(s.config.AppID, encryptedXMLMessage.Encrypt, s.config.EncodingAESKey)
			}
		} else {
		}
	} else {
		// 普通模式
		signature, _ := s.query(request, "signature")
		if pass := s.validate(signature, timestamp, nonce, s.config.Token); pass {

			// 处理验证服务器地址
			if echostr, ok := s.query(request, "echostr"); strings.ToUpper(request.Method) == "GET" && ok {
				response.WriteHeader(http.StatusOK)
				response.Write([]byte(echostr))
				return
			}

			if strings.ToUpper(request.Method) == "POST" {
				rawXMLMsgBytes, _ = ioutil.ReadAll(request.Body)
			}
		}
	}

	if rawXMLMsgBytes != nil {
		if message, err := s.makeMessage(rawXMLMsgBytes); err == nil {
			ret := s.callHandlers(&Context{
				Message: message,
			})
			if ret.Reply {
				s.replyXML(response, ret.Data)
				return
			}
		}
	}

	s.replyNone(response)
}

func (s *server) validate(signature string, params ...string) bool {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	signatureGen := fmt.Sprintf("%x", h.Sum(nil))
	return signature == signatureGen
}

func (s *server) makeMessage(body []byte) (ReceiveMessage, error) {
	mv, err := mxj.NewMapXml(body)

	if err != nil {
		return nil, err
	}

	msg, _ := mv["xml"].(map[string]interface{})

	return msg, nil
}

func (s *server) callHandlers(context *Context) *ReplyMessage {
	var reply *ReplyMessage

	for _, h := range s.handlers {
		reply = h.Handler(context)

		if reply != nil && (reply.Reply || reply.Exit) {
			break
		}
	}

	if reply != nil {
		return reply
	}

	return &ReplyMessage{
		Reply: false,
	}
}

func (s *server) replyXML(response http.ResponseWriter, data map[string]interface{}) {
	mv := mxj.Map(data)
	rawXML, err := mv.Xml("xml")

	if err != nil {
		s.replyNone(response)
	} else {
		s.reply(response, rawXML)
	}
}

func (s *server) replyNone(response http.ResponseWriter) {
	s.reply(response, []byte("success"))
}

func (s *server) reply(response http.ResponseWriter, data []byte) {
	response.WriteHeader(http.StatusOK)
	response.Write(data)
}

func (s *server) query(request *http.Request, key string) (string, bool) {
	if values, ok := request.URL.Query()[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}
