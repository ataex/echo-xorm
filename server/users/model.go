package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity (here are DB definitions)
type User struct {
	ID            uint64 `xorm:"'id' pk autoincr unique notnull" json:"id"`
	Email         string `xorm:"text index not null unique 'email'" json:"email"`
	DisplayName   string `xorm:"'display_name'" json:"display_name"`
	Password      string `xorm:"text not null 'password'" json:"-"`
	PasswordEtime uint64 `xorm:"'password_etime'" json:"password_etime"`
	PasswordURL   string `xorm:"'password_url'" json:"password_url"`
	Created       uint64 `xorm:"created" json:"created"`
	Updated       uint64 `xorm:"updated" json:"updated"`
}

// TableName used by xorm to set table name for entity
func (u *User) TableName() string {
	return "users"
}

// FindAll users in database
func (u *User) FindAll(orm *xorm.Engine) ([]User, error) {
	var (
		users []User
		err   error
	)
	err = orm.Find(&users)
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
		hash     []byte
		affected int64
	)
	affected, err = orm.Where("email = ?", u.Email).Count(&User{})
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	if affected != 0 {
		return http.StatusConflict, errors.New("such user always exists")
	}

	// encrypt password
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}
	u.Password = string(hash[:])

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
