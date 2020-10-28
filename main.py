import flask
from flask import request
import datetime; 

from flaskext.mysql import MySQL
from hashlib import pbkdf2_hmac
import binascii
import base64


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


@app.route('/login', methods=['POST'])
def login():

    # cursor.execute("select * from user where")
    # data = cursor.fetchall()
    # print(data)
    # content = request.json
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
    print("bytes(password, 'utf-8'): ", bytes(password, 'utf-8').hex())
    
    key = pbkdf2_hmac(
        hash_name = 'sha1', 
        password = password.encode(), 
        salt = data[1].encode(),
        iterations = 10000, 
        dklen = 32
    )
    # check if new created hashed key equals to password saved in database
    if base64.b64encode(key).decode("utf-8") == data[0]:
        resp = {'status': 200, 'error': False, 'message': 'Logged In', 'time': datetime.datetime.now()} 
        return resp
    else :
        resp = {'status': 400, 'error': True, 'message': 'Login Failed', 'time': datetime.datetime.now()} 
        return resp



# Run the applicaiton
app.run(host='127.0.0.1', port=5005)
