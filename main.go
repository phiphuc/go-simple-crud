package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

//CREATE TABLE `restaurants` (
//`id` int(11) NOT NULL AUTO_INCREMENT,
//`owner_id` int(11) DEFAULT NULL,
//`name` varchar(50) NOT NULL,
//`addr` varchar(255) NOT NULL,
//`city_id` int(11) DEFAULT NULL,
//`lat` double DEFAULT NULL,
//`lng` double DEFAULT NULL,
//`cover` json DEFAULT NULL,
//`logo` json DEFAULT NULL,
//`shipping_fee_per_km` double DEFAULT '0',
//`status` int(11) NOT NULL DEFAULT '1',
//`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
//`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
//PRIMARY KEY (`id`),
//KEY `owner_id` (`owner_id`) USING BTREE,
//KEY `city_id` (`city_id`) USING BTREE,
//KEY `status` (`status`) USING BTREE
//) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

type Restaurants struct {
	Id   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
	Addr string `json:"address" gorm:"column:addr"`
}

func (Restaurants) TableName() string {
	return "restaurants"
}

func main() {
	dsn := os.Getenv("DBConnectionStr")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	if err := runService(db); err != nil {
		log.Println(err)
	}

	//newNote := Note{Title: "Demo", Content: "This is content"}
	//if err := db.Create(&newNote); err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(newNote)

	//var notes []Note
	//db.Where("status=?", 1).Find(&notes)
	//fmt.Println(notes)
	//
	//var note Note
	//if err := db.Where("id", 5).First(&note); err != nil {
	//	log.Println(err)
	//}
	//fmt.Println(note)

	//db.Table(Note{}.TableName()).Where("id", 4).Delete(nil)

	//note.Title = "Demo 3"
	//db.Table(Note{}.TableName()).
	//	Where("id", 5).Updates(&note)
}

func runService(db *gorm.DB) error {
	r := gin.Default()
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//CRUD
	restaurants := r.Group("/restaurants")
	{
		restaurants.POST("", func(c *gin.Context) {
			var data Restaurants

			if err := c.ShouldBind(&data); err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			if err := db.Create(&data).Error; err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, data)
		})

		restaurants.GET("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))

			if err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			var data Restaurants

			if err := db.Where("id", id).First(&data).Error; err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, data)
		})

		restaurants.GET("", func(c *gin.Context) {

			var data []Restaurants

			if err := db.Table(Restaurants{}.TableName()).Find(&data).Error; err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, data)
		})

		restaurants.PUT("", func(c *gin.Context) {
			var data Restaurants

			if err := c.ShouldBind(&data); err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			if err := db.Table(Restaurants{}.TableName()).Where("id=?", data.Id).Updates(&data).Error; err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, data)
		})

		restaurants.DELETE("/:id", func(c *gin.Context) {
			id, err := strconv.Atoi(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			if err := db.Table(Restaurants{}.TableName()).Where("id", id).Delete(nil).Error; err != nil {
				c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, true)
		})
	}

	return r.Run()
}
