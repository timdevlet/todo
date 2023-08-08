package cards

type CardsService struct {
	repo *CardsRepository
}

func NewCardService(repo *CardsRepository) *CardsService {
	return &CardsService{
		repo: repo,
	}
}
