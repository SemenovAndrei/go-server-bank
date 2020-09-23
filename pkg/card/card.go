package card

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrType      = errors.New("wrong card type")
	ErrSystem    = errors.New("wrong card system")
	ErrUserId    = errors.New("user ID not found")
	ErrBasicCard = errors.New("no basic card")
)

type Card struct {
	Id     int64
	UserId string
	Number int64
	Type   string
	System string
}

type Service struct {
	mu     sync.RWMutex
	cards  []*Card
	mun    sync.Mutex
	number int64
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AllCards() []*Card {
	s.mu.RLock()
	defer s.mu.RLock()
	return s.cards
}

func (s *Service) Add(userId, cardType, cardSystem string) (*Card, error) {

	err := getTypeCard(cardType)
	if err != nil {
		return &Card{}, err
	}

	err = getSystemCard(cardSystem)
	if err != nil {
		return &Card{}, err
	}

	err = s.isBasicCard(userId)
	if err != nil && cardType != "basic" {
		return &Card{}, err
	}

	number := s.getNumber()

	card := &Card{
		Id:     number,
		UserId: userId,
		Number: number,
		Type:   cardType,
		System: cardSystem,
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cards = append(s.cards, card)
	return card, nil
}

func getTypeCard(typeCard string) error {
	typesCard := []string{"basic", "additional", "virtual"}
	for _, value := range typesCard {
		if value == typeCard {
			return nil
		}
	}
	return ErrType
}

func getSystemCard(systemCard string) error {
	systemsCard := []string{"Visa", "MasterCard"}
	for _, value := range systemsCard {
		if strings.ToLower(value) == strings.ToLower(systemCard) {
			return nil
		}
	}
	return ErrSystem
}

func (s *Service) getNumber() int64 {
	s.mun.Lock()
	defer s.mun.Unlock()
	s.number += 1
	return s.number
}

func (s *Service) isBasicCard(userId string) error {
	for _, value := range s.cards {
		if value.UserId == userId {
			return nil
		}
	}
	return ErrBasicCard
}
