package logs

type LogService struct {
	repo ILogRepository
}

func NewLogService(repo ILogRepository) ILogService {
	return &LogService{
		repo: repo,
	}
}

type ILogService interface {
	InsertLog(l Log) (string, error)
}

// ----------------------------

func (s *LogService) InsertLog(l Log) (string, error) {
	return s.repo.InsertLog(l)
}
