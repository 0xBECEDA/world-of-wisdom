package quotes

type QuoteService interface {
	GetQuote() Quote
}

type Service struct {
	repo QuoteRepo
}

func NewService(repository QuoteRepo) *Service {
	return &Service{
		repo: repository,
	}
}

func (r *Service) GetQuote() Quote {
	return r.repo.GetQuote()
}
