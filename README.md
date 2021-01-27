# BeeWar Backend

## local quickstart

## deploy to heroku, db4free.net

install heroku cli, add remote to the heroku git repo (`heroku/master`)

in `origin/master` branch,

```cassandraql
// deploy
git push heroku master
// see logs
heroku logs --tail
```

Database used is db4free.net. Credentials available in heroku Procfile. To migrate, run migration manually in provided phpMyAdmin from db4free

```cassandraql
// seed from local
heroku run seed
```
