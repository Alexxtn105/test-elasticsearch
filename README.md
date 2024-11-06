# Тестовый проект с использованием ElasticSearch

Клиент для ElasticSearch
```bash
go get github.com/elastic/go-elasticsearch/v8@latest
```

Индексирование:
```bash
curl -X POST localhost:8087/blogs/index
```


# Установка ElasticSearch в Docker:
Тянем образ (latest не работает, нужно тянуть конкретную версию)
```bash
docker pull elasticsearch:8.15.3
```
Создаем сеть
```bash
docker network create somenetwork
```
Запускаем контейнер
```bash
docker run -d --name elasticsearch --net somenetwork -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" elasticsearch:8.15.3
```

В запущенном контейнере ```elasticsearch``` отключаем аутентификацию путем редактирования файла ```elasticsearch.yml``` в папке ```/usr/share/elasticsearch/config```:
```yaml
xpack.security.enabled: false
xpack.security.enrollment.enabled: false
```
Перезапустить контейнер