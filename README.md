```
git init
git remote add ruby-psql git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push ruby-psql master:ruby-psql
```

## Starting the server
```
rails server
```
## Ruby issues with mysql
Had to recreate project using command
```
rails new ruby -d mysql
```
Following line resolved the issue of ruby gems like mysql not getting installed
```
sudo apt install default-libmysqlclient-dev
```
