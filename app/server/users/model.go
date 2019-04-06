package users

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

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
func NewUser(b *PostBody) *User {
	u := &User{
		Email:       b.Email,
		DisplayName: b.DisplayName,
	}
	// password
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	u.Password = string(hash[:])
	// passwordURL
	switch {
	case b.PasswordURL == nil:
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: "",
		}
	case *b.PasswordURL == "":
		u.PasswordURL = sql.NullString{
			Valid:  false,
			String: "",
		}
	default:
		u.PasswordURL = sql.NullString{
			Valid:  true,
			String: *b.PasswordURL,
		}

	}
	return u
}

// FindAll users in database
func FindAll(orm *xorm.Engine) ([]User, error) {
	users := []User{}
	err := orm.Find(&users)
	return users, err
}

// Find user in database
func (u *User) Find(orm *xorm.Engine) (int, error) {
	found, err := orm.Get(u)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("user not found")
	}
	return http.StatusOK, nil
}

// Save user to database
func (u *User) Save(orm *xorm.Engine) (int, error) {
	var (
		err      error
		affected int64
	)
	affected, err = orm.Where("email = ?", u.Email).Count(&User{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected != 0 {
		return http.StatusConflict, errors.New("such user always exists")
	}

	u.Created = uint64(time.Now().UTC().Unix())
	u.Updated = u.Created
	affected, err = orm.InsertOne(u)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected == 0 {
		return http.StatusUnprocessableEntity, errors.New("db refused to insert such user")
	}

	return http.StatusCreated, nil
}

// Update user in database
func (u *User) Update(orm *xorm.Engine) (int, error) {
	var (
		err      error
		found    bool
		user     User
		affected int64
	)
	// get old user data (and check if user exists)
	found, err = orm.ID(u.ID).Get(&user)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("user not found")
	}
	err = u.setFieldsFrom(user)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	u.Updated = uint64(time.Now().UTC().Unix())
	affected, err = orm.Update(u)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected != 0 {
		return http.StatusUnprocessableEntity, errors.New("db refused to update")
	}
	return http.StatusOK, nil
}

// Delete user from database
func (u *User) Delete(orm *xorm.Engine) (int, error) {
	var (
		err      error
		found    bool
		affected int64
		user     User
	)
	// check if user exists
	found, err = orm.ID(u.ID).Get(&user)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if !found {
		return http.StatusNotFound, errors.New("user not exists")
	}
	//delete
	affected, err = orm.ID(u.ID).Delete(&User{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected == 0 {
		return http.StatusUnprocessableEntity, errors.New("db refused to delete user")
	}
	return http.StatusOK, nil
}

//------------------------------------------------------------------------------
func (u *User) setFieldsFrom(user User) error {
	if len(u.Email) == 0 {
		u.Email = user.Email
	}
	if len(u.Password) == 0 {
		u.Password = user.Password
	} else {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hash[:8])
	}
	return nil
}
