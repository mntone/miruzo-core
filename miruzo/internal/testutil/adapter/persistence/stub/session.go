package stub

import "github.com/mntone/miruzo-core/miruzo/internal/persist"

type session struct {
	ActionStub *actionRepository
	JobStub    *jobRepository
	StatsStub  *statsRepository
	UserStub   *userRepository
}

func (s *session) ImageList() persist.ImageListRepository {
	return nil
}

func (s *session) Job() persist.JobRepository {
	return s.JobStub
}

func (s *session) Settings() persist.SettingsRepository {
	return nil
}

func (s *session) User() persist.UserRepository {
	return s.UserStub
}

type txSession struct {
	*session
}

func (s *txSession) Action() persist.ActionRepository {
	return s.ActionStub
}

func (s *txSession) Stats() persist.StatsRepository {
	return s.StatsStub
}

func (s *txSession) User() persist.SessionUserRepository {
	return s.UserStub
}

func (s *txSession) View() persist.ViewRepository {
	return nil
}
