package service_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/linuxsuren/wechat-backend/pkg/config"
	"github.com/linuxsuren/wechat-backend/pkg/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"net/http"
	"net/http/httptest"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Test")
}

var _ = Describe("", func() {
	var (
		weConfig *config.WeChatConfig
		ctrl     *gomock.Controller
		request  *http.Request
		writer   http.ResponseWriter
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		weConfig = config.NewConfig()

		request = httptest.NewRequest("GET", "/test?valid=true", nil)
		writer = httptest.NewRecorder()
	})

	It("turn on the valid", func() {
		Expect(weConfig.Valid).To(Equal(false))

		service.HandleConfig(writer, request, weConfig)

		Expect(weConfig.Valid).To(Equal(true))
	})

	Context("turn off the valid", func() {
		JustBeforeEach(func() {
			request = httptest.NewRequest("GET", "/test", nil)
		})

		It("turn off the valid", func() {
			weConfig.Valid = true
			service.HandleConfig(writer, request, weConfig)
			Expect(weConfig.Valid).To(Equal(false))
		})
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
