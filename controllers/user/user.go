package user

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pg "github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	guuid "github.com/google/uuid"
	basketController "github.com/mataliksamil/Go_Bootcamp_Final/controllers/basket"
	entities "github.com/mataliksamil/Go_Bootcamp_Final/entities"
)

// Create User Table
func CreateUserTable(db *pg.DB) error {
	opts := &orm.CreateTableOptions{
		FKConstraints: true,
		IfNotExists:   true,
	}
	createError := db.Model(&entities.User{}).CreateTable(opts)

	//createError := db.CreateTable(&User{}, opts)
	if createError != nil {
		log.Printf("Error while creating User table, Reason: %v\n", createError)
		return createError
	}
	log.Printf("User table created")
	return nil
}

func GetAllUsers(c *gin.Context) {
	var users []entities.User

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
	var user entities.User
	c.BindJSON(&user)
	name := user.Name
	id := guuid.New().String()
	_, insertError := dbConnect.Model(&entities.User{
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
	user := &entities.User{ID: userId}
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
	var user entities.User
	c.BindJSON(&user)
	name := user.Name
	_, err := dbConnect.Model(&entities.User{}).
		Set("name = ?", name).
		Where("id = ?", userId).
		Update()

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
	user := &entities.User{ID: userId}
	_, err := dbConnect.Model(user).
		WherePK().
		Delete()

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
	user := &entities.User{ID: userId}
	// JOIN for whole structure under a user
	err := dbConnect.Model(user).
		Relation("Baskets").
		Relation("Baskets.BasketProducts").
		Relation("Baskets.BasketProducts.Product").
		WherePK().
		Select()
	//basket.BasketProducts[].Product
	if err != nil {
		log.Printf("Error while getting a user baskets, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Baskets not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": user,
	})

}

func GetUsersActiveBasket(c *gin.Context) {

	userId := c.Param("userId")
	err := basketController.ApplyDiscount(userId)
	if err != nil {
		log.Printf("Error while Applying discounts, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Baskets not found",
		})
		return
	}

	var myBasket = &entities.Basket{}
	err = dbConnect.Model(myBasket).Where("user_id=?", userId).Where("basket_status=?", 1).Select()
	if err != nil {
		log.Printf("Error while getting a user baskets, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Baskets not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ACTIVE BASKET : ",
		"data":    myBasket,
	})
}
