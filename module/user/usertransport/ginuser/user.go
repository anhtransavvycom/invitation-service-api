package ginuser

import (
	"app-invite-service/common"
	"app-invite-service/component"
	"app-invite-service/component/hash"
	"app-invite-service/component/tokenprovider/jwt"
	"app-invite-service/module/user/userbiz"
	"app-invite-service/module/user/usermodel"
	"app-invite-service/module/user/userstorage"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      Login
// @Description  Login with existing account
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email     formData  string  true  "email"
// @Param        password  formData  string  true  "password"
// @Success      200       {object}  common.SuccessRes{data=usermodel.Account}
// @Failure      500       {object}  common.AppError
// @Failure      400       {object}  common.AppError
// @Router       /login [post]
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

// Register godoc
// @Summary      Register
// @Description  Register for new account
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        email     formData  string  true  "email"
// @Param        password  formData  string  true  "password"
// @Success      200       {object}  common.SuccessRes{data=int}
// @Failure      500       {object}  common.AppError
// @Failure      400       {object}  common.AppError
// @Router       /register [post]
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

// GenerateInviteToken godoc
// @Summary      Generate invitation token
// @Description  Generate an invitation token
// @Tags         token
// @Param        Authorization  header    string   true  "Authorization header"
// @Success      200            {object}  common.SuccessRes{data=usermodel.InvitationToken}
// @Failure      500            {object}  common.AppError
// @Router       /tokens/generate [post]
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

// LoginWithInviteToken godoc
// @Summary      Login with invitation token
// @Description  Login with invitation token
// @Tags         auth
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        invitation_token  formData  string  true  "invitation token"
// @Success      200               {object}  common.SuccessRes{data=usermodel.Account}
// @Failure      500               {object}  common.AppError
// @Router       /login/invitation [post]
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

// ValidateInvitationToken godoc
// @Summary      Validate invitation token
// @Description  check weather invitation token is valid
// @Tags         token
// @Param        token  path  string  true  "invitation token to be validated"
// @Success      200
// @Failure      500            {object}  common.AppError
// @Router       /tokens/{token}/validation [post]
func ValidateInvitationToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewValidateInviteTokenBiz(redis)
		if err := biz.ValidateInvitationToken(c.Request.Context(), c.Param("token")); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]bool{"success": true}))
	}
}

// ListInvitationToken godoc
// @Summary      List invitation tokens
// @Description  List invitation tokens
// @Tags         token
// @Param        status         query     integer  true  "token status"  Enums(0, 1)
// @Param        Authorization  header    string   true  "Authorization header"
// @Success      200            {object}  common.SuccessRes{data=[]usermodel.InvitationToken}
// @Failure      500  {object}  common.AppError
// @Router       /tokens [get]
func ListInvitationToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var filter usermodel.InvitationTokenFilter
		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewListInvitationTokenBiz(redis)

		result, err := biz.ListInvitationToken(c.Request.Context(), &filter)
		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, nil, filter))
	}
}

// UpdateInvitationToken godoc
// @Summary      Update an invitation token
// @Description  Update an invitation token
// @Tags         token
// @Accept       x-www-form-urlencoded
// @Produce      json
// @Param        token          path      string   true  "token"
// @Param        status         formData  integer  true  "token status"  Enums(0, 1)
// @Param        Authorization  header    string  true  "Authorization header"
// @Success      200
// @Failure      500  {object}  common.AppError
// @Failure      400  {object}  common.AppError
// @Router       /tokens/{token} [put]
func UpdateInvitationToken(appCtx component.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data usermodel.InvitationTokenUpdate
		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		redis := appCtx.GetRedisConnection()
		biz := userbiz.NewUpdateInvitationTokenBiz(redis)
		if err := biz.UpdateInvitationToken(c.Request.Context(), c.Param("token"), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(map[string]bool{"success": true}))
	}
}
