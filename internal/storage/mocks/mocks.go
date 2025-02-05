package mocks

import (
	"context"
	"github.com/Denis-Kuso/rss_collector/server/internal/storage"
	"github.com/google/uuid"
)

// TODO are these two structs really neccessary?
type UsersModel struct {
}

type FeedsModel struct {
}

//func (m *UsersModel) Create(ctx context.Context, username string) (storage.User, error) {
// could create some scenarios here depending on what the db might return
//	return storage.User{}, nil
//}

func (m *UsersModel) Get() (storage.User, error)
func (f *FeedsModel) Create(ctx context.Context, name, URL string) error
func (f *FeedsModel) Get(ctx context.Context, userID uuid.UUID) error
func (f *FeedsModel) Follow(ctx context.Context, userID, feedID uuid.UUID) error // good name?

func (f *FeedsModel) ShowAvailable(ctx context.Context) error
func (f *FeedsModel) Delete(ctx context.Context) error
func (f *FeedsModel) GetPosts(ctx context.Context, userID uuid.UUID, limit int) error // should this be in its own file
