package hashcash

type Service interface {
	AddHashIndicator(indicator uint64)
	IndicatorExists(indicator uint64) bool
	DeleteIndicator(indicator uint64) error
}

type PowService struct {
	repo Repository
}

func NewService(hashRepo Repository) *PowService {
	return &PowService{
		repo: hashRepo,
	}
}

func (s *PowService) AddHashIndicator(indicator uint64) {
	s.repo.AddIndicator(indicator)
}

func (s *PowService) IndicatorExists(indicator uint64) bool {
	return s.repo.IndicatorExists(indicator)
}

func (s *PowService) DeleteIndicator(indicator uint64) error {
	return s.repo.DeleteIndicator(indicator)
}
