# Tour Service 

Go + Gin mikroservis za upravljanje turama i ključnim tačkama, sa MongoDB dokumentnom bazom.

- NoSQL dokumentna baza: MongoDB

## Pokretanje MongoDB

```powershell
& "C:\Program Files\MongoDB\Server\7.0\bin\mongod.exe" --dbpath "C:\mongodb-data\db" --bind_ip 127.0.0.1 --port 27017
```

## Pokretanje Go backend-a

```bash
cd services/tour-service
go mod tidy
go run main.go
```

Servis sluša na `http://localhost:8085`.

`MONGO_URI` je opciona env promenljiva. Ako nije postavljena, koristi se:
- `mongodb://localhost:27017`

Mongo konfiguracija:
- Database: `tourist_app_tours`
- Collection: `tours`

## Pokretanje Angular frontend-a

```bash
cd frontend/tourist-app-ui
npm install
ng serve
```

