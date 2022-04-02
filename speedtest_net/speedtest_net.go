package speedtest_net

// I moved different implementations to different packages due to 3rd party imports
// to minimize binary size in case of client-app uses only one kind of providers, it can import only required package

import (
	"github.com/pkg/errors"
	"github.com/showwin/speedtest-go/speedtest"
	"speedtester"
	"sync"
)

// returns speedtest.net implementation initialized with provided options
func NewSpeedTestNetProvider(options ...speedtest.Option) *speedTestNetProvider {
	return &speedTestNetProvider{client: speedtest.New(options...)}
}

type speedTestNetProvider struct {
	client  *speedtest.Speedtest
	servers speedtest.Servers
}

func (s *speedTestNetProvider) TestDownload() (float64, error) {
	maxDownloadSpeed := float64(0)
	err := s.TestDownloadAsync(func(mbpsSpeed float64) {
		if mbpsSpeed > maxDownloadSpeed {
			maxDownloadSpeed = mbpsSpeed
		}
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return maxDownloadSpeed, nil
}
func (s *speedTestNetProvider) TestDownloadAsync(callback speedtester.MeasuredSpeedCallback) error {
	err := s.fetchServers()
	if err != nil {
		return errors.WithStack(err)
	}
	errChan := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(s.servers))
	for _, serverInList := range s.servers {
		go func(server *speedtest.Server) {
			defer wg.Done()
			err = server.DownloadTest(false)
			if err != nil {
				errChan <- err
			}
			callback(server.DLSpeed)
		}(serverInList)
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return errors.WithStack(err)
	default:
		return nil
	}
}

func (s *speedTestNetProvider) TestUpload() (float64, error) {
	maxUploadSpeed := float64(0)
	err := s.TestUploadAsync(func(mbpsSpeed float64) {
		if mbpsSpeed > maxUploadSpeed {
			maxUploadSpeed = mbpsSpeed
		}
	})
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return maxUploadSpeed, nil
}
func (s *speedTestNetProvider) TestUploadAsync(callback speedtester.MeasuredSpeedCallback) error {
	err := s.fetchServers()
	if err != nil {
		return errors.WithStack(err)
	}
	errChan := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(s.servers))
	for _, serverInList := range s.servers {
		go func(server *speedtest.Server) {
			defer wg.Done()
			err = server.UploadTest(false)
			if err != nil {
				errChan <- err
			}
			callback(server.ULSpeed)
		}(serverInList)
	}
	wg.Wait()
	select {
	case err := <-errChan:
		return errors.WithStack(err)
	default:
		return nil
	}
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
