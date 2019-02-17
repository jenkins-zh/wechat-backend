package reply

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGitHubBind(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unknown keywords")
}

var _ = Describe("github bind", func() {
	var (
		binder GitHubBind
		data   GitHubBindData
		ctrl   *gomock.Controller
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		binder = &GitHubBinder{
			File: "github_bind.yaml",
		}
		data = GitHubBindData{
			WeChatID: "WeChatID",
			GitHubID: "GitHubID",
		}
	})

	It("should not error", func() {
		Expect(binder.Exists("none")).To(Equal(false))

		Expect(binder.Add(data)).To(BeNil())

		Expect(binder.Exists("WeChatID")).To(Equal(true))

		Expect(binder.Find("WeChatID").GitHubID).To(Equal("GitHubID"))
	})

	It("non-repetitive", func() {
		Expect(binder.Add(data)).To(BeNil())

		Expect(binder.Count()).To(Equal(1))
	})

	AfterEach(func() {
		ctrl.Finish()
	})
})
