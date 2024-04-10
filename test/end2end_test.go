package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"testing"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/app"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/tprepo"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/storage"
	"github.com/efremovich/data-receiver/pkg/broker"
	anats "github.com/efremovich/data-receiver/pkg/anats"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceivePacket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, pgConnStr, err := postgresdb.GetMockConn("../migrations/package_receiver_db")
	if err != nil {
		t.Fatalf(err.Error())
	}

	connStr, err := anats.CreateNatsTempContainer(ctx)
	if err != nil {
		t.Fatalf(err.Error())
	}

	tpRepo, err := tprepo.NewTransportPackageRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	storagePath := "./storage"
	_ = os.RemoveAll(storagePath)

	err = os.Mkdir(storagePath, 0777)
	if err != nil {
		t.Fatalf(err.Error())
	}

	storageMock, err := storage.NewMockStorageClient(storagePath)
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer func() { _ = os.RemoveAll(storagePath) }()

	// Тестовая конфигурация приложения.
	// ДБ и НАТС подняли в контейнерах. Клиент к оператору сломанный - используем бекап операторов.
	// Сторадж - заглушка в памяти.
	cfg := config.Config{
		ServiceName:    "e2e_test_receiver",
		SelfOperatorID: "2AE",
		NatsURL:        []string{connStr},
		OperatorAPI: config.OperatorAPI{
			BaseURL:  "http//asd", // Укажем сломанный, чтобы использовать бекап из файла.
			Login:    "l",
			Password: "p",
		},
		PGWriterConn: pgConnStr,
		PGReaderConn: pgConnStr,
		LogLevel:     -4,
		Packer: config.Packer{
			CertData: crt,
			KeyData:  key,
		},
		Storage: config.Storage{
			URL:            "./storage", // Для мока - это путь к директории.
			UseMockStorage: true,
		},
		Gateway: config.Gateway{
			AuthToken: "secret",
			HTTP: config.Adr{
				Host: "localhost",
				Port: "12321",
			},
			GRPC: config.Adr{
				Host: "localhost",
				Port: "12322",
			},
		},
	}

	testAPP, err := app.New(ctx, cfg)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Запускаем инстант приложения локально.
	go func() {
		err = testAPP.Start(ctx)
		if err != nil {
			panic(err)
		}
	}()

	// Сделали подписку на NATS, чтобы считать сообщения, которые посылает приёмник.
	var natsMessagesCounter = new(int32)
	*natsMessagesCounter = 0
	err = natsSucsribe(natsMessagesCounter, connStr)

	if err != nil {
		t.Fatalf(err.Error())
	}

	type testCase struct {
		testName string

		fileName string // Имя CMS-файла.

		respCode            int                 // Ожидаемый статус-код.
		tpStatusInDB        entity.TpStatusEnum // Статус записи в БД.
		natsMessagesCounter int                 // Количество сообщений в NATS после теста.
		filesInStorage      []string            // Проверить наличие этих файлов в сторадж
		filesInDB           int                 // Количество записей в таблице с файлами
	}

	testCases := []testCase{
		{
			// CMS c невалидным сертом.
			testName:            "cms_сломанный_серт",
			fileName:            "dadfe07f757a436aad6f77bbebe6c022.cms",
			respCode:            http.StatusBadRequest,
			tpStatusInDB:        entity.TpStatusEnumFailed,
			natsMessagesCounter: 0,
			filesInDB:           0,
		},
		{
			// Валидный CMS с 1 ЛС.
			testName:            "cms_валидный_1_лс",
			fileName:            "d831d274d9f248d791e34bf20943511b.cms",
			respCode:            http.StatusOK,
			tpStatusInDB:        entity.TpStatusEnumSuccess,
			natsMessagesCounter: 1,
			filesInStorage: []string{"d831d274d9f248d791e34bf20943511b.cms", "77755f3ecdda43d9b36d4f3a0b13e29e.p7s", "da9208b807284135ba2dffaaf08bff27.bin", "dmcd00011c5d943b2986b2a845bd1d1a.bin",
				"file2002f7db4e748bced6fc17486501.bin", "mcd00011c5d943b2986b2a845bd1d1af.bin", "pdmcd0001c5d943b2986b2a845bd1d1a.p7s", "pmcd0001c5d943b2986b2a845bd1d1af.p7s"},
			filesInDB: 8,
		},
		{
			// Отправим тот же CMS ещё раз. Должен вернуться тот же ответ, новое сообщение отправиться не должно.
			testName:            "cms_валидный_1_лс_повторно",
			fileName:            "d831d274d9f248d791e34bf20943511b.cms",
			respCode:            http.StatusOK,
			tpStatusInDB:        entity.TpStatusEnumSuccess,
			natsMessagesCounter: 1,
			filesInDB:           8,
		},
		{
			// Валидный CMS с ТК
			testName:            "cms_валидный_тк",
			fileName:            "e8bc7dd3a26c4947bea2a6938a04ceb1.cms",
			respCode:            http.StatusOK,
			tpStatusInDB:        entity.TpStatusEnumSuccess,
			natsMessagesCounter: 2,
			filesInStorage:      []string{"e8bc7dd3a26c4947bea2a6938a04ceb1.cms", "6eec0e069a3c485f9c80459f3556e3a1.p7s", "86badd7c788c4b4cb75c09cff0548ecf.bin"},
			filesInDB:           3,
		},
		{
			// ТП с вложенностью директорий больше 1
			testName:            "cms_с_вложенностью_глубже_1",
			fileName:            "186e88f22cc542e0bef59de5267cf766.cms",
			respCode:            http.StatusBadRequest,
			tpStatusInDB:        entity.TpStatusEnumFailed,
			natsMessagesCounter: 2,
			filesInDB:           0,
		},
		{
			// ТП с 3 лс
			testName:            "cms_3_лс",
			fileName:            "11b7c042c146416b9892658df2065d45.cms",
			respCode:            http.StatusOK,
			tpStatusInDB:        entity.TpStatusEnumSuccess,
			natsMessagesCounter: 5,
			filesInStorage: []string{"11b7c042c146416b9892658df2065d45.cms", "2c73a264943c475ea8a8a2a008af68ea.p7s", "d0ed0f809d7f4ee5893ddf443c03ad21.bin",
				"870f00d721e046c2beecdf1900f65f5d.p7s", "ea4d3bfe4c454ae5b28bc71648a42d01.bin", "7aab072c2d764650b04693d05b314231.p7s", "1958a66bf43a45e58c02db3cc5c51522.bin"},
			filesInDB: 9,
		},
	}

	url := fmt.Sprintf("http://%s/receiver/cms/v1/", net.JoinHostPort(cfg.Gateway.HTTP.Host, cfg.Gateway.HTTP.Port))

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			resp, err := sendCMS(ctx, url, test.fileName, test.fileName)
			if err != nil {
				t.Fatalf(err.Error())
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				b, _ := io.ReadAll(resp.Body)
				fmt.Printf("body: %s\n", string(b))
			}

			assert.Equal(t, test.respCode, resp.StatusCode)

			resultTP, err := tpRepo.SelectByName(ctx, test.fileName)
			if err != nil {
				t.Fatalf(err.Error())
			}

			if test.respCode != http.StatusOK {
				require.NotEmpty(t, resultTP.ErrorCode)
			}

			assert.Equal(t, test.tpStatusInDB, resultTP.Status)
			assert.Equal(t, int32(test.natsMessagesCounter), *natsMessagesCounter)

			for _, v := range test.filesInStorage {
				_, err = storageMock.GetFile(ctx, v)
				require.NoError(t, err)
			}

			var res []struct {
				Dir string `db:"dir"`
				Doc string `db:"doc"`
			}

			err = conn.GetReadConnection().Select(&res, "SELECT d.name as dir, doc.name as doc FROM tp_directory d INNER JOIN tp_document doc ON d.id = doc.directory_id WHERE d.tp_id = $1", resultTP.ID)
			if err != nil {
				t.Fatalf(err.Error())
			}

			assert.Len(t, res, test.filesInDB)
		})
	}
}

