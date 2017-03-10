package users

import (
	"github.com/corvinusz/echo-xorm/utils"
	"github.com/go-xorm/xorm"
	"golang.org/x/crypto/bcrypt"
)

// User is an entity (here are DB definitions)
type User struct {
	ID       uint64 `xorm:"'id' pk autoincr unique notnull" json:"id"`
	Login    string `xorm:"text index not null unique 'login'" json:"login"`
	Hash     string `xorm:"'hash' text unique" json:"hash"`
	Password string `xorm:"text not null 'password'" json:"-"`
	Created  uint64 `xorm:"created" json:"-"` // too lazy to extract
	Updated  uint64 `xorm:"updated" json:"-"` // too lazy to extract
}

// TableName used by xorm to set table name for entity
func (u *User) TableName() string {
	return "users"
}

// FindAll saves user to database
func (u *User) FindAll(orm *xorm.Engine) ([]User, error) {
	var (
		users []User
		err   error
	)
	err = orm.Find(&users)
	return users, err
}

// Find extract 1 user from database
func (u *User) Find(orm *xorm.Engine) (bool, error) {
	return orm.ID(u.ID).Get(u)
}

// Save inserts user to database
func (u *User) Save(orm *xorm.Engine) (int64, error) {
	var (
		err      error
		hash     []byte
		affected int64
	)
	affected, err = orm.Where("login = ?", u.Login).Count(&User{})
	if err != nil {
		return 0, err
	}
	if affected != 0 {
		return 0, nil
	}

	// encrypt password
	hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	u.Password = string(hash[:])
	u.Hash = utils.GetSHA3Hash(u.Login)
	affected, err = orm.InsertOne(u)
	return affected, err
}

// Update updates user to database
func (u *User) Update(orm *xorm.Engine) (int64, error) {
	var (
		err   error
		found bool
		user  User
	)
	// check if user exists
	found, err = orm.ID(u.ID).Get(&user)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, nil
	}
	//update
	err = u.updateFields(&user)
	if err != nil {
		return 0, nil
	}
	return orm.ID(user.ID).Update(&user)
}

// Delete updates user to database
func (u *User) Delete(orm *xorm.Engine) (int64, error) {
	var (
		err   error
		found bool
		user  User
	)
	// check if user exists
	found, err = orm.ID(u.ID).Get(&user)
	if err != nil {
		return 0, err
	}
	if !found {
		return 0, nil
	}
	//delete
	return orm.ID(u.ID).Delete(&User{})
}

//------------------------------------------------------------------------------
func (u *User) updateFields(user *User) error {
	if len(u.Login) != 0 {
		user.Login = u.Login
		user.Hash = utils.GetSHA3Hash(u.Login)
	}
	if len(u.Password) != 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hash[:])
	}
	return nil
}
