package mail_sender

type sender struct {
}

func NewSender() *sender {
	return &sender{}
}

func (s *sender) Send(_ string, _ int) error {
	// choo choo i'm a sender
	// i'm sending requests to a smtp-server

	return nil
}
