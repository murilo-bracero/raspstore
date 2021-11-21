package repository

import (
	"context"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/uuid"
	"raspstore.github.io/authentication/db"
	"raspstore.github.io/authentication/model"
)

const UsersKindName = "users"

type usersDatastoreRespository struct {
	ctx    context.Context
	client *datastore.Client
}

func NewDatastoreUsersRepository(ctx context.Context, conn db.DatastoreConnection) UsersRepository {
	return &usersDatastoreRespository{client: conn.Client(), ctx: ctx}
}

func (r *usersDatastoreRespository) Save(user *model.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.UserId = uuid.NewString()

	key := datastore.NameKey(UsersKindName, user.UserId, nil)

	if _, err := r.client.Put(r.ctx, key, user); err != nil {
		return err
	}

	return nil
}

func (r *usersDatastoreRespository) FindById(id string) (user *model.User, err error) {

	found := new(model.User)
	key := datastore.NameKey(UsersKindName, id, nil)

	if err := r.client.Get(r.ctx, key, found); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

	return found, nil
}

func (r *usersDatastoreRespository) FindByEmailOrUsername(email string, username string) (*model.User, error) {

	var users []*model.User

	query := datastore.NewQuery(UsersKindName).Filter("email =", email).Filter("username =", username)

	if _, err := r.client.GetAll(r.ctx, query, &users); err != nil {
		if err == datastore.ErrNoSuchEntity {
			return nil, nil
		}
		return nil, err
	}

	return users[len(users)-1], nil
}

func (r *usersDatastoreRespository) DeleteUser(id string) error {
	key := datastore.NameKey(UsersKindName, id, nil)

	return r.client.Delete(r.ctx, key)
}

func (r *usersDatastoreRespository) UpdateUser(user *model.User) error {
	key := datastore.NameKey(UsersKindName, user.UserId, nil)

	found, err := r.FindById(user.UserId)

	if err != nil {
		return err
	}

	if found == nil {
		return datastore.ErrNoSuchEntity
	}

	found.UpdatedAt = time.Now()
	found.Email = user.Email
	found.Username = user.Username
	found.PhoneNumber = user.PhoneNumber

	if _, err := r.client.Put(r.ctx, key, found); err != nil {
		return err
	}

	return nil
}

func (r *usersDatastoreRespository) FindAll() (users []*model.User, err error) {

	query := datastore.NewQuery(UsersKindName)
	if _, err := r.client.GetAll(r.ctx, query, &users); err != nil {
		return nil, err
	}

	return users, nil
}
