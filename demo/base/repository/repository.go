package repository

import "fmt"

type Repository interface {
	Insert(text string) bool
}

type FirstRepository struct {
}

func (fr *FirstRepository) Insert(text string) bool {
	fmt.Println("first insert:", text)
	return true
}

type SecondRepository struct {
}

func (sr *SecondRepository) Insert(text string) bool {
	fmt.Println("second insert:", text)
	return true
}
