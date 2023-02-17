package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database schemas
type User struct {
	ID        uint
	Name      string
	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Event     []Event `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type Event struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint
	Name        string
	Description string
	Timing      time.Time `gorm:"type:time"`
	Place       string
	NumSeats    int
	Price       int
	User        User `gorm:"foreignKey:UserID"`
}

type Ticket struct {
	ID      int `gorm:"primaryKey"`
	UserID  uint
	EventID uint
	User    User  `gorm:"foreignKey:UserID"`
	Event   Event `gorm:"foreignKey:EventID"`
}

// connection to database
func ConnectDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("db_url")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening DB", err)
	}
	fmt.Println("connected to DB")

	return db
}

// creates a new user if not already existing
func CreateUser(name string, email string, firstName string, lastName string, db *gorm.DB) {
	user := User{
		Name:      name,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	// checking if user already exists
	result := db.Where("email = ?", email).First(&user)

	if result.RowsAffected == 0 {
		result = db.Create(&user)
		if result.Error != nil {
			log.Fatal("Error creating user", result.Error)
		}
		fmt.Println("New user created", user)
	} else {
		fmt.Println(user, "already exists")
	}
}

// sets up database schema
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Event{}, &Ticket{})
}
