## storage-client

Данный пакет является оберткой поверх gRPC клиента единого файлового хранилища (сервис storage). Это нужно для более
удобного использования API хранилища без необходимости каждый раз реализовывать gRPC клиент.

Про взаимодействие с единым файловым хранилищем можно подробнее почитать
тут: https://track.astral.ru/soft/wiki/pages/viewpage.action?pageId=3163259010


Как сделать Docker-образ с storage-client
=============================================

По запросам коллег https://t.me/c/1517677573/15407

```shell

  # cd ~/projects/astral/Astral.Edo.Backeng.Go
  # docker build -t harbor.infra.yandex.astral-dev.ru/astral-edo/go/storage-cli:latest -f modules/storage-client/Dockerfile  .
  # docker push harbor.infra.yandex.astral-dev.ru/astral-edo/go/storage-cli:latest
```

Примеры использования storage-client
=============================================

Достать файл с ингреса демо стенда без авторизации, но с TLS. Сохранить файл в директорию /home/vodolaz095/temp/

```shell
storage-client \
    --addr=storage-grpc.edo-demo.cloud.astral-dev.ru:443 \
    --tls=true \
    --id=7b6b3a57-62b2-46f7-a88c-ab6b8a92f9bb \
    --token="" \
    --output=/home/vodolaz095/temp/
```

Загрузить файл в storage, порт 8090 которого прокинут через lens + port-forwardig

```shell

storage-client \
 --addr 127.0.0.1:8090 \
 --id 37f1c8b0-16f0-4a04-ac44-765b20339572 \
 --file ~/temp/poa_37f1c8b0-16f0-4a04-ac44-765b20339572.xml


```

Получить атрибуты файла из igress'а storage демо стенда без авторизации, но с TLS (в предположении, 
что истёк сертификат, и мы его не проверяем) 

```shell

storage-client \
    --addr=storage-grpc.edo-demo.cloud.astral-dev.ru:443 \
    --tls=true \
    --insecure=true \
    --id=7b6b3a57-62b2-46f7-a88c-ab6b8a92f9bb \
    --token="" \
    -a

```

Пример ответа
```json
{
  "Id": "7b6b3a57-62b2-46f7-a88c-ab6b8a92f9bb",
  "Attrs": {
    "Created": 1666340608,
    "Expires": 1824107008,
    "Creator": "edo-document-store",
    "CustomId": "",
    "Filename": "DP_IZVPOL_2AE272E0C82-EE25-40DC-9610-DFDA77FA574A_2AE4D776B8E-09CA-416B-ACD0-3B9762375CE7_20221021_2E4788E0-44BE-439F-BD60-9D9A41C705AB.xml",
    "Size": 5229,
    "StorageType": 0,
    "Type": "Document",
    "SubType": "Izvpol",
    "Readonly": true,
    "Protected": true
  },
  "ServiceAttrs": {
    "edo-document-store": {
      "TTL": 1824107008
    }
  }
}

```

Удалить файл из igress'а storage демо стенда c авторизацией, TLS

```shell

storage-client \
    --addr=storage-grpc.edo-demo.cloud.astral-dev.ru:443 \
    --tls=true \
    --id=7b6b3a57-62b2-46f7-a88c-ab6b8a92f9bb \
    --del=true \
    --token="eyJhbGciOiJIUzI1NiJ9.eyJzZXJ2aWNlIjoic3RvcmF0ZS1jbGkiLCJpc3MiOiJhc3RyYWwiLCJpYXQiOjE1MTYyMzkwMjJ9.orAlNlr_7aNM-5xJE88fhJInSw3_xtbgA60hZazmiXo"

```
