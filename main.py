import flask
from flask import request
import datetime; 

from flaskext.mysql import MySQL
from hashlib import pbkdf2_hmac
import binascii
import base64
from random import randrange
import smtplib, ssl
import uuid


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

    sendNewUserEmail(email, True, "")

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


@app.route('/change-password', methods=['POST'])
def changePassword():

    try:
        email = request.form['email']
        password = request.form['password']
        newPassword = request.form['newPassword']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
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
        salt = generateSalt()
        dbPassword = base64.b64encode(generateKey(newPassword, salt)).decode("utf-8") 
        stmt = ("UPDATE user SET salt=%s, password=%s where email=%s")
        data = (salt, dbPassword, email)
        cursor.execute(stmt, data)
        conn.commit()
        
        resp = {'status': 200, 'error': False, 'message': 'Password changed', 'time': datetime.datetime.now()} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Original password not right', 'time': datetime.datetime.now()} 
        return resp

@app.route('/forgot-password', methods=['POST'])
def forgotPassword():

    try:
        email = request.form['email']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Email is required', 'time': datetime.datetime.now()} 
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

    uuidNumber = uuid.uuid4()
    stmt = ("UPDATE user SET reset_token=%s where email=%s")
    data = (str(uuidNumber), email)
    cursor.execute(stmt, data)
    conn.commit()

    sendNewUserEmail(email, False, str(uuidNumber))
    
    resp = {'status': 200, 'error': False, 'message': 'Reset password is sent', 'time': datetime.datetime.now()} 
    return resp

@app.route('/reset-password', methods=['POST'])
def resetPassword():

    try:
        token = request.form['token']
        password = request.form['password']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
        return resp
    stmt = ("select id from user where reset_token=%s")
    data = (token)
    cursor.execute(stmt, data)
    data = cursor.fetchone()

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'Token is not valid', 
        'time': datetime.datetime.now()} 
        return resp
    userID = data[0]

    salt = generateSalt()
    dbPassword = base64.b64encode(generateKey(password, salt)).decode("utf-8") 
    stmt = ("UPDATE user SET salt=%s, password=%s, token=%s where id=%s")
    data = (salt, dbPassword, None, userID)
    cursor.execute(stmt, data)
    conn.commit()
    
    resp = {'status': 200, 'error': False, 'message': 'Password successfully reset', 'time': datetime.datetime.now()} 
    return resp

@app.route('/all-friends', methods=['POST'])
def allFriends():

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
        stmt = ("select first_name, last_name, email FROM user where email!=%s")
        data = (email)
        cursor.execute(stmt, data)
        rows = cursor.fetchall()
        users = []
        for r in rows: 
            users.append({"firstName" : r[0], "lastName" : r[1], "email" : r[2]})

        resp = {'status': 200, 'error': False, 'message': 'Friends attached', 'time': datetime.datetime.now(), 'usersAll': users} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp


def sendNewUserEmail(email, accountCreated, token):

    port = 465  # For SSL
    smtp_server = "smtp.gmail.com"
    sender_email = "example@gmail.com"  # Enter your address
    receiver_email = email  # Enter receiver address
    password = input("Type your password and press enter: ")

    message = ""
    if accountCreated:
        message = """\
        Subject: Welcome to Friends

        New Account Created."""
    else:
        message = """\
        Subject: Password Reset Request

        To reset your password, click the link below:
        http://localhost:3000/#/reset-password?token=""" 
        message = message + token



    context = ssl.create_default_context()
    with smtplib.SMTP_SSL(smtp_server, port, context=context) as server:
        server.login(sender_email, password)
        server.sendmail(sender_email, receiver_email, message)



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