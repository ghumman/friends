var express = require('express');
var app = express();
var bodyParser = require('body-parser');
var fs = require("fs");
const crypto = require('crypto');
const { v4: uuidv4 } = require('uuid');

// create application/x-www-form-urlencoded parser
var urlencodedParser = bodyParser.urlencoded({ extended: false });

var mysql = require('mysql');

var con = mysql.createConnection({
  host: "localhost",
  user: "ghumman",
  password: "ghumman",
  database : "friends_mysql"
});

con.connect(function(err) {
  if (err) throw err;
});

// User Related Endpoints
// add-user
app.post('/add-user', urlencodedParser, function (req, res) {
	const email = req.body.email;
	const password = req.body.password;
	const firstName = req.body.firstName;
	const lastName = req.body.lastName;

	if (email == undefined || password == undefined || firstName == undefined || lastName == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}

	var sqlCheckUser = mysql.format("select password, salt from user where email=?", [email]);
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.length == 0) {

			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(password, salt);
			var sqlGetID = mysql.format("select id from user order by id desc limit 1");
			con.query(sqlGetID, function (errorGetID, resultsGetID, fields) {
				if (errorGetID) throw errorGetID;
				var newUserID = 1;
				console.log('value of resultsGetID.id: ', resultsGetID[0].id );
				if (resultsGetID.length > 0) {
					newUserID = resultsGetID[0].id + 1;
				}

				console.log('newUserID: ', newUserID);
				console.log('new Date(Date.now()): ', new Date(Date.now()));
				console.log('dbPassword: ', dbPassword);
				console.log('salt: ', salt);
				console.log('firstName: ', firstName);
				console.log('lastName: ', lastName);
				console.log('email: ', email);

				// create new user using id, auth_type, created_at, email, firstName, lastName, dbPassword, salt 
				var sqlCreateUser = mysql.format("insert into user (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", [newUserID, 0, new Date(Date.now()) , email, firstName, lastName, dbPassword, salt]);
				con.query(sqlCreateUser, function (errorCreateUser, resultsCreateUser, fields) {
					if (errorCreateUser) throw errorCreateUser;
					// sendNewUserEmail(email, true, "")
					return res.send({'status': 200, 'error': false, 'message': 'User Created', 'time': new Date(Date.now())} );
				});
			});
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'User Already Exists', 'time': new Date(Date.now())} );
		}
	});
})

// login
app.post('/login', urlencodedParser, function (req, res) {

	const email = req.body.email;
	const password = req.body.password;

	if (email == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Email or Password missing', 'time': new Date(Date.now())} );
	}
	var sql = mysql.format("select password, salt from user where email=?", [email]);
	con.query(sql, function (error, results, fields) {
		if (error) throw error;
		if (results.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, results[0].salt);
		if (results[0].password === key) {
			return res.send({'status': 200, 'error': false, 'message': 'Logged In', 'time': new Date(Date.now())} );
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
		}
	});
})

// change-password
app.post('/change-password', urlencodedParser, function (req, res) {
	const email = req.body.email;
	const password = req.body.password;
	const newPassword = req.body.newPassword;

	if (email == undefined || password == undefined || newPassword == undefined ) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}

	var sqlLogin = mysql.format("select password, salt from user where email=?", [email]);
	con.query(sqlLogin, function (errorLogin, resultsLogin, fields) {
		if (errorLogin) throw errorLogin;
		if (resultsLogin.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsLogin[0].salt);
		if (resultsLogin[0].password === key) {
			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(newPassword, salt);
			
			// update user with newly generated salt and dbPassword
			var sqlUpdateUser = mysql.format("UPDATE user SET salt=?, password=? where email=?", [salt, dbPassword, email]);
			con.query(sqlUpdateUser, function (errorUpdateUser, resultsUpdateUser, fields) {
				if (errorUpdateUser) throw errorUpdateUser;
				return res.send({'status': 200, 'error': false, 'message': 'Password changed', 'time': new Date(Date.now())} );
			});
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'Original password not right', 'time': new Date(Date.now())} );
		}
	});
})

