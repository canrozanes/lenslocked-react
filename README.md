To restart docker:
```
docker compose down
docker compose up -d
```


To login into psql database after running docker
```
docker exec -it lenslocked-db-1 /usr/bin/psql -U baloo -d lenslocked
```

To start the Go Server with file watching:

```
modd
```

To start the dev server for frontend development
```
cd frontend
npm run dev
```

and then visit `localhost:5173`

Above command starts the Vite server. When the vite server is running, the frontend is served directly from Vite but all request made to the `/api` endpoint is proxied to the Go backend
by Vite.

To build the frontend to serve it directly from the Go server

```
npm run build
```

and then visit `localhost:3000`
