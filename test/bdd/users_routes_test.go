package bddtests_test

import (
	"math/rand"
	"net/http"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/corvinusz/echo-xorm/server/users"
)

var _ = Describe("Test GET /users", func() {
	Context("Get all users", func() {
		It("should respond properly", func() {
			var orig, result []users.User
			// get orig
			err := suite.app.C.Orm.Omit("password").Find(&orig)
			Expect(err).NotTo(HaveOccurred())
			// get resp
			resp, err := suite.rc.R().SetResult(&result).Get("/users")
			Expect(err).NotTo(HaveOccurred())
			Expect(http.StatusOK).To(Equal(resp.StatusCode()))
			Expect(len(orig)).To(BeNumerically(">=", 5))
			Expect(len(result)).To(Equal(len(orig)))
			Expect(result).To(BeEquivalentTo(orig))
		})
	})
})

var _ = Describe("Test GET /users/:id", func() {
	Context("with 3 random id", func() {
		It("should respond properly", func() {
			for i := 0; i < 3; i++ {
				id := rand.Int()%7 + 1
				orig := new(users.User)
				result := new(users.User)
				// get orig
				found, err := suite.app.C.Orm.ID(id).Omit("password").Get(orig)
				Expect(err).NotTo(HaveOccurred())
				Expect(found).To(BeTrue())
				// get resp
				resp, err := suite.rc.R().SetResult(result).Get("/users/" + strconv.Itoa(id))
				Expect(err).NotTo(HaveOccurred())
				Expect(http.StatusOK).To(Equal(resp.StatusCode()))
				Expect(result).To(BeEquivalentTo(orig))
			}
		})
	})
})

var _ = Describe("Test POST /users", func() {
	Context("Post predefined user", func() {
		It("should respond properly", func() {
			result := new(users.User)
			passUrl := "a_test_user/password/url"
			payload := users.PostBody{
				Email:       "a_test_user_01_email",
				DisplayName: "a_test_user_01_name",
				Password:    "a_test_user_01_password",
				PasswordURL: &passUrl,
			}
			// http request
			resp, err := suite.rc.R().SetBody(payload).SetResult(result).Post("/users")
			Expect(err).NotTo(HaveOccurred())
			Expect(http.StatusCreated).To(Equal(resp.StatusCode()))
			Expect(result.ID).NotTo(BeZero())
			Expect(result.Email).To(Equal(payload.Email))
			Expect(result.DisplayName).To(Equal(payload.DisplayName))
			Expect(result.Created).NotTo(BeZero())
			Expect(result.Updated).NotTo(BeZero())
			// get original user
			fromDb := new(users.User)
			found, err := suite.app.C.Orm.ID(result.ID).Omit("password").Get(fromDb)
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(result).To(BeEquivalentTo(fromDb))
		})
	})
})
