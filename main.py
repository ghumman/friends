import flask
from flask import request
import datetime; 

from flaskext.mysql import MySQL
from hashlib import pbkdf2_hmac
import binascii
import base64
from random import randrange


app = flask.Flask(__name__)
app.config["DEBUG"] = True

mysql = MySQL()
app.config['MYSQL_DATABASE_USER'] = 'ghumman'
app.config['MYSQL_DATABASE_PASSWORD'] = 'ghumman'
app.config['MYSQL_DATABASE_DB'] = 'friends_mysql'
app.config['MYSQL_DATABASE_HOST'] = 'localhost'
mysql.init_app(app)

conn = mysql.connect()
cursor =conn.cursor()

@app.route('/add-user', methods=['POST'])
def addUser():

    try:
        email = request.form['email']
        password = request.form['password']
        firstName = request.form['first_name']
        lastName = request.form['last_name']

    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
        return resp
    stmt = ("select password, salt from user where email=%s")
    data = (email)
    cursor.execute(stmt, data)
    data = cursor.fetchone()

    if data is not None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Already Exists', 
        'time': datetime.datetime.now()} 
        return resp
    salt = generateSalt()
    dbPassword = base64.b64encode(generateKey(password, salt)).decode("utf-8") 

    # Get the last id to create id for user as our table is not auto increment on id column
    stmt = ("select id from user order by id desc limit 1")
    cursor.execute(stmt)
    data = cursor.fetchone()
    userID = 0
    if data is not None: 
        userID = data[0] + 1
    else: 
        userID = 1

    # Create the new user
    stmt = ("insert into user (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES (%s, %s, %s, %s, %s, %s, %s, %s)")
    data = (userID, 0, datetime.datetime.now() , email, firstName, lastName, dbPassword, salt)
    cursor.execute(stmt, data)
    conn.commit()

    resp = {'status': 200, 'error': False, 'message': 'User Created', 'time': datetime.datetime.now()} 
    return resp



@app.route('/login', methods=['POST'])
def login():

    try:
        email = request.form['email']
        password = request.form['password']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Email or Password missing', 'time': datetime.datetime.now()} 
        return resp
    stmt = ("select password, salt from user where email=%s")
    data = (email)
    cursor.execute(stmt, data)
    data = cursor.fetchone()

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, data[1])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == data[0]:
        resp = {'status': 200, 'error': False, 'message': 'Logged In', 'time': datetime.datetime.now()} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp

def generateKey(password, salt):
    return pbkdf2_hmac(
        hash_name = 'sha1', 
        password = password.encode(), 
        salt = salt.encode(),
        iterations = 10000, 
        dklen = 32
    )
def generateSalt(): 
    ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    returnValue = ""
    for x in range(30):
        returnValue = returnValue + ALPHABET[randrange(len(ALPHABET)) ]
    return returnValue


# Run the applicaiton
app.run(host='127.0.0.1', port=5005)
