package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"GoProject/controllers"
	"GoProject/models"
	"GoProject/repositories"
	"GoProject/service"

	_ "github.com/go-sql-driver/mysql"
)

func main() {

	dsn := "root:@tcp(127.0.0.1:3307)/pustaka-barang?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal Connect Ke Database")
	}
	if err := db.AutoMigrate(&models.Barang{}, &models.User{}); err != nil {
		panic("Failed to run migration: " + err.Error())
	}
	barangRepo := repositories.NewRepository(db)
	barangService := service.NewService(barangRepo)
	barangController := controllers.NewBarangController(barangService)

	authController := controllers.NewAuthController(db)

	router := gin.Default()

	store := cookie.NewStore([]byte("secret_key"))
	router.Use(controllers.NoCache())
	router.Use(sessions.Sessions("session", store))

	authGroup := router.Group("/api")
	authGroup.Use(controllers.AuthMiddleware())
	authGroup.GET("/home", barangController.Home)
	authGroup.GET("/index", barangController.Index)
	authGroup.GET("/create", barangController.CreatePage)
	authGroup.POST("/create", barangController.Create)
	authGroup.GET("/edit/:id", barangController.EditPage)
	authGroup.POST("/edit/:id", barangController.Edit)
	authGroup.GET("/delete/:id", barangController.Delete)

	router.StaticFS("/api/assets", http.Dir("./views/assets"))

	authGroup.GET("/register", authController.RegisterPage)
	router.POST("/register", authController.Register)
	// router.GET("/", authController.LoginPage)
	authGroup.GET("/login", authController.LoginPage)
	router.POST("/login", authController.Login)
	router.GET("/logout", authController.Logout)

	router.Run(":8686")
	fmt.Println("Server Berjalan di 8686")

}
