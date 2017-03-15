package bddtests

/*
import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Do not use name starting with Test... to avoid automatic call of function
func (suite *LsxTestSuite) testGetGroups(t *testing.T) {
	//work part
	Convey("GET /groups", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 4)
	})
	Convey("GET /groups with limit", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?limit=3")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 3)
		So(result[0].ID, ShouldEqual, 1)
		So(result[1].ID, ShouldEqual, 4)
		So(result[2].ID, ShouldEqual, 5)
	})
	Convey("GET /groups with offset", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?offset=2")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 2)
		So(result[0].ID, ShouldEqual, 5)
		So(result[1].ID, ShouldEqual, 10)
	})
	Convey("GET /groups with limit and offset", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?limit=2&offset=1")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 2)
		So(result[0].ID, ShouldEqual, 4)
		So(result[1].ID, ShouldEqual, 5)
	})

	Convey("GET /groups?id=10", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?id=10")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(result[0].Name, ShouldNotBeBlank)
	})
	Convey("GET /groups?name=operators", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?name=operators")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 1)
		So(result[0].Name, ShouldEqual, "operators")
	})
	//error checks
	Convey("GET /groups?id=err", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?id=1005001")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 0)
	})
	Convey("GET /groups?name=not-existing-name", t, func() {
		result := []ctx.Group{}
		resp, err := suite.rc.R().SetResult(&result).Get("/groups?name=not-existing-name")
		So(err, ShouldBeNil)
		So(resp.StatusCode(), ShouldEqual, 200)
		So(len(result), ShouldEqual, 0)
	})

}
*/
