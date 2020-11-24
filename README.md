```
git init
git remote add node-js git@github.com:ghumman/friends.git
touch README.md
git add README.md
git commit -m "Initial commit"
git push node-js master:node-js
```

## Starting the server
```
node server.js
```

## Error when connecting to local mysql instant
```
Followed this link to fix this error
https://stackoverflow.com/questions/50093144/mysql-8-0-client-does-not-support-authentication-protocol-requested-by-server
```
