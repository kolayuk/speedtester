package speedtester_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"speedtester"
	"speedtester/fast_com"
	"speedtester/speedtest_net"
)

var _ = Describe("Benchmarking", func() {
	BeforeEach(func() {

	})
	const COUNTS = 2

	Measure("fast.com implementation", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
			Expect(err).NotTo(HaveOccurred())
			Expect(download).To(BeNumerically(">", 0))
			Expect(upload).To(BeNumerically(">", 0))
			fmt.Println("fast.com", download, upload)
		})
		// not sure if is it critical to long execution
		Expect(runtime.Seconds()).Should(BeNumerically("<", 60))
	}, COUNTS)
	Measure("speedtest.net implementation", func(b Benchmarker) {
		_ = b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(speedtest_net.NewSpeedTestNetProvider())
			Expect(err).NotTo(HaveOccurred())
			Expect(download).To(BeNumerically(">", 0))
			Expect(upload).To(BeNumerically(">", 0))
			fmt.Println("speedtest.net", download, upload)
		})
	}, COUNTS)
})