func sendCMS(ctx context.Context, url string, fileName string, fileNameIsHeader string) (*http.Response, error) {
	cmsBytes, err := os.ReadFile("./testdata/" + fileName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(cmsBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Send-Receipt-To", "http://localhost:123/receiver/cms/v1/")
	req.Header.Add("Content-Disposition", "attachment;filename="+fileNameIsHeader)

	return http.DefaultClient.Do(req)
}

func natsSucsribe(counter *int32, natsConnString string) error {
	handler := func(_ context.Context, _ anats.Message) anats.MessageResultEnum {
		_ = atomic.AddInt32(counter, 1)
		return anats.MessageResultEnumSuccess
	}

	natsClient, err := anats.NewNatsClient(context.Background(), anats.NatsClientConfig{
		Urls:       []string{natsConnString},
		StreamName: broker.ReceiverStreamName,
		Subjects:   []string{broker.ReceiverSubjectNormalPriority},
	})

	if err != nil {
		return err
	}

	err = natsClient.Subscribe(context.Background(), "test_consumer", broker.ReceiverSubjectNormalPriority, handler, anats.SubscribeOptions{
		Workers:    1,
		MaxDeliver: 1,
	})

	if err != nil {
		return err
	}

	return nil
}

var crt = []byte(`-----BEGIN CERTIFICATE-----
MIIDbTCCAlUCFG3UrDPvF9Z5AqOBOtyvNLdoUZ8mMA0GCSqGSIb3DQEBCwUAMHMx
CzAJBgNVBAYTAlJVMQ8wDQYDVQQIDAZLYWx1Z2ExDzANBgNVBAcMBkthbHVnYTEW
MBQGA1UECgwNS2FsdWdhIEFzdHJhbDEqMCgGA1UEAwwhcmVncmVzc2lvbi4yYWUu
c3RhZ2luZy5rZXlkaXNrLnJ1MB4XDTIxMDcxNjA2MTcxNVoXDTMxMDcxNDA2MTcx
NVowczELMAkGA1UEBhMCUlUxDzANBgNVBAgMBkthbHVnYTEPMA0GA1UEBwwGS2Fs
dWdhMRYwFAYDVQQKDA1LYWx1Z2EgQXN0cmFsMSowKAYDVQQDDCFyZWdyZXNzaW9u
LjJhZS5zdGFnaW5nLmtleWRpc2sucnUwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQC6XqVAGHJA8qt49SE545C6RobY4IZFYExqLquz/qStzGWBOsTNvelF
Ub9W2SNvW2F2sG7avEw+JLaubrJ522tWB/bFC5dooQ9iY377cC76dJBavJ18hlDD
O9/lcWyMpZIw8mvTv2XuJMpxvauGRVhNqKF1mcv9eH+vYAO/0To+1YDCJjntr+p4
9mtqUCb26EfKEMcgwpfg4RbqGKItKfRyV5/2nO9tolhG23P9Hr8ruhD+7ckQg95l
J5BiDRLh/FIL8msRWTO//1e0z3VEWjy7zVB9OJ6ru8/V0LQxJSGbFI0d5kZ0WmMo
G3CG2czyF73Hdive7tCnF1Ts/jw+4NLNAgMBAAEwDQYJKoZIhvcNAQELBQADggEB
AJpdgDmvsqDF4mXu4inrapXcpGBYsBZB1piav4y5mW2V2hOwV1odvQXZ2F6JUrEj
bafMOLAmpI8D95ngT1AxQaUf7Fu4BUVdg112c8X83HBWI5rY0nRekzcP/PF9od2h
Xy7cUmivxiP0edAFW7v8SZrlECZgVFcpaFL/fcFcSAu//bQlP77RvUzFka3Z6Izg
Lh3L+w1KsE6VKbMWKqgCrgo6lxx+dn9Jpz+XH3vkT3nezj+NdHSx2VK39Ce3kgtS
XOqdieVIGi2Vy5NFe2miWYUqrZnzHQ5iq2ncrtVuG1GZRBkZgcRA7XjLj32Rjq7z
sFf4IF3IF5wcEJq3HK9aoVg=
-----END CERTIFICATE-----
`)

var key = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC6XqVAGHJA8qt4
9SE545C6RobY4IZFYExqLquz/qStzGWBOsTNvelFUb9W2SNvW2F2sG7avEw+JLau
brJ522tWB/bFC5dooQ9iY377cC76dJBavJ18hlDDO9/lcWyMpZIw8mvTv2XuJMpx
vauGRVhNqKF1mcv9eH+vYAO/0To+1YDCJjntr+p49mtqUCb26EfKEMcgwpfg4Rbq
GKItKfRyV5/2nO9tolhG23P9Hr8ruhD+7ckQg95lJ5BiDRLh/FIL8msRWTO//1e0
z3VEWjy7zVB9OJ6ru8/V0LQxJSGbFI0d5kZ0WmMoG3CG2czyF73Hdive7tCnF1Ts
/jw+4NLNAgMBAAECggEAL6h6S5y0wuULcmAFAh+uZ+QXIaiwzVbl05VhtSKeDA+j
uVtE7nPtDhvseRIH3LcaglZ10puqR43t5UdLfpvco+Bfe14OduQQ4hEdbMDwUn2y
WHG5OBnE11gdVjgeEQ1aCAhGCJz+PNrNpi5hiXF8Nke0GjWdE5FX3YoJC2k+oshY
cDzzE/ubtNVHSdDSKwUOgWrS12WybNnwmh4faYk74N/3kk34Yh6s8JnYk8+26gjO
dFLqDkWAboOpZi7B1Wnyr7ldNvFi8YJ2E90bC4GmoS4D0/hl/5HzYWql3Zmf5gmc
CSblS74C7TmPMKvGoaTcsas17GWFY0zU4bdPyzkkpQKBgQDcKoQokZ09COCQkCjV
QM9nP79127IE8KNXd4VsOBZLVGoM06ErgmkLaDX8e1OAzSMkX+HDZTfD2yTGJGx5
Q8wYBkkiC6PtfrM2x87OrIUkVG3815jKr1sMpHGd5N0u2KqEVR7dC6m7DJh6tkxR
RBblX9Kc/bXo75mwfi7HTe8FmwKBgQDYs/aFQUT76ny3qAC73MF2oP8d4UsFMJsm
W1xVYk0W9eDb28q8cqOdr/+mYIjz4/5gf48MawKOqYmwKwUBPlLO/hlkZB2y2Wss
quc3+YnQ8rwrnEW4F+9/Dx+k57NTmF5c4P3Q2BLKbXjHzT0J02+Kck11AmDagZx3
GjpI9IkDtwKBgHg+XDmP9bGM9KDfqv11TREV1up2l45tIri1lVAafcqciuMAfki2
C8roGnwPmvaAkw3ds/60fDVirX3uDLRaG9CPNkf61YfzJ8vmaoOj43+JAR0TXuZr
yS1pbogOo+JfARoPJzEQmp2G7owunhXQOzUBFZUaV8yld2nWMJQ3czC7AoGBAIJG
lnJ/zaAA1R94AZDu4uOVYCmvcnFZSjyh+f1ezmd6Q8cI+HWYGaLH1tJIAK1WqGuM
5AucHXp0k9Dz29tmg1PrUIqY4X3O1W6SA9UT0HVsKBGzrfpBcXqaNfTmUll0JW6C
2DQAYjON4mmDiilpEgpSMxyf5GgYOV8kxltrnx87AoGAH2CyRjCHER0h3Ky1z5Nk
JeAXtCLG1KAz3i/cdgn+2/Z49O2W6qwz/A9EchqByAzhn+fwNRGRRDpv9WXcCdn/
l6OvshyQx6HyAdTccJg6R1DAXe0puZmtw+0xNDaeGtt0vH9O9NwNcCkbNVZbxGRj
DKeg6mlKO0o+YQ+OVC+7wCM=
-----END PRIVATE KEY-----`)
