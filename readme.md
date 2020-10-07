# Maga Auctions - API Facade

## Como rodar

```
docker-compose up
```

As variáveis de ambiente em desenvolvimento são setadas em `api/cmd/config.yml` e para rodar em outros contextos basta exportá-las:
```
API_ENV: <environment>
API_PORT: <port>
LEGACY_URI: <uri>
```
___

## Como usar

### Dev Portal
```
http://localhost:8081
```

### API
```
http://localhost:8080/maga-auctions/v1/health-check
```