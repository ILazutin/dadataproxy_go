package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/ilazutin/dadataproxy_go/internal/service/dadata"
	"github.com/stretchr/testify/suite"
)

type ProxyServiceTestSuite struct {
	suite.Suite

	service *ProxyService
}

type MockStorage struct{}

func (st MockStorage) Save(string, interface{}) error {
	return nil
}

func (st MockStorage) Read(string) (interface{}, error) {
	return nil, nil
}

func (st MockStorage) ReadAllKeys() ([]string, error) {
	return nil, nil
}

func (suite *ProxyServiceTestSuite) SetupTest() {
	storageMock := MockStorage{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	suite.service = New(dadata.New("ComeAndGetYourOwnSecretKey", "ComeAndGetYourOwnSecretKey"), storageMock, logger)
}

func (suite *ProxyServiceTestSuite) TestKeepConnections() {
	addresses := make([]string, 102)
	for index := range addresses {
		addresses[index] = fmt.Sprintf("[ \"мск сухонская %d%d/-%d\" ]", index+1, index, index+2)
	}

	for index, address := range addresses {
		result, err := suite.service.CleanValue("clean/address", address, false, suite.service.logger)
		suite.Require().NoError(err)
		body, err := json.Marshal(result)
		if err == nil {
			suite.service.logger.Info(fmt.Sprintf("Result %d: %s", index, string(body)))
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func TestServiceTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ProxyServiceTestSuite))
}
