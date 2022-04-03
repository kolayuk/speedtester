package fast_com

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

type mockedFastComClient struct {
	returnErrorOnInit   bool
	targetDownloadSpeed []float64 // kbps
}

func (m *mockedFastComClient) Init() error {
	if m.returnErrorOnInit {
		return errors.Errorf("Mocked client setup to return error")
	}
	return nil
}

func (m *mockedFastComClient) GetUrls() (urls []string, err error) {
	return []string{"mockedUrl"}, nil
}

func (m *mockedFastComClient) Measure(urls []string, KbpsChan chan<- float64) (err error) {
	Expect(len(urls)).To(BeNumerically(">", 0))
	for _, speed := range m.targetDownloadSpeed {
		KbpsChan <- speed
	}
	return nil
}

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
	It("syncronous functions returns maximum speed of all measurements in case of download", func() {
		targetSpeeds := []float64{
			float64(512678.784), // kbps, so 512.678784 mbps
			float64(487678.784),
			float64(563678.784),
		}
		provider := fastComProvider{fastComClient: &mockedFastComClient{returnErrorOnInit: false, targetDownloadSpeed: targetSpeeds}}
		maxDownloadSpeed := float64(0)
		Expect(provider.TestDownloadAsync(func(mbpsSpeed float64) {
			if mbpsSpeed > maxDownloadSpeed {
				maxDownloadSpeed = mbpsSpeed
			}
		})).NotTo(HaveOccurred())
		// get value syncronous way
		downloadSpeed, err := provider.TestDownload()
		Expect(err).NotTo(HaveOccurred())
		Expect(downloadSpeed).To(Equal(float64(563.678784)))
		Expect(maxDownloadSpeed).To(Equal(float64(563.678784)))
	})
	// same as previous but for upload speed (but based on real speed, so we're getting 30% accuracy)
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
		const ACCURACY = 0.3
		// uploadSpeed +- 30% of maxUploadSpeed calculated before
		Expect(uploadSpeed).To(BeNumerically(">", maxUploadSpeed-(maxUploadSpeed*ACCURACY)))
		Expect(uploadSpeed).To(BeNumerically("<", maxUploadSpeed+(maxUploadSpeed*ACCURACY)))
	})
	It("testing returning error baseds on the mock implementation", func() {
		provider := fastComProvider{fastComClient: &mockedFastComClient{returnErrorOnInit: true}}
		_, err := provider.TestDownload()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Mocked client setup to return error"))
	})
})
