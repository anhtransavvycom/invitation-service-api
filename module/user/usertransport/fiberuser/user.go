package fiberuser

import (
	"app-invite-service/common"
	"app-invite-service/component"
	"app-invite-service/component/hash"
	"app-invite-service/component/tokenprovider/jwt"
	"app-invite-service/module/user/userbiz"
	"app-invite-service/module/user/usermodel"
	"app-invite-service/module/user/userstorage"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func Login(appCtx component.AppContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data usermodel.UserLogin

		if err := c.BodyParser(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()
		store := userstorage.NewSQLStore(db)
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
		md5 := hash.NewMd5Hash()
		tokenConfig := appCtx.GetTokenConfig()

		biz := userbiz.NewLoginBiz(store, tokenProvider, md5, tokenConfig)

		account, err := biz.Login(c.UserContext(), &data)
		if err != nil {
			panic(err)
		}

		return c.Status(http.StatusOK).JSON(common.SimpleSuccessResponse(account))
	}
}
