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
	ID        uint    `json:"-"`
	Email     string  `gorm:"unique" json:"email"`
	FirstName string  `json:"given_name"`
	LastName  string  `json:"family_name"`
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

var DB *gorm.DB

// connection to database
func ConnectDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("db_url")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening DB", err)
	}
	fmt.Println("Connected to Database")
}

// add a new event
func AddEvent(event Event) error {
	result := DB.Create(&event)
	return result.Error
}

func BookTicket(ticket Ticket) error {
	result := DB.Create(&ticket)
	return result.Error
}

// creates a new user if not already existing
func CreateUser(user User) error {
	// checking if user already exists
	result := DB.Where("email = ?", user.Email).First(&user)

	if result.RowsAffected == 0 {
		result = DB.Create(&user)
		fmt.Println("New user created", user)
		return result.Error
	} else {
		fmt.Println(user, "already exists")
	}
	return nil
}

// sets up database schema
func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Event{}, &Ticket{})
	fmt.Println("Migrated to Database")
}
