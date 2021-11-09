package nlp

import (
	"fmt"
	"log"
	"strconv"
)

/*
	The nlp (Natural Language Processing) service.
*/

type Service struct {
	defaultLength string
	bucketName    string
}

func NewService() *Service {
	return &Service{
		defaultLength: "120",
	}
}

func (s *Service) Learn(data string) error {
	if data == "" {
		return nil
	}

	return nil
}

func (s *Service) Generate() (*string, error) {
	lengthVal, err := strconv.Atoi(s.defaultLength)
	if err != nil {
		return nil, fmt.Errorf("failed to convert length to int: %w", err)
	}
	log.Printf("length: %v", lengthVal)

	example := "generate\n"

	return &example, nil
}
