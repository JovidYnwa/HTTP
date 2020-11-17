package banners

import (
	"context"
	"errors"
	"os"
	"sync"
)

type Service struct {
	nextBannerID int64
	mu           sync.RWMutex
	items        []*Banner
}

type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//All возращает все существующие баннеры
func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	return s.items, nil
}

//ByID возращает баннер по идентификатору
func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}

//Save сохраняет/обновляет баннер
func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, banner := range s.items {
		if banner.ID == item.ID {
			banner.Title = item.Title
			banner.Content = item.Content
			banner.Link = item.Link
			banner.Button = item.Button
			return banner, nil
		}
	}

	s.nextBannerID++
	item.ID = s.nextBannerID
	s.items = append(s.items, item)

	return s.items[len(s.items)-1], nil
}

//RemoveByID удаляет баннер по идентификатору
func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()
	for index, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:index], s.items[index+1:]...)
			pathToImage := "web/banners/" + banner.Image
			err := os.Remove(pathToImage)
			if err != nil {
				return nil, err
			}
			return banner, nil
		}
	}

	return nil, errors.New("item not found")
}
