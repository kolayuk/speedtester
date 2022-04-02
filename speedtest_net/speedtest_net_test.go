package speedtest_net

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// tests are the same like in fast_com package (logic is the same but different implementations)
var _ = Describe("speedtest.net support", func() {
	It("async callback should be called for each measurement", func() {
		provider := NewSpeedTestNetProvider()
		Expect(provider.fetchServers()).NotTo(HaveOccurred())
		// we have servers to connect to download data
		Expect(len(provider.servers)).To(BeNumerically(">", 0))
		callBackCalls := 0
		Expect(provider.TestDownloadAsync(func(mbpsSpeed float64) {
			fmt.Println(mbpsSpeed)
			Expect(mbpsSpeed).To(BeNumerically(">", 0)) // speed tested, value positive
			callBackCalls += 1
		})).NotTo(HaveOccurred())
		Expect(callBackCalls).To(Equal(len(provider.servers))) // called for each server, so count should be equal

		// the same for upload
		callBackCalls = 0
		Expect(provider.TestUploadAsync(func(mbpsSpeed float64) {
			Expect(mbpsSpeed).To(BeNumerically(">", 0)) // speed tested, value positive
			callBackCalls += 1
		})).NotTo(HaveOccurred())
		Expect(callBackCalls).To(Equal(len(provider.servers))) // called for each server, so count should be equal
	})
	It("syncronous functions returns maximum speed of all measurements in case of download", func() {
		provider := NewSpeedTestNetProvider()
		maxDownloadSpeed := float64(0)
		Expect(provider.TestDownloadAsync(func(mbpsSpeed float64) {
			if mbpsSpeed > maxDownloadSpeed {
				maxDownloadSpeed = mbpsSpeed
			}
		})).NotTo(HaveOccurred())
		// get value syncronous way
		downloadSpeed, err := provider.TestDownload()
		Expect(err).NotTo(HaveOccurred())
		// I have no idea how to compare internet speed it depends of environment
		// like ISP, country, etc, and event these two value will never be equal
		// so checking them with 20% accuracy
		const ACCURACY = 0.2
		// downloadSpeed +- 20% of maxDownloadSpeed calculated before
		Expect(downloadSpeed).To(BeNumerically(">", maxDownloadSpeed-(maxDownloadSpeed*ACCURACY)))
		Expect(downloadSpeed).To(BeNumerically("<", maxDownloadSpeed+(maxDownloadSpeed*ACCURACY)))
	})
	// same as previous but for upload speed
	It("syncronous functions returns maximum speed of al measurements in case of upload measurement", func() {
		provider := NewSpeedTestNetProvider()
		maxUploadSpeed := float64(0)
		Expect(provider.TestUploadAsync(func(mbpsSpeed float64) {
			if mbpsSpeed > maxUploadSpeed {
				maxUploadSpeed = mbpsSpeed
			}
		})).NotTo(HaveOccurred())
		uploadSpeed, err := provider.TestUpload()
		Expect(err).NotTo(HaveOccurred())
		const ACCURACY = 0.3
		// uploadSpeed +- 30% of maxUploadSpeed calculated before
		Expect(uploadSpeed).To(BeNumerically(">", maxUploadSpeed-(maxUploadSpeed*ACCURACY)))
		Expect(uploadSpeed).To(BeNumerically("<", maxUploadSpeed+(maxUploadSpeed*ACCURACY)))
	})

	// TODO: negative cases ( return error) is not covered, but I have no idea how can I get network error
	// without any hardware breaking (like deattaching an ethernet cable), but in case of that
	// test consistency would be low, positive cases will fail too
})
