package fast_com

// I moved different implementations to different packages due to 3rd party imports
// to minimize binary size in case of client-app uses only one kind of providers, it can import only required package

import (
	"github.com/ddo/go-fast"
	"github.com/pkg/errors"
	"net/http"
	"speedtester"
	"strings"
	"sync"
	"time"
)

// returns fast.com implementation of speed test
func NewFastComProvider() speedtester.SpeedTestProvider {
	return &fastComProvider{}
}

type fastComProvider struct {
}

func (f *fastComProvider) TestDownload() (float64, error) {
	maxMeasuredSpeed := float64(0)
	err := f.TestDownloadAsync(func(mbpsSpeed float64) {
		// saving max speed in our callback to return it later
		if mbpsSpeed > maxMeasuredSpeed {
			maxMeasuredSpeed = mbpsSpeed
		}
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return maxMeasuredSpeed, nil
}

func (f *fastComProvider) TestDownloadAsync(callback speedtester.MeasuredSpeedCallback) error {
	// we are free to use 3rd party libraries by the task description, so we have a library with a download speed test from fast.com
	// using it here, so dont reinvent wheel
	fastCom := fast.New()
	err := fastCom.Init()
	if err != nil {
		return errors.WithStack(err)
	}

	// get urls to download from
	urls, err := fastCom.GetUrls()
	if err != nil {
		return errors.WithStack(err)
	}

	// measure speed
	KbpsChan := make(chan float64)
	go func() {
		for Kbps := range KbpsChan {
			// call
			callback(Kbps / 1000)
		}
	}()

	err = fastCom.Measure(urls, KbpsChan)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (f *fastComProvider) TestUpload() (float64, error) {
	maxUploadSpeed := float64(0)
	err := f.TestUploadAsync(func(mbpsSpeed float64) {
		// saving max speed
		if mbpsSpeed > maxUploadSpeed {
			maxUploadSpeed = mbpsSpeed
		}
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return maxUploadSpeed, nil
}
func (f *fastComProvider) TestUploadAsync(callback speedtester.MeasuredSpeedCallback) error {
	// get URLs the same way as in download test
	fastCom := fast.New()
	err := fastCom.Init()
	if err != nil {
		return errors.WithStack(err)
	}
	// get urls to download from
	urls, err := fastCom.GetUrls()
	if err != nil {
		return errors.WithStack(err)
	}
	// but 3rd party library we found it seems does not support upload speed test
	// so reinventing wheel here
	var wg sync.WaitGroup
	wg.Add(len(urls))
	errChan := make(chan error)
	for _, url := range urls {
		go func(urlToSendPostRequest string) {
			const megabyte = 1000 * 1000      // TODO: 1024? but it seems we're calculating network throughput, not disk space, so it seems ok
			const payloadSize = 10 * megabyte // TODO: not sure about size, is 10 mb okay?
			timeStart := time.Now()
			// according to dev tools on fast.com frontend sends post requests to urls from the lists and writes data, payload does not matter
			payload := strings.NewReader(strings.Repeat("0", payloadSize))
			resp, err := http.Post(urlToSendPostRequest, "text/plain", payload)
			if err != nil {
				errChan <- err
				wg.Done()
				return
			}
			defer func() {
				_ = resp.Body.Close()
				wg.Done()
			}()
			// cheching how much time is takes to send such amount of data
			spentTime := time.Since(timeStart).Seconds()
			bitsPerSec := (payloadSize * float64(8)) / spentTime // *8 to get bits
			callback((bitsPerSec / megabyte))                    // divide bits/s to megabyte (1000*1000) to get mbps
		}(url)
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return errors.WithStack(err)
	default:
		return nil
	}

}
