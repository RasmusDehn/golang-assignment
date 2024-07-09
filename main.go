package main

import (
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

var DB *gorm.DB

type Movie struct {
	ID       uint    `json:"id" gorm:"primary_key"`
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Year     int     `json:"year"`
	Price    float64 `json:"price"`
}

type MovieStruct struct {
	Title    string  `json:"title" binding:"required"`
	Director string  `json:"director" binding:"required"`
	Year     int     `json:"year"`
	Price    float64 `json:"price"`
}

type UpdateMovieStruct struct {
	Title    string  `json:"title"`
	Director string  `json:"director"`
	Year     int     `json:"year"`
	Price    float64 `json:"price"`
}

func ConnectDatabase() {

	database, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Movie{})
	if err != nil {
		return
	}

	DB = database
}

func FindMovies(c *gin.Context) {
	var movies []Movie
	DB.Find(&movies)

	c.JSON(http.StatusOK, gin.H{"data": movies})
}

func FindMovie(c *gin.Context) {
	var movie Movie
	if err := DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": movie})

}

func UpdateMovie(c *gin.Context) {
	var movie Movie
	if err := DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	var input UpdateMovieStruct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	DB.Model(&movie).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": movie})

}

func CreateMovie(c *gin.Context) {
	var input MovieStruct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie := Movie{Title: input.Title, Director: input.Director}
	DB.Create(&movie)

	c.JSON(http.StatusOK, gin.H{"data": movie})
}

func DeleteMovie(c *gin.Context) {
	var movie Movie
	if err := DB.Where("id = ?", c.Param("id")).First(&movie).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}
	DB.Delete(&movie)
	c.JSON(http.StatusOK, gin.H{"data": true})
}

func main() {
	r := gin.Default()

	ConnectDatabase()

	r.GET("/movies", FindMovies)
	r.GET("/movie/:id", FindMovie)
	r.DELETE("/movie/:id", DeleteMovie)
	r.POST("/movie", CreateMovie)
	r.PATCH("/movie/:id", UpdateMovie)
	r.Run()
}
