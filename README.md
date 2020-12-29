```
git init
git remote add python-django git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push python-django master:python-django
```

## Starting the server
```
python manage.py runserver
```

## Starting server at different port
```
python manage.py runserver 8080
```

## Calling APIs
Instead of calling endpoints like 
```
localhost:8000/login
```
call endpoints like following
```
localhost:8000/login/
```
