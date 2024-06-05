package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/farid21ola/forum/domain"
	"github.com/farid21ola/forum/model"
	"sync"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Domain      *domain.Domain
	NewComments []*model.Comment
	// All active subscriptions
	CommentsObservers map[string][]chan *model.Comment
	mu                sync.Mutex
}
