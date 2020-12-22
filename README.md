```
git init
git remote add python-mongo git@github.com:ghumman/friends.git
touch README.md
git add README.md
git commit -m "Initial commit"
git push python-mongo master:python-mongo
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

## Create Django/Python Application [This is not going to be Django project, instead it will be a Flask project]
```
python manage.py startapp friends_app
```
## Installing Flask and Libraries
```
pip install Flask

```

## Running application
```
python main.py
```
