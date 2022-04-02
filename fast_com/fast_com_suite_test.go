package fast_com

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFastCom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "FastCom Suite")
}
