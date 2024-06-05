package domain

import (
	"context"
	"errors"
	"github.com/farid21ola/forum/model"
	"log"
)

func (d *Domain) Login(ctx context.Context, input *model.LoginInput) (*model.AuthResponse, error) {
	user, err := d.Storage.UserByUsername(ctx, input.Username)
	if err != nil {
		return nil, ErrBadCredentials
	}

	err = user.ComparePassword(input.Password)
	if err != nil {
		return nil, ErrBadCredentials
	}

	token, err := user.GenToken()
	if err != nil {
		return nil, errors.New("something went wrong")
	}

	return &model.AuthResponse{
		AuthToken: token,
		User:      user,
	}, nil
}

func (d *Domain) Register(ctx context.Context, input *model.RegisterInput) (*model.AuthResponse, error) {
	_, err := d.Storage.UserByUsername(ctx, input.Username)
	if err == nil {
		return nil, errors.New("username is already in use")
	}

	user := &model.User{
		Username:  input.Username,
		Password:  input.Password,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err = user.HashPassword(input.Password)
	if err != nil {
		log.Printf("can't hash password: %v", err)
		return nil, errors.New("something went wrong")
	}

	tx, err := d.Storage.Begin(ctx)
	if err != nil {
		log.Printf("error creating a transaction: %v", err)
		return nil, errors.New("something went wrong")
	}
	defer tx.Rollback(ctx)

	if _, err = d.Storage.CreateUser(ctx, tx, user); err != nil {
		log.Printf("error creating a user: %v", err)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("error while commiting tx: %v", err)
		return nil, err
	}

	token, err := user.GenToken()
	if err != nil {
		log.Printf("error while generating the token: %v", err)
		return nil, errors.New("something went wrong")
	}

	return &model.AuthResponse{
		AuthToken: token,
		User:      user,
	}, nil
}
