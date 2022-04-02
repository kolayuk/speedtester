package speedtester_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestSpeedTester(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "speedtester package tests")
}