// forgot-password
app.post('/forgot-password', urlencodedParser, function (req, res) {
	const email = req.body.email;

	if (email == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Email is required', 'time': new Date(Date.now())} );
	}

	var sqlLogin = mysql.format("select password, salt from user where email=?", [email]);
	con.query(sqlLogin, function (errorLogin, resultsLogin, fields) {
		if (errorLogin) throw errorLogin;
		if (resultsLogin.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		var uuidNumber = uuidv4();
		var sqlUpdateUser = mysql.format("UPDATE user SET reset_token=? where email=?", [uuidNumber, email]);
		con.query(sqlUpdateUser, function (errorUpdateUser, resultsUpdateUser, fields) {
			if (errorUpdateUser) throw errorUpdateUser;

			// sendNewUserEmail(email, false, uuidNumber)
			return res.send({'status': 200, 'error': false, 'message': 'Reset password is sent', 'time': new Date(Date.now())} );
		});
	});

})

// reset-password
app.post('/reset-password', urlencodedParser, function (req, res) {
	const token = req.body.token;
	const password = req.body.password;

	if (token == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}
	var sqlCheckToken = mysql.format("select id from user where reset_token=?", [token]);
	con.query(sqlCheckToken, function (errorCheckToken, resultsCheckToken, fields) {
		if (errorCheckToken) throw errorCheckToken;
		if (resultsCheckToken.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Token is not valid', 'time': new Date(Date.now())} );
		} else {
			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(password, salt);
			// update user using salt, password, token using id
			var sqlUpdateUser = mysql.format("UPDATE user SET salt=?, password=?, reset_token=? where id=?", [salt, dbPassword, null, resultsCheckToken[0].id]);
			con.query(sqlUpdateUser, function (errorCreateUser, resultsCreateUser, fields) {
				if (errorCreateUser) throw errorCreateUser;

				return res.send({'status': 200, 'error': false, 'message': 'Password successfully reset', 'time': new Date(Date.now())} );
			});
		}
	});
})

// all-friends
app.post('/all-friends', urlencodedParser, function (req, res) {
	const email = req.body.email;
	const password = req.body.password;

	if (email == undefined || password == undefined ) {
		return res.send({'status': 400, 'error': true, 'message': 'Email or Password missing', 'time': new Date(Date.now())} );
	}

	var sqlCheckUser = mysql.format("select password, salt from user where email=?", [email]);
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckUser[0].salt);
		if (resultsCheckUser[0].password === key) {
			var sqlFriends = mysql.format("select first_name, last_name, email FROM user where email!=?", [email]);
			con.query(sqlFriends, function (errorFriends, resultsFriends, fields) {
				if (errorFriends) throw errorFriends;
				var users = [];
				for (i=0; i<resultsFriends.length; i++) {
					users.push({"firstName" : resultsFriends[i].first_name, "lastName" : resultsFriends[i].last_name, "email" : resultsFriends[i].email})
				}
				return res.send({'status': 200, 'error': false, 'message': 'Friends attached', 'time': new Date(Date.now()), 'usersAll' : users} );
			});
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
		}
	});
})

// Messages Related Endpoints
// send-message
app.post('/send-message', urlencodedParser, function (req, res) {

})

// messages-user-and-friend
app.post('/messages-user-and-friend', urlencodedParser, function (req, res) {

})


// helper functions
// generates Key which is saved as database password
const generateKey = (password, salt) => {
	return crypto.pbkdf2Sync(
		password, 
		salt, 
		10000, 
		32, 
		'sha1'
		).toString('base64');
}

// generates salt
const generateSalt = () => {
	const ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    returnValue = ""
    for(x=0; x<30;x++)
        returnValue = returnValue + ALPHABET[Math.floor(Math.random() * Math.floor(ALPHABET.length))]
    return returnValue
}

// send email
/*
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
*/


// start the app
app.listen(8080)
