package speedtester_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestInvoicesRouter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "invoices router tests")
}
