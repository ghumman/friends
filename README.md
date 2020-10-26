```
git init
git remote add python git@github.com:ghumman/friends.git
touch README.md
git add README.md
git commit -m "Initial commit"
git push python master:python
```

## Installing Django
Use following guide to make sure python is set to 3.8
```
https://unix.stackexchange.com/questions/410579/change-the-python3-default-version-in-ubuntu
```
If Python 2 >=2.7.9 or Python 3 >=3.4 no need to install pip
```
python -m pip install Django
```
Otherwise
```
curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py
python get-pip.py
python -m pip install Django

```

## Initializing Django Project
```
django-admin startproject friends
```

## Create Django/Python Application
```
python manage.py startapp friends_app
```

## Running application
```
python manage.py runserver
```
