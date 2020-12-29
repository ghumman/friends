```
git init
git remote add php-psql git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push php-psql master:php-psql
```

## Starting the server
```
php -S localhost:8000 main.php
```

## To send email using smtp gmail email
```
sudo pear install --alldeps Mail
```

## To use postgreSQL with php PDO
```
sudo apt-get install php7.4-pgsql
```

