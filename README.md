# Parser Module

Этот модуль предназначен для парсинга событий с веб-сайта и отправки их в Kafka.

## Установка

1. Установите модуль:
    ```sh
    go get github.com/notickets/datamining
    ```

2. Импортируйте модуль в файле :
    ```go
    import "github.com/notickets/datamining"
    ```

3. Создайте файл `.env` и добавьте необходимые переменные окружения:
    ```env
    KAFKA_BROKER=your_kafka_broker
    KAFKA_TOPIC=your_kafka_topic
    PROXY_URL=your_proxy_url (если требуется)
    ```

4. Готово!