package handlers

import (
	"net/http"
	"os"
	"twitter-bookmarks-api/auth"
	"twitter-bookmarks-api/database"
	"twitter-bookmarks-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TwitterAuth(c *gin.Context) {
	state := auth.GenerateStateToken()
	authURL := auth.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

func TwitterCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if !auth.ValidateStateToken(state) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid state parameter"})
		return
	}

	token, err := auth.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to exchange code for token"})
		return
	}

	twitterUser, err := auth.GetTwitterUserInfo(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user info from Twitter"})
		return
	}

	user, err := database.CreateUser(
		c.Request.Context(),
		twitterUser.Data.ID,
		twitterUser.Data.Username,
		twitterUser.Data.Name,
		twitterUser.Data.ProfileImageURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to create user"})
		return
	}

	jwtToken, err := auth.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to generate token"})
		return
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	redirectURL := frontendURL + "/auth/callback?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func GetMe(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)

	user, err := database.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse{Message: "Logged out successfully"})
}
