package fast_com

import "speedtester"

// requires api token from fast.com
func NewFastComProvider(token string) speedtester.SpeedTestProvider {
	return &fastComProvider{apiToken: token}
}

type fastComProvider struct {
	apiToken string
}

func (f *fastComProvider) TestDownload() (float64, error) {
	//TODO implement me
	panic("implement me")
}

func (f *fastComProvider) TestUpload() (float64, error) {
	//TODO implement me
	panic("implement me")
}
