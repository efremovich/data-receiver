package operator_client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Config struct {
	Addr               string        `yaml:"addr" conf:"addr" validate:"required"`
	Timeout            time.Duration `yaml:"timeout" conf:"timeout"`
	InsecureSkipVerify bool          `yaml:"insecure_skip_verify" conf:"insecure_skip_verify"`
	Auth               Auth          `yaml:"auth" conf:"auth"`
}

type Auth struct {
	IsRequired bool   `yaml:"is_required" conf:"is_required"`
	Login      string `yaml:"login" conf:"login"`
	Password   string `yaml:"password" conf:"password"`
}

type Client struct {
	addr string // Адрес оператора в формате https://edo.keydisk.ru
	auth Auth

	httpClient *http.Client
}

func NewClient(conf Config) (*Client, error) {
	u, err := url.Parse(conf.Addr)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга поля Addr из конфигурации: %w", err)
	}

	timeout := conf.Timeout
	if timeout == 0 {
		timeout = time.Second * 5
	}

	httpClient := &http.Client{Timeout: timeout}
	if conf.InsecureSkipVerify {
		httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	return &Client{
		addr:       u.String(),
		auth:       conf.Auth,
		httpClient: httpClient,
	}, nil
}

const (
	fileTypeDocument = "doc"
	fileTypeSign     = "sign"
	fileTypeMeta     = "meta"
)

type OperatorSystemResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}
