package model

import (
	"errors"
	"fmt"

	"github.com/hb-go/echo-web/module/auth"
	"github.com/hb-go/echo-web/module/log"
	"github.com/hb-go/gorm"

	"golang.org/x/crypto/bcrypt"
)

// GetUserByID is function find user by id
func (u *User) GetUserByID(id uint64) *User {
	user := User{}
	// var count int64
	// db := DB().Where("id = ?", id)
	// if err := Cache(db).First(&user).Count(&count).Error; err != nil {
	// 	log.Debugf("GetUserById error: %v", err)
	// 	return nil
	// }

	return &user
}

// GetUserByUserName is sign in function
func (u *User) GetUserByUserName() *User {
	newDb, err := newDB()
	if err := newDb.Preload("Account").Where("username = ?", u.UserName).First(&u).Error; err != nil {
		log.Debugf("GetUserByUserNamePwd error: %v", err)
		return nil
	}
	u.CheckPassword(u.Password)
	if err != nil {
		return nil
	}
	return u
}

// SetPassword function convert password to hash string.
func (u *User) SetPassword() error {
	if len(u.Password) == 0 {
		return errors.New("Password should not be empty!")
	}
	bytePassword := []byte(u.Password)
	// Make sure the second param `bcrypt generator cost` between [4, 32)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	u.Password = string(passwordHash)
	return nil
}

// CheckPassword check password
func (u *User) CheckPassword(password string) error {
	bytePassword := []byte(password)
	byteHashedPassword := []byte(u.Password)
	return bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)
}

// AddUserWithUserNamePwd create user.
func (u *User) AddUserWithUserNamePwd() *User {
	fmt.Println(u)
	u.SetPassword()
	newDb, err := newDB()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if err := newDb.Create(&u).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(u)
	return u
}

// User prepeo use system.
type User struct {
	gorm.Model
	// ID            uint64  `json:"id,omitempty"`
	UserName      string  `form:"username" json:"username,omitempty" gorm:"unique; type:varchar(25)"`
	Password      string  `form:"password" json:"password" gorm:"type:varchar(255)"`
	Email         string  `json:"email" gorm:"type:varchar(100);unique_index"`
	PhoneNumber   string  `json:"phone_number" gorm:"type:varchar(10)"`
	LastName      string  `json:"last_name" gorm:"type:varchar(25)"`
	FirstName     string  `json:"first_name" gorm:"type:varchar(25)"`
	Admin         bool    `json:"admin"`
	AccountID     uint    `form:"account_id" json:"account_id" gorm:"not null;"`
	Roles         []*Role `gorm:"many2many:user_role" json:"roles"`
	Account       Account `gorm:"ForeignKey:AccountID"`
	Authenticated bool    `form:"-" db:"-" json:"-"`
}

type Role struct {
	gorm.Model

	Name  string  `json:"name"`
	Users []*User `gorm:"many2many:user_role" json:"users"`
}

// GenerateAnonymousUser should generate an anonymous user model
// for all sessions. This should be an unauthenticated 0 value struct.
func GenerateAnonymousUser() auth.User {
	return &User{}
}

// TableName users
func (u User) TableName() string {
	return "user"
}

// Login will preform any actions that are required to make a user model
// officially authenticated.
func (u *User) Login() {
	u.Authenticated = true
}

func (u *User) GetAccount() string {
	return u.Account.Name
}

func (u *User) GetAccountID() uint {
	return u.Account.ID
}

// Logout will preform any actions that are required to completely
// logout a user.
func (u *User) Logout() {
	// Remove from logged-in user's list
	// etc ...
	u.Authenticated = false
}

// IsAuthenticated check status auth.
func (u *User) IsAuthenticated() bool {
	return u.Authenticated
}

// UniqueId insterface UniqeaID
func (u *User) UniqueId() interface{} {
	return u.ID
}

// GetById will populate a user object from a database model witha matching id.
func (u *User) GetById(id interface{}) error {
	if err := DB().Preload("Account").First(&u, id).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) GetUserByEmailPwd(email string, pwd string) *User {
	user := User{}
	if err := DB().Preload("Account").Where("email = ? ", email).First(&user).Error; err != nil {
		log.Debugf("GetUserByNicknamePwd error: %v", err)
		return nil
	}
	return &user
}

func GetUserList() []*User {
	users := []*User{}
	if err := DB().Find(&users).Error; err != nil {
		return nil
	}
	return users
}
