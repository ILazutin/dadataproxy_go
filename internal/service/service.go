package service

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"reflect"

	"github.com/ilazutin/dadataproxy_go/internal/service/dadata"
	"github.com/ilazutin/dadataproxy_go/internal/storage"
)

type ProxyService struct {
	dadata  *dadata.DaData
	storage storage.Storage
	logger  *slog.Logger
}

func New(dadata *dadata.DaData, storage storage.Storage, logger *slog.Logger) *ProxyService {
	return &ProxyService{
		dadata:  dadata,
		storage: storage,
		logger:  logger,
	}
}

func (s ProxyService) CleanValue(path string, data string, ignoreCache bool, logger *slog.Logger) (*dadata.DaDataClean, error) {

	queryString, err := s.formatQuery(data)
	if err != nil {
		return nil, err
	}

	storageKey := s.makeHash(path, queryString)

	if !ignoreCache {
		storagedValue, err := s.storage.Read(storageKey)
		if err != nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
				slog.String("error", err.Error()),
			)
		}
		if storagedValue == nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
			)
		} else {
			var res *dadata.DaDataClean
			err := json.Unmarshal([]byte(storagedValue.(string)), &res)
			if err != nil {
				logger.Error("Decode cached value", slog.String("error", err.Error()))
			}

			logger.Info("Get from cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey))

			return res, nil
		}
	}

	result, err := s.dadata.GetCleanValue(path, queryString)
	if err != nil {
		logger.Error("Get value from dadata api", slog.String("query", queryString))
		return nil, err
	}

	encodingResult, _ := json.Marshal(result)
	go s.storage.Save(storageKey, encodingResult)

	return result, nil
}

func (s ProxyService) SuggestValue(path string, data string, ignoreCache bool, logger *slog.Logger) (*dadata.DaDataSuggest, error) {

	queryString, err := s.formatQuery(data)
	if err != nil {
		return nil, err
	}

	storageKey := s.makeHash(path, queryString)

	if !ignoreCache {
		storagedValue, err := s.storage.Read(storageKey)
		if err != nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
				slog.String("error", err.Error()),
			)
		}
		if storagedValue == nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
			)
		} else {
			var res *dadata.DaDataSuggest
			err := json.Unmarshal([]byte(storagedValue.(string)), &res)
			if err != nil {
				logger.Error("Decode cached value", slog.String("error", err.Error()))
			}

			logger.Info("Get from cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
			)

			return res, nil
		}
	}

	result, err := s.dadata.GetSuggestValue(path, queryString)
	if err != nil {
		logger.Error("Get value from dadata api", slog.String("query", queryString))
		return nil, err
	}

	encodingResult, _ := json.Marshal(result)
	go s.storage.Save(storageKey, encodingResult)

	return result, nil
}

func (s ProxyService) IpLocateValue(path string, data string, ignoreCache bool, logger *slog.Logger) (*dadata.DaDataIpLocate, error) {

	queryString, err := s.formatQuery(data)
	if err != nil {
		return nil, err
	}

	storageKey := s.makeHash(path, queryString)

	if !ignoreCache {
		storagedValue, err := s.storage.Read(storageKey)
		if err != nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
				slog.String("error", err.Error()),
			)
		}
		if storagedValue == nil {
			logger.Info("Key not found in cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
			)
		} else {
			var res *dadata.DaDataIpLocate
			err := json.Unmarshal([]byte(storagedValue.(string)), &res)
			if err != nil {
				logger.Error("Decode cached value", slog.String("error", err.Error()))
			}

			logger.Info("Get from cache",
				slog.String("path", path),
				slog.String("query", queryString),
				slog.String("key", storageKey),
			)

			return res, nil
		}
	}

	result, err := s.dadata.GetIpLocateValue(path, queryString)
	if err != nil {
		logger.Error("Get value from dadata api", slog.String("query", queryString))
		return nil, err
	}

	encodingResult, _ := json.Marshal(result)
	go s.storage.Save(storageKey, encodingResult)

	return result, nil
}

func (s ProxyService) SaveToCache(path string, query interface{}, body interface{}, logger *slog.Logger) error {

	queryString, err := json.Marshal(query)
	if err != nil {
		return err
	}

	storageKey := s.makeHash(path, string(queryString))

	var encodingResult []byte
	if reflect.TypeOf(body).Kind() == reflect.String {
		encodingResult = []byte(body.(string))
	} else {
		encodingResult, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	err = s.storage.Save(storageKey, encodingResult)
	logger.Info("Key save to cache",
		slog.String("key", storageKey),
		slog.String("path", path),
		slog.String("query", string(queryString)),
		slog.String("body", string(encodingResult)[0:20]+"..."),
	)

	return err
}

func (s ProxyService) MigrateCache(logger *slog.Logger) error {

	keys, err := s.storage.ReadAllKeys()
	if err != nil {
		return err
	}

	const path = "/clean/address"

	for _, key := range keys {
		value, err := s.storage.Read(key)
		if err != nil {
			logger.Error("Cannot convert key's value",
				slog.String("key", key),
				slog.String("value", value.(string)),
				slog.String("error", err.Error()),
			)
			continue
		}

		var structValue dadata.DaDataClean
		err = json.Unmarshal([]byte(value.(string)), &structValue)
		if err != nil {
			logger.Error("Cannot convert key's value",
				slog.String("key", key),
				slog.String("value", value.(string)),
				slog.String("error", err.Error()),
			)
			continue
		}

		if len(structValue) < 1 {
			logger.Error("Key is empty",
				slog.String("key", key),
				slog.String("value", value.(string)),
			)
			continue
		}

		firstValue := structValue[0]

		s.SaveToCache(path, []string{firstValue["source"].(string)}, value, logger)
	}

	return nil
}

func (s ProxyService) makeHash(values ...string) string {
	h := sha256.New()
	for _, value := range values {
		io.WriteString(h, value)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s ProxyService) formatQuery(data string) (string, error) {
	var query interface{}
	err := json.Unmarshal([]byte(data), &query)
	if err != nil {
		return "", err
	}

	queryString, err := json.Marshal(query)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	err = json.Compact(buffer, queryString)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}
