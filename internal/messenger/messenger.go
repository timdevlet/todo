package messenger

type MessengerService struct {
	repo *MessageRepository
}

func NewMessengerService(repo *MessageRepository) *MessengerService {
	return &MessengerService{
		repo: repo,
	}
}

// ----------------------------

func (ts *MessengerService) SendMessage(from_uuid, to_uuid string, text ...string) (string, error) {
	return ts.repo.InsertMessage(from_uuid, to_uuid, text...)
}

func (ts *MessengerService) GetChannels(user_uuid string) ([]ChannelDirect, error) {
	return ts.repo.FetchChannels(user_uuid)
}

func (ts *MessengerService) CountChannelMessages(channel_uuid string) (int, error) {
	return ts.repo.CountChannelMessages(channel_uuid)
}

func (ts *MessengerService) GetMessages(channel_uuid string) ([]Message, error) {
	return ts.repo.FetchMessages(channel_uuid)
}

func (ts *MessengerService) FetchMessagesByUser(user_uuid string) ([]Message, error) {
	return ts.repo.FetchMessagesByUser(user_uuid)
}
