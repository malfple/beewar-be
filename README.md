# BeeWar Backend

## Quickstart

Setup local config file: Copy `config.template.json` to `config.json`.
This template config may be different for each person, so change it accordingly.
You can also commit your own config file if you want (e.g. `config.malfple.json`).

Get golang (version, see `go.mod`)

Get mysql server.

```cassandraql
// install dep
go mod vendor
// migration and seed. make a db named 'beewar' first
mysql -u root(or username) -p beewar < tools/db/migration.sql
go run tools/db/seeder_go.go
// run the main service
make run
```

Before commit, you can check everything using:
```cassandraql
make check
```

Check out `Makefile` to see what this command does and to see other useful commands!

## deploy to heroku, db4free.net

install heroku cli, add remote to the heroku git repo (`heroku/master`)

in `origin/master` branch,

```cassandraql
// add remote
heroku git:remote -a beewar-be
// deploy
git push heroku master
// see logs
heroku logs --tail -a beewar-be
```

Database used is db4free.net. Credentials available in heroku Procfile. To migrate, run migration manually in provided phpMyAdmin from db4free.
Seeding however, is done inside the heroku container:

```cassandraql
// seed from local
heroku run seed
```
