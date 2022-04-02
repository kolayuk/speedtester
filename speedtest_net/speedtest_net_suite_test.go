package speedtest_net_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSpeedtestNet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SpeedtestNet Suite")
}
