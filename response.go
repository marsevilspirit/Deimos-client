package deimosclient

type Response struct {
	Action    string `json:"action"`
	Node      *Node  `json:"node"`
	PrevNode  *Node  `json:"prevNode,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
}
