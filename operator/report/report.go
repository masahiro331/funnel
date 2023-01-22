package report

import "sync"

var (
	MemStore Store
)

func init() {
	MemStore = Store{reports: make(map[string][]Report)}
}

type Store struct {
	reports map[string][]Report
	mut     sync.Mutex
}

type Report struct {
	Id      string
	Content []byte
}

func (s *Store) Len(OperationId string) int {
	reports, ok := s.reports[OperationId]
	if !ok {
		return 0
	}
	return len(reports)
}

func (s *Store) Pop(OperationId string) []Report {
	s.mut.Lock()
	defer s.mut.Unlock()

	reports, ok := s.reports[OperationId]
	if !ok {
		return []Report{}
	}
	s.reports[OperationId] = []Report{}
	return reports
}

func (s *Store) Push(OperationId string, r Report) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.reports[OperationId] = append(MemStore.reports[OperationId], r)
}
