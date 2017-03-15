package bddtests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/resty.v0"
	"gopkg.in/testfixtures.v2"

	"fmt"
	"net"
	"testing"
	"time"

	"github.com/corvinusz/echo-xorm/app"
	"github.com/corvinusz/echo-xorm/ctx"
)

func TestBddtests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bddtests Suite")
}

var suite *LsxTestSuite

var _ = BeforeSuite(func() {
	suite = new(LsxTestSuite)
	err := suite.setupSuite()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	if suite.app.C.Orm != nil {
		suite.app.C.Orm.Close()
	}
})

const cfgFileName = "./test-config/echo-xorm-test-config.toml"
const fixturesFolder = "./fixtures"

// LsxTestSuite is testing context for app
type LsxTestSuite struct {
	app     *app.Application
	baseURL string
	rc      *resty.Client
}

// SetupTest called once before test
func (suite *LsxTestSuite) setupSuite() error {
	err := suite.setupServer()
	if err != nil {
		return err
	}
	suite.baseURL = "http://localhost:" + suite.app.C.Config.Port
	// create and setup resty client
	suite.rc = resty.DefaultClient
	suite.rc.SetHeader("Content-Type", "application/json")
	suite.rc.SetHostURL(suite.baseURL)
	return nil
}

// setupServer prepares testing server with test data
func (suite *LsxTestSuite) setupServer() error {
	var err error
	// init test application
	appFlags := &ctx.Flags{
		CfgFileName: cfgFileName,
	}

	suite.app, err = app.New(appFlags)
	if err != nil {
		return err
	}

	// load test fixtures
	err = suite.setupFixtures()
	if err != nil {
		return err
	}

	// start test server
	go suite.app.Run()
	// wait til server started then return
	return suite.waitServerStart(3 * time.Second)
}

// setupFixtures writes test data to database from fixtures
func (suite *LsxTestSuite) setupFixtures() error {
	db := suite.app.C.Orm.DB().DB
	fixtures, err := testfixtures.NewFolder(db, &testfixtures.SQLite{}, fixturesFolder)
	if err != nil {
		return err
	}
	err = fixtures.Load()
	return err
}

// waitServerStart redials server til OK or timeout
func (suite *LsxTestSuite) waitServerStart(timeout time.Duration) error {
	const sleepTime = 300 * time.Millisecond
	dialer := &net.Dialer{
		DualStack: false,
		Deadline:  time.Now().Add(timeout),
		Timeout:   sleepTime,
		KeepAlive: 0,
	}
	done := time.Now().Add(timeout)
	for time.Now().Before(done) {
		c, err := dialer.Dial("tcp", ":"+suite.app.C.Config.Port)
		if err == nil {
			return c.Close()
		}
		fmt.Println(err.Error())
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("cannot connect %v for %v", suite.baseURL, timeout)
}
