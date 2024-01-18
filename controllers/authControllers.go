package controllers

import (
	"GoProject/middleware"
	"GoProject/models"
	"html/template"
	"net/http"
	"path"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthControllers struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthControllers {
	return &AuthControllers{
		DB: db,
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Exclude "/login" from redirection
		if ctx.Request.URL.Path == "/api/register" || ctx.Request.URL.Path == "/logout" {
			ctx.Next()
			return
		}

		if ctx.Request.URL.Path == "/api/login" || ctx.Request.URL.Path == "/logout" {
			ctx.Next()
			return
		}

		if !isUserLoggedIn(ctx) {
			// Redirect to the login page with an error message
			ctx.Redirect(http.StatusSeeOther, "/api/login")
			ctx.Abort()
			return
		}

		// Get the token from the cookie
		tokenString, err := ctx.Cookie("token")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ada Auth"})
			ctx.Abort()
			return
		}

		// Verify the token and get the claims.
		claims, err := middleware.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			ctx.Abort()
			return
		}

		ctx.Set("id", claims.ID)

		sessions := sessions.Default(ctx)
		sessions.Set("token", tokenString)
		sessions.Save()

		ctx.Next()
	}
}

func (c *AuthControllers) RegisterPage(ctx *gin.Context) {
	layout := path.Join("views/templates/pages/layout-auth.html")
	base := path.Join("views/templates/base/register.html")

	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.ExecuteTemplate(ctx.Writer, "register.html", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (c *AuthControllers) Register(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, e.Field()+" Password Minimal 8 Karakter")
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user.Password = string(hashPassword)

	if err := c.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	// ctx.JSON(http.StatusOK, gin.H{"token": token})
	ctx.Redirect(http.StatusSeeOther, "/api/login")
}

func isUserLoggedIn(ctx *gin.Context) bool {
	// Periksa apakah pengguna memiliki token dalam cookie atau session
	token, err := ctx.Cookie("token")
	if err != nil {
		return false
	}

	sessions := sessions.Default(ctx)
	savedToken := sessions.Get("token")
	if savedToken == nil || savedToken.(string) != token {
		return false
	}

	return true
}

func (c *AuthControllers) LoginPage(ctx *gin.Context) {
	layout := path.Join("views/templates/pages/layout-auth.html")
	base := path.Join("views/templates/base/login.html")

	if isUserLoggedIn(ctx) {
		// Jika sudah login, redirect ke rute "/api/home"
		ctx.Redirect(http.StatusSeeOther, "/api/home")
		return
	}

	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.ExecuteTemplate(ctx.Writer, "login.html", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (c *AuthControllers) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBind(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		var validationErrors []string
		for _, e := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, e.Field()+" Password Minimal 8 Karakter")
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	var dbUser models.User
	if err := c.DB.Where("Username = ?", user.Username).First(&dbUser).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := middleware.GenerateToken(dbUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Set the JWT token as a cookie
	secure := gin.Mode() == gin.ReleaseMode
	ctx.SetCookie("token", token, 3600, "/", "localhost", secure, true) // Expires in 1 hour

	// Tambahkan token ke session
	sessions := sessions.Default(ctx)
	sessions.Set("token", token)
	err = sessions.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	ctx.Set("flash", gin.H{"success": "Login successful!"})

	// Redirect to home
	ctx.Redirect(http.StatusSeeOther, "/api/home")
}

func (c *AuthControllers) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)

	// Clear the session
	sessions := sessions.Default(ctx)
	sessions.Clear()
	err := sessions.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	// Redirect to the login page or any other appropriate page
	ctx.Redirect(http.StatusSeeOther, "/api/login")
}

func NoCache() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.Header("Pragma", "no-cache")
		ctx.Header("Expires", "0")
		ctx.Next()
	}
}
