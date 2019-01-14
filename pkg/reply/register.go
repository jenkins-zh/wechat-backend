package reply

var autoReplyInitChains []Init

// Register add an implement of AutoReply
func Register(initFunc Init) {
	autoReplyInitChains = append(autoReplyInitChains, initFunc)
}

// AutoReplyChains return all implements of AutoReply
func AutoReplyChains() []Init {
	return autoReplyInitChains
}

// func init() {
// 	autoReplyInitChains = make([]Init, 3)
// }
