package userbiz

import (
	"app-invite-service/common"
	"app-invite-service/module/user/usermodel"
	"context"
)

type RegisterStore interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*usermodel.User, error)
	CreateUser(ctx context.Context, data *usermodel.UserCreate) error
}

type Hash interface {
	Hash(data string) string
}

type registerBiz struct {
	store RegisterStore
	hash  Hash
}

func NewRegisterBiz(store RegisterStore, hash Hash) *registerBiz {
	return &registerBiz{store: store, hash: hash}
}

func (biz *registerBiz) Register(ctx context.Context, data *usermodel.UserCreate) error {
	if err := data.Validate(); err != nil {
		return err
	}

	user, err := biz.store.FindUser(ctx, map[string]interface{}{"email": data.Email})

	if user != nil {
		return common.ErrEntityExisted(usermodel.EntityName, err)
	}

	if err == common.ErrRecordNotFound {
		salt := common.GenSalt(50)

		data.Password = biz.hash.Hash(data.Password + salt)
		data.Salt = salt

		if err := biz.store.CreateUser(ctx, data); err != nil {
			return common.ErrCannotCreateEntity(usermodel.EntityName, err)
		}
	} else {
		return common.ErrDB(err)
	}

	return nil
}
