package speedtester

import "github.com/pkg/errors"

// Tests internet connection speed, accepts implementation of test speed provider (for the ability to create external 3rd party providers outside of our package)
// returns download speed (mbps), upload speed(mbps) and error if occured
// TODO: it would be good to separate download & upload tests but we have a requirment to have 1 external api for both
func TestSpeed(provider SpeedTestProvider) (downloadSpeed float64, uploadSpeed float64, err error) {
	// test max download speed
	downloadSpeed, err = provider.TestDownload()
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	// test max upload speed
	uploadSpeed, err = provider.TestUpload()
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	return downloadSpeed, uploadSpeed, nil
}

// function signature will be called on measured speed
type MeasuredSpeedCallback func(mbpsSpeed float64)

// Interface to implement for different speed tests providers (speedtest.com, fast.com, etc...)
type SpeedTestProvider interface {
	// returns max download speed in mbps of all tries
	TestDownload() (float64, error)
	// callback will be called for every measurement
	TestDownloadAsync(callback MeasuredSpeedCallback) error
	// returns max upload speed in mbps of all tries
	TestUpload() (float64, error)
	// callback will be called for every measurement
	TestUploadAsync(callback MeasuredSpeedCallback) error
}
