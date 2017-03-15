package bddtests_test

/*
import (
	"encoding/json"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/corvinusz/echo-xorm/server/users"
)

// Do not use name starting with Test... to avoid automatic call of test function
func (suite *LsxTestSuite) testPostUsers(t *testing.T) {
	// working part
	Convey("POST /users", t, func() {
		input := users.UserInput{
			Login:    "a_test_user_100",
			Password: "a_test_user_100",
		}
		result := new(users.User)
		resp, err := suite.rc.R().SetBody(input).SetResult(result).Post("/users")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 201)
		So(result.Login, ShouldEqual, input.Login)
		So(result.Created, ShouldNotEqual, 0)
		So(result.Updated, ShouldEqual, 0)
	})
	//error checks
	Convey("POST /users with bad payload", t, func() {
		data := "this is not a json; ' select 1;"
		resp, err := suite.rc.R().SetBody(data).Post("/users")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("POST /users with deficient payload", t, func() {
		input := users.UserInput{
			Login: "Filler",
		}
		resp, err := suite.rc.R().SetBody(input).Post("/users")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("POST /users with existing login", t, func() {
		input := users.UserInput{
			Login:    "a_test_user_100",
			Password: "filler",
		}
		data, err := json.Marshal(input)
		So(err, ShouldBeNil)
		resp, err := suite.rc.R().SetBody(data).Post("/users")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 409)
	})
}

//------------------------------------------------------------------------------
func (suite *LsxTestSuite) testGetUsers(t *testing.T) {
	//work part
	Convey("GET /users", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 8)
	})
	Convey("GET /users with limit", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?limit=3")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 3)
		So(result[0].ID, ShouldEqual, 1)
		So(result[1].ID, ShouldEqual, 2)
		So(result[2].ID, ShouldEqual, 3)
	})
	Convey("GET /users with offset", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?offset=2")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 6)
		So(result[0].ID, ShouldEqual, 3)
		So(result[1].ID, ShouldEqual, 4)
	})
	Convey("GET /users with limit and offset", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?limit=3&offset=4")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 3)
		So(result[0].ID, ShouldEqual, 5)
		So(result[1].ID, ShouldEqual, 6)
		So(result[2].ID, ShouldEqual, 7)
	})
	Convey("GET /users?id=8", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?id=8")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(result[0].ID, ShouldEqual, 8)
		So(result[0].Login, ShouldEqual, "a_test_operator_08")
	})
	Convey("GET /users?login=", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?login=a_test_user_06")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 1)
		So(result[0].ID, ShouldEqual, 6)
		So(result[0].Login, ShouldEqual, "a_test_user_06")
	})
	//error checks
	Convey("GET /users?id=err", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?id=1005001")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 0)
	})
	Convey("GET /users?name=not-existing-name", t, func() {
		result := []users.User{}
		resp, err := suite.rc.R().SetResult(&result).Get("/users?login=not-existing-name")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 0)
	})
}

//------------------------------------------------------------------------------
func (suite *LsxTestSuite) testPutUsers(t *testing.T) {
	// working part
	Convey("PUT /users/{id}", t, func() {
		input := users.UserInput{
			Login:    "a_updated_test_user_20",
			Password: "a_updated_test_user_20",
		}
		result := new(users.User)
		data, err := json.Marshal(input)
		So(err, ShouldBeNil)
		resp, err := suite.rc.R().SetBody(data).SetResult(result).Put("/users/6")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(result.Login, ShouldEqual, input.Login)
		So(result.Updated, ShouldAlmostEqual, time.Now().UTC().Unix())
	})
	// errors check
	Convey("PUT /users with bad payload", t, func() {
		data := "this is not a json; ' select 1;"
		resp, err := suite.rc.R().SetBody(data).Put("/users/6")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("PUT /users with bad id", t, func() {
		data := `{"name":"filler","email":"filler","password":"filler","group_id":1}`
		resp, err := suite.rc.R().SetBody(data).Put("/users/blabla6")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("PUT /users with non-existent id", t, func() {
		data := `{"name":"filler","email":"filler","password":"filler","group_id":1}`
		resp, err := suite.rc.R().SetBody(data).Put("/users/100506")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 404)
	})
	Convey("PUT /users with non-existent group_id", t, func() {
		data := `{"name":"filler","email":"filler","password":"filler","group_id":1002}`
		resp, err := suite.rc.R().SetBody(data).Put("/users/6")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
}

//------------------------------------------------------------------------------
func (suite *LsxTestSuite) testDeleteUsers(t *testing.T) {
	// working part
	Convey("DELETE /users/{id}", t, func() {
		resp, err := suite.rc.R().Delete("/users/7")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		// check actual data
		resultFromGet := []users.User{}
		resp, err = suite.rc.R().SetResult(&resultFromGet).Get("/users?id=7")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(resultFromGet), ShouldEqual, 0)
	})
	// error testing
	Convey("DELETE /users with bad id", t, func() {
		resp, err := suite.rc.R().Delete("/users/blakdks")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 400)
	})
	Convey("DELETE non-existent /users", t, func() {
		resp, err := suite.rc.R().Delete("/users/1005001")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 404)
	})

}
*/
