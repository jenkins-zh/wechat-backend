package reply

// MatchAutoReply only reply for match
type WelcomeReply struct {
	AutoReply
}

func InitWelcomeReply() AutoReply {
	return &WelcomeReply{}
}

var _ AutoReply = &WelcomeReply{}

// Accept consider if it will accept the request
func (m *WelcomeReply) Accept(request *TextRequestBody) (ok bool) {

	if "event" == request.MsgType && "subscribe" == request.Event {
		request.Content = "welcome"
		m.AutoReply = InitMatchAutoReply()
		ok = m.AutoReply.Accept(request)
	}
	return
}

// Handle hanlde the request then return data
// func (m *WelcomeReply) Handle() (data []byte, err error) {
// 	return m.
// }
