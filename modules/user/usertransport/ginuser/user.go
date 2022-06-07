package ginuser

import (
	"app-invite-service/common"
	"app-invite-service/component"
	"app-invite-service/component/hash"
	"app-invite-service/component/tokenprovider/jwt"
	"app-invite-service/modules/user/userbiz"
	"app-invite-service/modules/user/usermodel"
	"app-invite-service/modules/user/userstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data usermodel.UserLogin

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()
		store := userstorage.NewSQLStore(db)
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
		md5 := hash.NewMd5Hash()
		tokenConfig := appCtx.GetTokenConfig()

		biz := userbiz.NewLoginBiz(store, tokenProvider, md5, tokenConfig)

		account, err := biz.Login(c.Request.Context(), &data)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}

func Register(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data usermodel.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := appCtx.GetMainDBConnection()
		store := userstorage.NewSQLStore(db)
		md5 := hash.NewMd5Hash()
		biz := userbiz.NewRegisterBiz(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}

func GenerateInviteToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewGenerateTokenBiz(redis)

		result, err := biz.GenerateToken(c.Request.Context())
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}

func LoginWithInviteToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data usermodel.UserLoginWithInviteToken

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		redis := appCtx.GetRedisConnection()
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())
		md5 := hash.NewMd5Hash()
		tokenConfig := appCtx.GetTokenConfig()

		biz := userbiz.NewLoginWithInviteTokenBiz(redis, tokenProvider, md5, tokenConfig)

		account, err := biz.LoginWithInviteToken(c.Request.Context(), &data)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}

func ValidateInviteToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("invitation_token")
		//var data usermodel.InvitationToken
		//
		//if err := c.ShouldBind(&data); err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}

		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewValidateInviteTokenBiz(redis)
		if err := biz.ValidateInviteToken(c.Request.Context(), token); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]bool{"data": true}))
	}
}

func ListInvitationToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter usermodel.InvitationTokenFilter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var paging common.Paging

		paging.Fulfill()

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewListInvitationTokenBiz(redis)

		result, err := biz.ListInvitationToken(c.Request.Context(), &filter, &paging)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
