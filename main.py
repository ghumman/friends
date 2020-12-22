import flask
from flask import request
import datetime
import pymongo

from hashlib import pbkdf2_hmac
import binascii
import base64
from random import randrange
import smtplib, ssl
import uuid


app = flask.Flask(__name__)
app.config["DEBUG"] = True

myclient = pymongo.MongoClient("mongodb://localhost:27017/")
mydb = myclient["friends_mongo"]
userCollection = mydb["user"]
messageCollection = mydb["message"]


@app.route('/add-user', methods=['POST'])
def addUser():
    try:
        email = request.form['email']
        password = request.form['password']
        firstName = request.form['firstName']
        lastName = request.form['lastName']

    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
        return resp

    data = userCollection.find_one({"email" : email})

    if data is not None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Already Exists', 
        'time': datetime.datetime.now()} 
        return resp
    salt = generateSalt()
    dbPassword = base64.b64encode(generateKey(password, salt)).decode("utf-8") 

    newUser = { "authType": 0, "firstName": firstName, "lastName": lastName, "createdAt": datetime.datetime.now(), "salt": salt, "password": dbPassword, "email": email}

    userCollection.insert_one(newUser)

    # sendNewUserEmail(email, True, "")

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


    data = userCollection.find_one({"email" : email})

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, data["salt"])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == data["password"]:
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

    data = userCollection.find_one({"email" : email})

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, data["salt"])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == data["password"]:
        salt = generateSalt()
        dbPassword = base64.b64encode(generateKey(newPassword, salt)).decode("utf-8") 

        query = { "email": email }
        newValues = { "$set": { "salt": salt, "password": dbPassword } }

        userCollection.update_one(query, newValues)
        
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

    data = userCollection.find_one({"email" : email})

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp

    uuidNumber = uuid.uuid4()

    query = { "email": email }
    newValues = { "$set": { "resetToken": str(uuidNumber)} }

    userCollection.update_one(query, newValues)

    # sendNewUserEmail(email, False, str(uuidNumber))
    
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

    data = userCollection.find_one({"resetToken" : token})

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'Token is not valid', 
        'time': datetime.datetime.now()} 
        return resp

    salt = generateSalt()
    dbPassword = base64.b64encode(generateKey(password, salt)).decode("utf-8") 

    query = { "email": data["email"] }
    newValues = { "$set": { "salt": salt, "password": dbPassword, "resetToken": None } }

    userCollection.update_one(query, newValues)

    
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

    data = userCollection.find_one({"email" : email})

    if data is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'User Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, data["salt"])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == data["password"]:

        friends = userCollection.find({"email" : {"$ne": email}})

        users = []
        for friend in friends: 
            users.append({"firstName" : friend['firstName'], "lastName" : friend['lastName'], "email" : friend['email']})

        resp = {'status': 200, 'error': False, 'message': 'Friends attached', 'time': datetime.datetime.now(), 'usersAll': users} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp

@app.route('/send-message', methods=['POST'])
def sendMessage():

    try:
        message = request.form['message']
        messageFromEmail = request.form['messageFromEmail']
        messageToEmail = request.form['messageToEmail']
        password = request.form['password']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
        return resp

    dataSender = userCollection.find_one({"email" : messageFromEmail})


    if dataSender is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'Sender Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, dataSender["salt"])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == dataSender["password"]:

        dataReceiver = userCollection.find_one({"email" : messageToEmail})

        if dataReceiver is None: 
            resp = {'status': 400, 
            'error': True, 
            'message': 'Receiver Does Not Exist', 
            'time': datetime.datetime.now()} 
            return resp

        # Sender credentials are correct and both sender and receiver exists
        saveMessage(message, dataSender, dataReceiver)

        resp = {'status': 200, 'error': False, 'message': 'Message sent', 'time': datetime.datetime.now()} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp

@app.route('/messages-user-and-friend', methods=['POST'])
def messagesUserAndFriend():

    try:
        userEmail = request.form['userEmail']
        friendEmail = request.form['friendEmail']
        password = request.form['password']
    except: 
        resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
        return resp

    dataSender = userCollection.find_one({"email" : userEmail})

    if dataSender is None: 
        resp = {'status': 400, 
        'error': True, 
        'message': 'Sender Does Not Exist', 
        'time': datetime.datetime.now()} 
        return resp
    
    key = generateKey(password, dataSender["salt"])
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == dataSender["password"]:

        # Check if receiver exists
        dataReceiver = userCollection.find_one({"email" : friendEmail})

        if dataReceiver is None: 
            resp = {'status': 400, 
            'error': True, 
            'message': 'Receiver Does Not Exist', 
            'time': datetime.datetime.now()} 
            return resp

        conversations = messageCollection.find({"$or" : [{"messageFrom": dataSender, "messageTo": dataReceiver}, {"messageFrom": dataReceiver, "messageTo": dataSender}]})

        messages = []
        for conversation in conversations: 
            messages.append({"message" : conversation["message"], "messageFromEmail" : conversation["messageFrom"]["email"], "messageToEmail" : conversation["messageTo"]["email"], "sentAt" : conversation["sentAt"]})

        resp = {'status': 200, 'error': False, 'message': 'Messages attached', 'time': datetime.datetime.now(), 'msgs': messages} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp


def saveMessage(message, sender, receiver):

    newMessage = { "message": message, "messageFrom": sender, "messageTo": receiver, "sentAt": datetime.datetime.now() }
    messageCollection.insert_one(newMessage)
    


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
