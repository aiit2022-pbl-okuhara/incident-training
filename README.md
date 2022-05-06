# incident-training

## Initialize
```
$ docker-compose up --build
```

はじめは `incident-training` データベースが存在しないためエラーになります。初回だけコンテナに入って `CREATE DATABASE` をします。
```
$ docker ps -a
$ docker exec -it {Container ID} bin/bash
```

`admin` ユーザで `Postgres` に接続する。
```
$ psql -U admin
```

次のクエリを実行する。
```
CREATE DATABASE "incident-training" OWNER = admin TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'ja_JP.UTF-8' LC_CTYPE = 'ja_JP.UTF-8';
```

## How to launch the application
```
$ docker-compose up --build
```

Access http://0.0.0.0:8080

## How to connect the database
After launching your application with docker-compose, connect to the `incident-training-db` container with the `admin` user.
```
$ docker exec -it incident-training-db psql -U admin
```