package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/golang/glog"

	"github.com/mayadata-io/kubera-auth/pkg/models"
	controller "github.com/mayadata-io/kubera-auth/versionedController/v1"
)

// UserController is a type to be accepted as input
type UserController struct {
	controller.GenericController
	routePath string
	model     *models.UserCredentials
}

// New creates a new User
func New() *UserController {
	return &UserController{
		routePath: controller.UserRoute,
		model:     &models.UserCredentials{},
	}
}

// Put updates a user details
func (user *UserController) Put(c *gin.Context) {
	err := c.BindJSON(user.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userModel := models.UserCredentials(*user.model)
	controller.Server.UpdateUserDetailsRequest(c, &userModel)
	return
}

//Patch updates the password of concerned user ggiven that request should be sent by admin
func (user *UserController) Patch(c *gin.Context) {
	err := c.BindJSON(user.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	controller.Server.ResetPasswordRequest(c, user.model.GetPassword(), user.model.GetUserName())
	return
}

//Post creates a user, request should be sent by admin
func (user *UserController) Post(c *gin.Context) {
	err := c.BindJSON(user.model)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotAcceptable, gin.H{
			"message": "Unable to parse JSON",
		})
		return
	}

	userModel := models.UserCredentials(*user.model)
	userModel.Kind = models.LocalAuth
	userModel.Role = models.RoleUser
	userModel.State = models.StateCreated
	controller.Server.CreateRequest(c, &userModel)
	return
}

// Get will respond with a particular user or all users
func (user *UserController) Get(c *gin.Context) {
	// Get all users
	controller.Server.GetUsersRequest(c)
	return
}

// GetByUID will respond with a particular user or all users
func (user *UserController) GetByUID(c *gin.Context) {
	userID := c.Param("userID")
	controller.Server.GetUserByUID(c, userID)
}

// GetByUsername will respond with a particular user or all users
func (user *UserController) GetByUsername(c *gin.Context) {
	userID := c.Param("username")
	controller.Server.GetUserByUserName(c, userID)
}

// Register will rsgister this controller to the specified router
func (user *UserController) Register(router *gin.RouterGroup) {
	controller.RegisterController(router, user, user.routePath)
	router.GET(user.routePath+"/uid/:userID", user.GetByUID)
	router.GET(user.routePath+"/username/:username", user.GetByUsername)
}
