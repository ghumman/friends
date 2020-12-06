```
git init
git remote add php git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push php master:php
```

## Starting the server
```
php -S localhost:8000 main.php
```

## To send email using smtp gmail email
```
sudo pear install --alldeps Mail
```

