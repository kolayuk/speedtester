package speedtester

import "github.com/pkg/errors"

// Tests internet connection speed, accepts implementation of test speed provider (for the ability to create extrernal 3rd party providers)
// returns download speed (mbps), upload speed(mbps) and error if occured
func TestSpeed(provider SpeedTestProvider) (downloadSpeed float64, uploadSpeed float64, err error) {
	return 0, 0, errors.Errorf("not implemented yet")
}

// Interface to implement for different speed tests providers (speedtest.com, fast.com, etc...)
type SpeedTestProvider interface {
	TestDownload() (float64, error)
	TestUpload() (float64, error)
}
