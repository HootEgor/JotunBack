package repository

import (
	"JotunBack/model"
	"cloud.google.com/go/firestore"
	"context"
)

type UserRepository struct {
	context context.Context
	c       *firestore.Client
}

func NewUserRepository(c *firestore.Client, ctx context.Context) *UserRepository {
	return &UserRepository{ctx, c}
}

func (r *UserRepository) CreateUser(user model.User) error {
	_, err := r.c.Collection("users").Doc(user.Username).Set(r.context, user)
	return err
}

func (r *UserRepository) GetUser(username string) (model.User, error) {
	doc, err := r.c.Collection("users").Doc(username).Get(r.context)
	if err != nil {
		return model.User{}, err
	}
	var user model.User
	doc.DataTo(&user)
	return user, nil
}

func (r *UserRepository) UpdateUser(user model.User) error {
	_, err := r.c.Collection("users").Doc(user.Username).Set(r.context, user)
	return err
}

func (r *UserRepository) DeleteUser(username string) error {
	_, err := r.c.Collection("users").Doc(username).Delete(r.context)
	return err
}

func (r *UserRepository) CreateACState(acConfig model.AirConditionerConfig) error {
	_, err := r.c.Collection("users").Doc(acConfig.Username).Collection("acState").Doc("state").Set(r.context, acConfig)
	return err
}

func (r *UserRepository) GetACState(username string) (model.AirConditionerConfig, error) {
	doc, err := r.c.Collection("users").Doc(username).Collection("acState").Doc("state").Get(r.context)
	if err != nil {
		return model.AirConditionerConfig{}, err
	}
	var acState model.AirConditionerConfig
	doc.DataTo(&acState)
	return acState, nil
}

func (r *UserRepository) UpdateACState(acConfig model.AirConditionerConfig) error {
	_, err := r.c.Collection("users").Doc(acConfig.Username).Collection("acState").Doc("state").Set(r.context, acConfig)
	return err
}

func (r *UserRepository) DeleteACState(username string) error {
	_, err := r.c.Collection("users").Doc(username).Collection("acState").Doc("state").Delete(r.context)
	return err
}

func (r *UserRepository) CreateTemp(temp model.TempDB, userName string) error {
	_, err := r.c.Collection("users").Doc(userName).Collection("temp").Doc(temp.TimeStamp.Format("2006-01-02 15:04:05")).Set(r.context, temp)
	return err
}

func (r *UserRepository) GetTemp(userName string) ([]model.TempDB, error) {
	iter := r.c.Collection("users").Doc(userName).Collection("temp").OrderBy("timeStamp", firestore.Desc).Limit(3).Documents(r.context)
	var temps []model.TempDB
	for {
		doc, err := iter.Next()
		if err == nil {
			var temp model.TempDB
			doc.DataTo(&temp)
			temps = append(temps, temp)
		} else {
			break
		}
	}
	return temps, nil
}

func (r *UserRepository) DeleteTemp(userName string) error {
	iter := r.c.Collection("users").Doc(userName).Collection("temp").Documents(r.context)
	for {
		doc, err := iter.Next()
		if err == nil {
			_, err = doc.Ref.Delete(r.context)
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}
