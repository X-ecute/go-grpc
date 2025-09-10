//go:generate mockgen -destination=rocket_mocks_test.go -package=rocket github.com/X-ecute/go-grpc/internal/rocket Store
package rocket

import "context"

// Store - the interface we expect
// our db implementation to follow
type Store interface {
	GetRocketById(id int) (Rocket, error)
	InsertRocket(rocket Rocket) (Rocket, error)
	UpdateRocket(rocket Rocket) (Rocket, error)
	DeleteRocket(id string) error
}

// Rocket - Definition of our rocket
type Rocket struct {
	ID      string
	Name    string
	Type    string
	Flights int
}

// Service  - responsible gor updating the rocker inventory
type Service struct {
	Store Store
}

// NewService - returns a new instance of our rocket service
func NewService(store Store) Service {
	return Service{
		Store: store,
	}
}

// GetRocketById - retrieves a rocket based on ID
func (s Service) GetRocketById(ctx context.Context, id int) (Rocket, error) {
	rkt, err := s.Store.GetRocketById(id)
	if err != nil {
		return Rocket{}, err
	}
	return rkt, nil
}

// InsertRocket - inserts a rocket into the store
func (s Service) InsertRocket(ctx context.Context, rocket Rocket) (Rocket, error) {
	return s.Store.InsertRocket(rocket)
}

// UpdateRocket - updates a rocket in store
func (s Service) UpdateRocket(ctx context.Context, rocket Rocket) (Rocket, error) {
	return s.Store.UpdateRocket(rocket)
}

// DeleteRocket - deletes a rocket from the store
func (s Service) DeleteRocket(ctx context.Context, id string) error {
	return s.Store.DeleteRocket(id)
}
