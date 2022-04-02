package speedtest_net

import "speedtester"

// returns speedtest.net implementation
func NewSpeedTestNetProvider() speedtester.SpeedTestProvider {
	return &speedTestNetProvider{}
}

type speedTestNetProvider struct {
}

func (s *speedTestNetProvider) TestDownload() (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *speedTestNetProvider) TestUpload() (float64, error) {
	//TODO implement me
	panic("implement me")
}
