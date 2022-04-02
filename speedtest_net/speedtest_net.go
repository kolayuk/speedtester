package speedtest_net

// I moved different implementations to different packages due to 3rd party imports
// to minimize binary size in case of client-app uses only one kind of providers, it can import only required package

import (
	"github.com/pkg/errors"
	"github.com/showwin/speedtest-go/speedtest"
	"speedtester"
)

// returns speedtest.net implementation initialized with provided options
func NewSpeedTestNetProvider(options ...speedtest.Option) speedtester.SpeedTestProvider {
	return &speedTestNetProvider{client: speedtest.New(options...)}
}

type speedTestNetProvider struct {
	client  *speedtest.Speedtest
	servers speedtest.Servers
}

func (s *speedTestNetProvider) TestDownload() (float64, error) {
	err := s.fetchServers()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	maxDownloadSpeed := float64(0)
	for _, server := range s.servers {
		err = server.DownloadTest(false)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		if server.DLSpeed > maxDownloadSpeed {
			maxDownloadSpeed = server.DLSpeed
		}
	}
	return maxDownloadSpeed, nil
}
func (s *speedTestNetProvider) TestDownloadAsync(callback speedtester.MeasuredSpeedCallback) error {
	return errors.Errorf("Not implemented yet")
}

func (s *speedTestNetProvider) TestUpload() (float64, error) {
	err := s.fetchServers()
	if err != nil {
		return 0, errors.WithStack(err)
	}
	maxUploadSpeed := float64(0)
	for _, server := range s.servers {
		err = server.UploadTest(false)
		if err != nil {
			return 0, errors.WithStack(err)
		}
		if server.ULSpeed > maxUploadSpeed {
			maxUploadSpeed = server.ULSpeed
		}
	}
	return maxUploadSpeed, nil
}
func (s *speedTestNetProvider) TestUploadAsync(callback speedtester.MeasuredSpeedCallback) error {
	//TODO implement me
	panic("implement me")
}

func (s *speedTestNetProvider) fetchServers() error {
	user, err := s.client.FetchUserInfo()
	if err != nil {
		return errors.WithStack(err)
	}

	serverList, err := speedtest.FetchServers(user)
	if err != nil {
		return errors.WithStack(err)
	}
	targets, err := serverList.FindServer([]int{})
	if err != nil {
		return errors.WithStack(err)
	}
	s.servers = targets
	return nil
}
