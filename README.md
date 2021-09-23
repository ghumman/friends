```
git init
git remote add node-js-psql git@github.com:ghumman/friends.git
touch README.md
git add README.md
git commit -m "Initial commit"
git push node-js-psql master:node-js-psql
```

## Starting the server
```
node server.js
```

## Setting up psql on ubuntu
```
Make sure you have psql installed on ubuntu, you can make sure by typing
psql --version
Basic commands: 
sudo -i -u postgres
psql
To list database
\l
To connect to a database
\c friends_psql
To list all tables
\d
```
