package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	guuid "github.com/google/uuid"
)

// Create User Table
func CreateUserTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}
	createError := db.Model(&User{}).CreateTable(opts)

	//createError := db.CreateTable(&User{}, opts)
	if createError != nil {
		log.Printf("Error while creating User table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("User table created")
	return nil
}

func GetAllUsers(c *gin.Context) {
	var users []User

	err := dbConnect.Model(&users).Select()
	if err != nil {
		log.Printf("Error while getting all users, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "All Users",
		"data":    users,
	})

}

func CreateUser(c *gin.Context) {
	var user User
	c.BindJSON(&user)
	name := user.Name
	id := guuid.New().String()
	_, insertError := dbConnect.Model(&User{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Insert()

	if insertError != nil {
		log.Printf("Error while inserting new user into db, Reason: %v\n", insertError)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":  http.StatusCreated,
		"message": "User created Successfully",
	})

}

func GetSingleUser(c *gin.Context) {
	userId := c.Param("userId")
	user := &User{ID: userId}
	err := dbConnect.Model(user).WherePK().Select()
	//err := dbConnect.Select(user)
	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}
}

func EditUserName(c *gin.Context) {
	userId := c.Param("userId")
	var user User
	c.BindJSON(&user)
	name := user.Name
	_, err := dbConnect.Model(&User{}).Set("name = ?", name).Where("id = ?", userId).Update()
	if err != nil {
		log.Printf("Error, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "User Address Edited Successfully",
	})

}

func DeleteUser(c *gin.Context) {
	userId := c.Param("userId")
	user := &User{ID: userId}
	_, err := dbConnect.Model(user).WherePK().Delete()
	//err := dbConnect.Delete(user)
	if err != nil {
		log.Printf("Error while deleting a single user, Reason: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Something went wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "User deleted successfully",
	})

}

func GetUsersAllBaskets(c *gin.Context) {

	userId := c.Param("userId")
	user := &User{ID: userId}

	err := dbConnect.Model(user).Relation("Baskets").Relation("Baskets.BasketProducts").Relation("Baskets.BasketProducts.Product").WherePK().Select()
	//basket.BasketProducts[].Product
	if err != nil {
		log.Printf("Error while getting a user baskets, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Baskets not found",
		})
		return
	}

}

func GetUsersActiveBasket(c *gin.Context) {

	userId := c.Param("userId")
	user := &User{ID: userId}

	err := dbConnect.Model(user).Relation("Baskets").Relation("Baskets.BasketProducts").Relation("Baskets.BasketProducts.Product").WherePK().Select()
	//basket.BasketProducts[].Product
	if err != nil {
		log.Printf("Error while getting a user baskets, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Baskets not found",
		})
		return
	}

	nameString := "Active basket of :" + user.Name
	for _, b := range user.Baskets {
		if b.BasketStatus == 1 {
			c.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": nameString,
				"data":    b,
			})
		}
	}
}
