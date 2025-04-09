package repository

import "mygo-ioc/demo/mapper"

type DemoRepository struct {
	DemoMapper *mapper.DemoMapper `autowired:"mapper"`
}

func (r *DemoRepository) Print(text string) {
	r.DemoMapper.Print(text)
}
