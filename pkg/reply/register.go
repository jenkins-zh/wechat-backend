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

type ByWeight []AutoReply

func (b ByWeight) Len() int {
	return len(b)
}

func (b ByWeight) Less(i, j int) bool {
	return b[i].Weight() < b[j].Weight()
}

func (b ByWeight) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// func init() {
// 	autoReplyInitChains = make([]Init, 3)
// }
