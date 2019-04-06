package users

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/corvinusz/echo-xorm/pkg/errors"

	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity (here are DB definitions)
type User struct {
	ID            uint64         `xorm:"'id' pk autoincr unique notnull" json:"id"`
	Email         string         `xorm:"'email' text index not null unique" json:"email"`
	DisplayName   string         `xorm:"'display_name' text" json:"display_name"`
	Password      string         `xorm:"'password' text not null" json:"-"`
	PasswordEtime uint64         `xorm:"'password_etime'" json:"password_etime"`
	PasswordURL   sql.NullString `xorm:"'password_url' text unique" json:"password_url"`
	Created       uint64         `xorm:"created" json:"created"`
	Updated       uint64         `xorm:"updated" json:"updated"`
}

// TableName used by xorm to set table name for entity
func (u *User) TableName() string {
	return "users"
}

// NewUser creates user from request body
// returns *User with data from body
// returns nil if error occured
func NewUser(b *PostBody) (*User, error) {
	u := &User{
		Email:       b.Email,
		DisplayName: b.DisplayName,
		Password:    b.Password,
	}
	// passwordURL
	switch {
	case b.PasswordURL == nil:
		u.PasswordURL = sql.NullString{
			Valid: false,
		}
	case *b.PasswordURL == "":
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: "",
		}
	default:
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: *b.PasswordURL,
		}

	}
	return u, nil
}

// FindAll users in database
func FindAll(orm *xorm.Engine) ([]User, error) {
	users := []User{}
	err := orm.Find(&users)
	if err != nil {
		return nil, errors.New("database error; " + err.Error())
	}
	return users, nil
}

// FindOne finds in database first user with matched struct fields
func (u *User) FindOne(orm *xorm.Engine) error {
	found, err := orm.Get(u)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}
	return nil
}

// Save user to database
func (u *User) Save(orm *xorm.Engine) error {
	affected, err := orm.Where("email = ?", u.Email).Count(&User{})
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected != 0 {
		return errors.NewWithCode(http.StatusConflict, "email already exists")
	}

	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.NewWithPrefix(err, "generate hash from password error")
	}
	u.Password = string(hash[:])
	// created/updated
	u.Created = uint64(time.Now().UTC().Unix())
	u.Updated = u.Created
	// save to storage
	affected, err = orm.InsertOne(u)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected == 0 {
		return errors.New("database error; db refused to insert")
	}

	return nil
}

// Update user in database
func (u *User) Update(orm *xorm.Engine) error {
	var oldUser User
	// get old user data (and check if user exists)
	found, err := orm.ID(u.ID).Get(&oldUser)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}
	err = u.setDataToUpdate(oldUser)
	if err != nil {
		return err
	}
	u.Updated = uint64(time.Now().UTC().Unix())
	affected, err := orm.Update(u)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected != 0 {
		return errors.New("database error; db refused to update")
	}
	return nil
}

// Delete user from database
func (u *User) Delete(orm *xorm.Engine) error {
	var user User
	// check if user exists
	found, err := orm.ID(u.ID).Get(&user)
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if !found {
		return errors.NewWithCode(http.StatusNotFound, "user not found")
	}
	//delete
	affected, err := orm.ID(u.ID).Delete(&User{})
	if err != nil {
		return errors.NewWithPrefix(err, "database error")
	}
	if affected == 0 {
		return errors.New("database error; db refused to delete")
	}
	return nil
}

//------------------------------------------------------------------------------
func (u *User) setDataToUpdate(user User) error {
	if len(u.Email) == 0 {
		u.Email = user.Email
	}
	if len(u.Password) == 0 {
		u.Password = user.Password
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.NewWithPrefix(err, "password hash generation error")
		}
		u.Password = string(hash[:8])
	}
	return nil
}
