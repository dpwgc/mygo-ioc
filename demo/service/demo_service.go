package service

import "mygo-ioc/demo/repository"

type DemoService struct {
	DemoRepository *repository.DemoRepository `autowired:"repository"`
}

func (s *DemoService) Print(text string) {
	s.DemoRepository.Print(text)
}
