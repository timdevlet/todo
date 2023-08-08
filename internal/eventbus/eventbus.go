package eventbus

type Event struct {
	Channel string

	Body string
}

type ProfileEvent struct {
	Event
	Name string
}

func NewUpdateProfileEvent(newName string) *ProfileEvent {
	return &ProfileEvent{
		Event: Event{
			Channel: "profile.update.fio",
		},
		Name: newName,
	}
}
