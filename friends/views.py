from django.http import JsonResponse
import datetime
from django.db import connection

from hashlib import pbkdf2_hmac
import binascii
import base64
from random import randrange
import smtplib, ssl
import uuid

def addUser(request):
    if request.method == "POST":
        try:
            email = request.POST.get('email')
            password = request.POST.get('password')
            firstName = request.POST.get('first_name')
            lastName = request.POST.get('last_name')

        except: 
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

        if email is None or password is None or firstName is None or lastName is None: 
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        
        with connection.cursor() as cursor:

            stmt = ("select password, salt from user where email=%s")
            data = [email]
            cursor.execute(stmt, data)
            data = cursor.fetchone()

            if data is not None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'User Already Exists', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)
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
            data = [userID, 0, datetime.datetime.now() , email, firstName, lastName, dbPassword, salt]
            cursor.execute(stmt, data)

            sendNewUserEmail(email, True, "")

            resp = {'status': 200, 'error': False, 'message': 'User Created', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)


def login(request):
    if request.method == "POST":
        try:
            email = request.POST.get('email')
            password = request.POST.get('password')
        except: 
            resp = {'status': 400, 'error': True, 'message': 'Email or Password missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        if email is None or password is None: 
            resp = {'status': 400, 'error': True, 'message': 'Email or Password missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)



        with connection.cursor() as cursor:

            stmt = ("select password, salt from user where email=%s")
            data = (email)
            cursor.execute(stmt, [email])
            data = cursor.fetchone()

            if data is None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'User Does Not Exist', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)

            key = generateKey(password, data[1])
            # check if new created hashed key equals to password saved in database
            if base64.b64encode(key).decode("utf-8") == data[0]:
                resp = {'status': 200, 'error': False, 'message': 'Logged In', 'time': datetime.datetime.now()} 
                return JsonResponse(resp)
            else :
                resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
                return JsonResponse(resp)

def changePassword(request):

    if request.method == "POST":

        try:
            email = request.POST.get('email')
            password = request.POST.get('password')
            newPassword = request.POST.get('newPassword')
        except: 
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)


        if email is None or password is None or newPassword is None: 
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        
        with connection.cursor() as cursor:
            stmt = ("select password, salt from user where email=%s")
            data = [email]
            cursor.execute(stmt, data)
            data = cursor.fetchone()

            if data is None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'User Does Not Exist', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)
            
            key = generateKey(password, data[1])
            # check if new created hashed key equals to password saved in database
            if base64.b64encode(key).decode("utf-8") == data[0]:
                salt = generateSalt()
                dbPassword = base64.b64encode(generateKey(newPassword, salt)).decode("utf-8") 
                stmt = ("UPDATE user SET salt=%s, password=%s where email=%s")
                data = [salt, dbPassword, email]
                cursor.execute(stmt, data)
                
                resp = {'status': 200, 'error': False, 'message': 'Password changed', 'time': datetime.datetime.now()} 
                return JsonResponse(resp)
            else :
                resp = {'status': 400, 'error': True, 'message': 'Original password not right', 'time': datetime.datetime.now()} 
                return JsonResponse(resp)

def forgotPassword(request):
    if request.method == "POST":
        try:
            email = request.POST.get('email')
        except: 
            resp = {'status': 400, 'error': True, 'message': 'Email is required', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        if email is None:
            resp = {'status': 400, 'error': True, 'message': 'Email is required', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

        with connection.cursor() as cursor:

            stmt = ("select password, salt from user where email=%s")
            data = [email]
            cursor.execute(stmt, data)
            data = cursor.fetchone()

            if data is None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'User Does Not Exist', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)

            uuidNumber = uuid.uuid4()
            stmt = ("UPDATE user SET reset_token=%s where email=%s")
            data = [str(uuidNumber), email]
            cursor.execute(stmt, data)

            sendNewUserEmail(email, False, str(uuidNumber))
            
            resp = {'status': 200, 'error': False, 'message': 'Reset password is sent', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

def resetPassword(request):
    if request.method == "POST":
        try:
            token = request.POST.get('token')
            password = request.POST.get('password')
        except: 
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        
        if token is None or password is None:
            resp = {'status': 400, 'error': True, 'message': 'Required data missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)
        
        with connection.cursor() as cursor:
        
            stmt = ("select id from user where reset_token=%s")
            data = [token]
            cursor.execute(stmt, data)
            data = cursor.fetchone()

            if data is None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'Token is not valid', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)
            userID = data[0]

            salt = generateSalt()
            dbPassword = base64.b64encode(generateKey(password, salt)).decode("utf-8") 
            stmt = ("UPDATE user SET salt=%s, password=%s, reset_token=%s where id=%s")
            data = [salt, dbPassword, None, userID]
            cursor.execute(stmt, data)
            
            resp = {'status': 200, 'error': False, 'message': 'Password successfully reset', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

def allFriends(request):
    if request.method == "POST":
        try:
            email = request.POST.get('email')
            password = request.POST.get('password')
        except: 
            resp = {'status': 400, 'error': True, 'message': 'Email or Password missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

        if email is None or password is None:
            resp = {'status': 400, 'error': True, 'message': 'Email or Password missing', 'time': datetime.datetime.now()} 
            return JsonResponse(resp)

        with connection.cursor() as cursor:
        
            stmt = ("select password, salt from user where email=%s")
            data = [email]
            cursor.execute(stmt, data)
            data = cursor.fetchone()

            if data is None: 
                resp = {'status': 400, 
                'error': True, 
                'message': 'User Does Not Exist', 
                'time': datetime.datetime.now()} 
                return JsonResponse(resp)
            
            key = generateKey(password, data[1])
            # check if new created hashed key equals to password saved in database
            if base64.b64encode(key).decode("utf-8") == data[0]:
                stmt = ("select first_name, last_name, email FROM user where email!=%s")
                data = [email]
                cursor.execute(stmt, data)
                rows = cursor.fetchall()
                users = []
                for r in rows: 
                    users.append({"firstName" : r[0], "lastName" : r[1], "email" : r[2]})

                resp = {'status': 200, 'error': False, 'message': 'Friends attached', 'time': datetime.datetime.now(), 'usersAll': users} 
                return JsonResponse(resp)
            else :
                resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
                return JsonResponse(resp)

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
