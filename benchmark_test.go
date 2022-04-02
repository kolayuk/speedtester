package speedtester_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"speedtester"
	"speedtester/fast_com"
)

var _ = Describe("Benchmarking", func() {
	BeforeEach(func() {

	})
	const COUNTS = 2

	Measure("FastCom implementation", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
			Expect(err).NotTo(HaveOccurred())
			Expect(download).To(BeNumerically(">", 0))
			Expect(upload).To(BeNumerically(">", 0))
		})
		// not sure if is it critical to long execution
		Expect(runtime.Seconds()).Should(BeNumerically("<", 60))
	}, COUNTS)
	Measure("SpeedTestNet implementation", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
			Expect(err).NotTo(HaveOccurred())
			Expect(download).To(BeNumerically(">", 0))
			Expect(upload).To(BeNumerically(">", 0))
		})

		Expect(runtime.Seconds()).Should(BeNumerically("<", 60), "PutInsuranceObject should be fast")
	}, COUNTS)
})
