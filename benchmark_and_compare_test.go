package speedtester_test

import (
	"fmt"
	"github.com/kolayuk/speedtester"
	"github.com/kolayuk/speedtester/fast_com"
	"github.com/kolayuk/speedtester/speedtest_net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Benchmarking", func() {
	BeforeEach(func() {

	})
	const COUNTS = 2
	// measure timings for fast.com
	Measure("fast.com implementation", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
			Expect(err).NotTo(HaveOccurred())          // no error
			Expect(download).To(BeNumerically(">", 0)) // speed was calculated
			Expect(upload).To(BeNumerically(">", 0))
			fmt.Println("fast.com", download, upload)
		})
		// not sure if is it critical to long execution
		Expect(runtime.Seconds()).Should(BeNumerically("<", 60))
	}, COUNTS)

	// measure timings for speedtest.net
	Measure("speedtest.net implementation", func(b Benchmarker) {
		_ = b.Time("runtime", func() {
			download, upload, err := speedtester.TestSpeed(speedtest_net.NewSpeedTestNetProvider())
			Expect(err).NotTo(HaveOccurred())
			Expect(download).To(BeNumerically(">", 0))
			Expect(upload).To(BeNumerically(">", 0))
			fmt.Println("speedtest.net", download, upload)
		})
	}, COUNTS)

	// difference between fast.com and speedtest results shouldn't be too much (const = 30%)
	It("compare results between speedtest.net and fast.com", func() {
		const ACCEPTABLE_DIFFERENCE = 0.3 // 30%
		downloadFastCom, uploadFastCom, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
		Expect(err).NotTo(HaveOccurred())
		downloadSpeedTest, uploadSpeedTest, err := speedtester.TestSpeed(speedtest_net.NewSpeedTestNetProvider())
		Expect(err).NotTo(HaveOccurred())

		// fastCom - 30% < speedtest < fastcom+30% for download
		Expect(downloadSpeedTest).To(BeNumerically(">", downloadFastCom-(downloadFastCom*ACCEPTABLE_DIFFERENCE)))
		Expect(downloadSpeedTest).To(BeNumerically("<", downloadFastCom+(downloadFastCom*ACCEPTABLE_DIFFERENCE)))
		// the same for upload
		Expect(uploadSpeedTest).To(BeNumerically(">", uploadFastCom-(uploadFastCom*ACCEPTABLE_DIFFERENCE)))
		Expect(uploadSpeedTest).To(BeNumerically("<", uploadFastCom+(uploadFastCom*ACCEPTABLE_DIFFERENCE)))
	})
	// TODO: negative cases ( return error) is not covered, but I have no idea how can I get network error
	// without any hardware breaking (like deattaching an ethernet cable), but in case of that
	// test consistency would be low, positive cases will fail too
})
