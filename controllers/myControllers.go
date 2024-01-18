package controllers

import (
	"GoProject/models"
	"GoProject/service"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BarangController interface {
	Index(ctx *gin.Context)
	Home(ctx *gin.Context)
	FindById(ctx *gin.Context)
	CreatePage(ctx *gin.Context)
	Create(ctx *gin.Context)
	EditPage(ctx *gin.Context)
	Edit(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type barangController struct {
	barangService service.BarangService
}

func NewBarangController(serv service.BarangService) BarangController {
	return &barangController{
		barangService: serv,
	}
}

func (c *barangController) Home(ctx *gin.Context) {

	var models models.User

	layout := path.Join("views/templates/pages/layout.html")
	base := path.Join("views/templates/base/home.html")

	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.ExecuteTemplate(ctx.Writer, "home.html", gin.H{"user": models})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func (c *barangController) Index(ctx *gin.Context) {
	layout := "views/templates/pages/layout.html"
	base := "views/templates/base/index.html"

	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	barangs, err := c.barangService.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.ExecuteTemplate(ctx.Writer, "index.html", gin.H{"barangs": barangs})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (c *barangController) FindById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	barang, err := c.barangService.FindById(int(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	ctx.JSON(http.StatusOK, barang)
}

func (c *barangController) CreatePage(ctx *gin.Context) {
	layout := path.Join("views/templates/pages/layout.html")
	base := path.Join("views/templates/base/create.html")

	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tmpl.ExecuteTemplate(ctx.Writer, "create.html", nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (c *barangController) Create(ctx *gin.Context) {

	nama_barang := ctx.PostForm("nama_barang")
	stockStr := ctx.PostForm("stock")
	hargaStr := ctx.PostForm("harga")

	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Stok harus berupa angka"})
		return
	}

	harga, err := strconv.Atoi(hargaStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Harga harus berupa angka"})
		return
	}

	barangg := models.Barang{
		NamaBarang: nama_barang,
		Stock:      stock,
		Harga:      harga,
	}

	erro := c.barangService.Create(&barangg)
	if erro != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Redirect(http.StatusSeeOther, "/api/index")

}

func (c *barangController) EditPage(ctx *gin.Context) {

	log.Println("EditPage function called")
	layout := path.Join("views/templates/pages/layout.html")
	base := path.Join("views/templates/base/edit.html")
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		log.Println("Invalid ID")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	existingBarang, err := c.barangService.FindById(int(id))
	if err != nil {
		log.Println("Barang not found")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}
	tmpl, err := template.ParseFiles(layout, base)
	if err != nil {
		log.Println("Error parsing template files")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = tmpl.ExecuteTemplate(ctx.Writer, "edit.html", gin.H{"barang": existingBarang})
	if err != nil {
		log.Println("Error executing template")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("EditPage function completed successfully")
}

// Edit method in barangController
func (c *barangController) Edit(ctx *gin.Context) {
	var barang models.Barang
	if err := ctx.ShouldBind(&barang); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingBarang, err := c.barangService.FindById(int(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Barang not found"})
		return
	}

	existingBarang.NamaBarang = barang.NamaBarang
	existingBarang.Stock = barang.Stock
	existingBarang.Harga = barang.Harga

	if err := c.barangService.Edit(existingBarang); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("After Edit: %+v\n", existingBarang)

	// Redirect to the appropriate route
	ctx.Redirect(http.StatusSeeOther, "/api/index")
}

func (c *barangController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	fmt.Println("Attempting to delete barang with ID:", id)

	if err := c.barangService.Delete(int(id)); err != nil {
		fmt.Println("Error deleting barang:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/api/index")
}
