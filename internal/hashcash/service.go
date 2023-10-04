package hashcash

type Service interface {
	AddHashIndicator(indicator uint64) error
	GetHashIndicator(indicator uint64) (uint64, error)
	RemoveHashIndicator(indicator uint64) error
}

type PowService struct {
	repo Repository
}

func NewService(hashRepo Repository) *PowService {
	return &PowService{
		repo: hashRepo,
	}
}

func (s *PowService) AddHashIndicator(indicator uint64) error {
	return s.repo.AddIndicator(indicator)
}

func (s *PowService) GetHashIndicator(indicator uint64) (uint64, error) {
	return s.repo.GetIndicator(indicator)
}

func (s *PowService) RemoveHashIndicator(indicator uint64) error {
	return s.repo.RemoveIndicator(indicator)
}
