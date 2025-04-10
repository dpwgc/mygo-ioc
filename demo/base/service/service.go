package service

import "mygo-ioc/demo/base/repository"

type DemoService struct {
	FirstRepository  repository.Repository `autowired:"repository" qualifier:"FirstRepository"`
	SecondRepository repository.Repository `autowired:"repository" qualifier:"SecondRepository"`
}

func (s *DemoService) FirstInsert(text string) bool {
	return s.FirstRepository.Insert(text)
}

func (s *DemoService) SecondInsert(text string) bool {
	return s.SecondRepository.Insert(text)
}
