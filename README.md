# go-speedtester

Small golang library to measure your internet connection speedwith one of two (built-in) services:

* fast.com (Netflix's)
* speedtest.net (Ookla's)

## Installation

```
go get github.com/kolayuk/speedtester
```

## Usage

Library exports one public function and two implementations for different services

```go
// using with fast.com implementation
downloadSpeed, uploaduploadSpeed, err := speedtester.TestSpeed(fast_com.NewFastComProvider())
if err != nil {
panic(err)
}
fmt.Println("fast.com", downloadSpeed, uploadSpeed)

// using with speedtest.net implementation
downloadSpeed, uploadSpeed, err := speedtester.TestSpeed(speedtest_net.NewSpeedTestNetProvider())
if err != nil {
panic(err)
}
fmt.Println("speedtest.net", downloadSpeed, uploadSpeed)
```

You can also add new speedtest providers of your choice without modification of the library. You will need to implement
an [interface]() and use your implementation in TestSpeed function call

```go
package my_awesome_package

import "github.com/kolayuk/speedtester"

// your custom implemetation
type myAwesomeSpeedTestService struct {
}

func (m *myAwesomeSpeedTestService) TestDownload() (float64, error) {
	panic("implement me")
}

func (m *myAwesomeSpeedTestService) TestDownloadAsync(callback speedtester.MeasuredSpeedCallback) error {
	panic("implement me")
}

func (m *myAwesomeSpeedTestService) TestUpload() (float64, error) {
	panic("implement me")
}

func (m *myAwesomeSpeedTestService) TestUploadAsync(callback speedtester.MeasuredSpeedCallback) error {
	panic("implement me")
}

// measurement using your custom implementation
downloadSpeed, uploadSpeed, err := speedtester.TestSpeed(&myAwesomeSpeedTestService{})
if err != nil {
panic(err)
}
fmt.Println("myAwesomeSpeedTestService", downloadSpeed, uploadSpeed)
```