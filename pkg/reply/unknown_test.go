package reply

import (
	"testing"

	"github.com/golang/mock/gomock"
	core "github.com/jenkins-zh/wechat-backend/pkg"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUnknown(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unknown keywords")
}

var _ = Describe("Unknon keywords", func() {
	var (
		reply AutoReply
		ctrl  *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		reply = &UnknownAutoReply{}
	})

	It("should not error", func() {
		Expect(reply.Accept(&core.TextRequestBody{})).To(Equal(true))
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
