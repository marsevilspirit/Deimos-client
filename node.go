package deimosclient

type Node struct {
	Key           string `json:"key"`
	Value         string `json:"value,omitempty"`
	Dir           bool   `json:"dir,omitempty"`
	ModifiedIndex uint64 `json:"modifiedIndex"`
	CreatedIndex  uint64 `json:"createdIndex"`
}
