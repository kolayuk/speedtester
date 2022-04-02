package fast_com

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("fast.com support", func() {
	It("async callback should be called for each measurement", func() {
		provider := NewFastComProvider()
		Expect(provider.fastComClient.Init()).NotTo(HaveOccurred())
		urls, err := provider.fastComClient.GetUrls()
		Expect(err).NotTo(HaveOccurred())
		// we have urls to download data
		Expect(len(urls)).To(BeNumerically(">", 0))
		callBackCalls := 0
		Expect(provider.TestDownloadAsync(func(mbpsSpeed float64) {
			fmt.Println(mbpsSpeed)
			Expect(mbpsSpeed).To(BeNumerically(">", 0)) // speed tested, value positive
			callBackCalls += 1
		})).NotTo(HaveOccurred())
		// called for each measurement (amount depends on ISP speed)
		// for 10s timeout (hardcoded in 3rd party library)
		Expect(callBackCalls).To(BeNumerically(">", 0))

		// the same for upload
		callBackCalls = 0
		Expect(provider.TestUploadAsync(func(mbpsSpeed float64) {
			Expect(mbpsSpeed).To(BeNumerically(">", 0)) // speed tested, value positive
			callBackCalls += 1
		})).NotTo(HaveOccurred())
		Expect(callBackCalls).To(BeNumerically(">", 0)) // called for each measurement
	})
	It("syncronous functions returns maximum speed of al measurements in case of download", func() {
		provider := NewFastComProvider()
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
		provider := NewFastComProvider()
		maxUploadSpeed := float64(0)
		Expect(provider.TestUploadAsync(func(mbpsSpeed float64) {
			if mbpsSpeed > maxUploadSpeed {
				maxUploadSpeed = mbpsSpeed
			}
		})).NotTo(HaveOccurred())
		uploadSpeed, err := provider.TestUpload()
		Expect(err).NotTo(HaveOccurred())
		const ACCURACY = 0.2
		// uploadSpeed +- 20% of maxUploadSpeed calculated before
		Expect(uploadSpeed).To(BeNumerically(">", maxUploadSpeed-(maxUploadSpeed*ACCURACY)))
		Expect(uploadSpeed).To(BeNumerically("<", maxUploadSpeed+(maxUploadSpeed*ACCURACY)))
	})

	// TODO: negative cases ( return error) is not covered, but I have no idea how can I get network error
	// without any hardware breaking (like deattaching an ethernet cable), but in case of that
	// test consistency would be low, positive cases will fail too
})
