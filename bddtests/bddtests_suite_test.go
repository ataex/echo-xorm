package bddtests_test

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/resty.v1"
	"gopkg.in/testfixtures.v2"

	"github.com/corvinusz/echo-xorm/app"
	"github.com/corvinusz/echo-xorm/ctx"
	"github.com/corvinusz/echo-xorm/server/auth"
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

const (
	cfgFileName    = "./test-config/echo-xorm-test-config.toml"
	fixturesFolder = "./fixtures"
)

// LsxTestSuite is testing context for app
type LsxTestSuite struct {
	app     *app.Application
	baseURL string
	rc      *resty.Client
}

// SetupTest called once before test
func (s *LsxTestSuite) setupSuite() error {
	err := s.setupServer()
	if err != nil {
		return err
	}
	s.baseURL = "http://localhost:" + s.app.C.Config.Port
	// create and setup resty client
	s.rc = resty.DefaultClient
	s.rc.SetHeader("Content-Type", "application/json")
	s.rc.SetHostURL(s.baseURL)
	// get auth token
	s.authorizeMe("admin", "admin")
	return nil
}

//------------------------------------------------------------------------------
func (s *LsxTestSuite) setupServer() error {
	var err error
	// init test application
	s.app, err = app.New(&ctx.Flags{CfgFileName: cfgFileName})
	if err != nil {
		return err
	}
	// load fixtures
	err = s.setupFixtures()
	if err != nil {
		return err
	}
	// start test server with go routine
	go s.app.Run()
	// wait til server started then return
	return s.waitServerStart(3 * time.Second)
}

//------------------------------------------------------------------------------
func (s *LsxTestSuite) setupFixtures() error {
	return testfixtures.LoadFixtures(
		fixturesFolder,
		s.app.C.Orm.DB().DB,
		&testfixtures.SQLite{})
}

//------------------------------------------------------------------------------
func (s *LsxTestSuite) waitServerStart(timeout time.Duration) error {
	const sleepTime = 300 * time.Millisecond
	dialer := &net.Dialer{
		DualStack: false,
		Deadline:  time.Now().Add(timeout),
		Timeout:   sleepTime,
		KeepAlive: 0,
	}
	done := time.Now().Add(timeout)
	for time.Now().Before(done) {
		c, err := dialer.Dial("tcp", ":"+s.app.C.Config.Port)
		if err == nil {
			return c.Close()
		}
		time.Sleep(sleepTime)
	}
	return fmt.Errorf("cannot connect %v for %v", s.baseURL, timeout)
}

//------------------------------------------------------------------------------
func (s *LsxTestSuite) authorizeMe(login, password string) error {
	// make authorization
	payload := auth.Input{
		Login:    login,
		Password: password,
	}
	result := new(auth.Result)
	response, err := s.rc.R().SetBody(payload).SetResult(result).Post("/auth")
	if err != nil {
		return err
	}

	// check response and set token
	if response.StatusCode() != 200 {
		return errors.New("auth response status is not 200 (not OK)")
	}
	// set auth token
	s.rc.SetAuthToken(result.Token)
	// return
	return nil
}
