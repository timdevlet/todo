package messenger

type Message struct {
	UUID        string `validate:"uuid4"`
	OwnerUuid   string `validate:"uuid4"`
	ChannelUuid string `validate:"uuid4"`
	Text        string
}

type ChannelDirect struct {
	UUID      string `validate:"uuid4"`
	FromUuid  string `validate:"uuid4"`
	ToUuid    string `validate:"uuid4"`
	CreatedAt string
}

type Total struct {
	Total int
}

type InsertMessageInput struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}
