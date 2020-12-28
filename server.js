var express = require('express');
var app = express();
var bodyParser = require('body-parser');
var fs = require("fs");
const crypto = require('crypto');
const { v4: uuidv4 } = require('uuid');
var nodemailer = require('nodemailer');

// create application/x-www-form-urlencoded parser
var urlencodedParser = bodyParser.urlencoded({ extended: false });

var mysql = require('mysql');
const { Client } = require('pg')


// var con = mysql.createConnection({
//   host: "localhost",
//   user: "ghumman",
//   password: "ghumman",
//   database : "friends_mysql"
// });


const con = new Client({
	user: 'postgres',
	host: 'localhost',
	database: 'friends_psql',
	password: 'postgres',
	port: 5432,
  })
// con.connect()
var transporter = nodemailer.createTransport({
	service: 'gmail',
	auth: {
	  user: 'example@gmail.com',
	  pass: 'xxxxxx'
	}
  });

con.connect(function(err) {
  if (err) throw err;
});

app.use(function (req, res, next) {

	    res.setHeader('Access-Control-Allow-Origin', 'http://localhost:3000');
	    res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS, PUT, PATCH, DELETE');
	    res.setHeader('Access-Control-Allow-Headers', 'X-Requested-With,content-type');
	    res.setHeader('Access-Control-Allow-Credentials', true);

	    next();
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

	// var sqlCheckUser = mysql.format("select password, salt from user where email=?", [email]);
	const sqlCheckUser = {
		// give the query a unique name
		name: 'fetch-user',
		text: 'select password, salt from users where email=$1',
		values: [email],
	}
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.rows.length == 0) {

			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(password, salt);
			// var sqlGetID = mysql.format("select id from user order by id desc limit 1");
			const sqlGetID = {
				// give the query a unique name
				name: 'fetch-id',
				text: 'select id from users order by id desc limit 1',
			}
			con.query(sqlGetID, function (errorGetID, resultsGetID, fields) {
				if (errorGetID) throw errorGetID;
				var newUserID = 1;
				if (resultsGetID.rows.length > 0) {
					newUserID = resultsGetID.rows[0].id + 1;
				}

				// create new user using id, auth_type, created_at, email, firstName, lastName, dbPassword, salt 
				// var sqlCreateUser = mysql.format("insert into user (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", [newUserID, 0, new Date(Date.now()) , email, firstName, lastName, dbPassword, salt]);
				const sqlCreateUser = {
					// give the query a unique name
					name: 'create-user',
					text: 'insert into users (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)',
					values: [newUserID, 0, new Date(Date.now()) , email, firstName, lastName, dbPassword, salt],
				}
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

	const sql = {
		// give the query a unique name
		name: 'fetch-user',
		text: 'select password, salt from users where email=$1',
		values: [email],
	}
	con.query(sql, function (error, results) {
		if (error) throw error;
		if (results.rows.length == 0 ) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, results.rows[0].salt);
		if (results.rows[0].password === key) {
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

	// var sqlLogin = mysql.format("select password, salt from users where email=?", [email]);
	const sqlLogin = {
		// give the query a unique name
		name: 'fetch-user',
		text: 'select password, salt from users where email=$1',
		values: [email],
	}
	con.query(sqlLogin, function (errorLogin, resultsLogin, fields) {
		if (errorLogin) throw errorLogin;
		if (resultsLogin.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsLogin.rows[0].salt);
		if (resultsLogin.rows[0].password === key) {
			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(newPassword, salt);
			
			// update user with newly generated salt and dbPassword
			// var sqlUpdateUser = mysql.format("UPDATE users SET salt=?, password=? where email=?", [salt, dbPassword, email]);
			const sqlUpdateUser = {
				// give the query a unique name
				name: 'update-user',
				text: 'UPDATE users SET salt=$1, password=$2 where email=$3',
				values: [salt, dbPassword, email],
			}
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

	// var sqlLogin = mysql.format("select password, salt from users where email=?", [email]);
	const sqlLogin = {
		// give the query a unique name
		name: 'fetch-user',
		text: 'select password, salt from users where email=$1',
		values: [email],
	}
	con.query(sqlLogin, function (errorLogin, resultsLogin, fields) {
		if (errorLogin) throw errorLogin;
		if (resultsLogin.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		var uuidNumber = uuidv4();
		// var sqlUpdateUser = mysql.format("UPDATE users SET reset_token=? where email=?", [uuidNumber, email]);
		const sqlUpdateUser = {
			// give the query a unique name
			name: 'update-user',
			text: 'UPDATE users SET reset_token=$1 where email=$2',
			values: [uuidNumber, email],
		}
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
	// var sqlCheckToken = mysql.format("select id from users where reset_token=?", [token]);
	const sqlCheckToken = {
		// give the query a unique name
		name: 'fetch-user2',
		text: 'select id from users where reset_token=$1',
		values: [token],
	}
	con.query(sqlCheckToken, function (errorCheckToken, resultsCheckToken, fields) {
		if (errorCheckToken) throw errorCheckToken;
		if (resultsCheckToken.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Token is not valid', 'time': new Date(Date.now())} );
		} else {
			// generate salt
			const salt = generateSalt();
			// generate database password
			const dbPassword = generateKey(password, salt);
			// update user using salt, password, token using id
			// var sqlUpdateUser = mysql.format("UPDATE users SET salt=?, password=?, reset_token=? where id=?", [salt, dbPassword, null, resultsCheckToken[0].id]);
			const sqlUpdateUser = {
				// give the query a unique name
				name: 'update-user2',
				text: 'UPDATE users SET salt=$1, password=$2, reset_token=$3 where id=$4',
				values: [salt, dbPassword, null, resultsCheckToken.rows[0].id],
			}
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

	// var sqlCheckUser = mysql.format("select password, salt from users where email=?", [email]);
	const sqlCheckUser = {
		// give the query a unique name
		name: 'fetch-user3',
		text: 'select password, salt from users where email=$1',
		values: [email],
	} 
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckUser.rows[0].salt);
		if (resultsCheckUser.rows[0].password === key) {
			// var sqlFriends = mysql.format("select first_name, last_name, email FROM users where email!=?", [email]);
			const sqlFriends = {
				// give the query a unique name
				name: 'fetch-friends',
				text: 'select first_name, last_name, email FROM users where email!=$1',
				values: [email],
			}
			con.query(sqlFriends, function (errorFriends, resultsFriends, fields) {
				if (errorFriends) throw errorFriends;
				var users = [];
				for (i=0; i<resultsFriends.rows.length; i++) {
					users.push({"firstName" : resultsFriends.rows[i].first_name, "lastName" : resultsFriends.rows[i].last_name, "email" : resultsFriends.rows[i].email})
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
	const message = req.body.message;
	const messageFromEmail = req.body.messageFromEmail;
	const messageToEmail = req.body.messageToEmail;
	const password = req.body.password;

	if (message == undefined || messageFromEmail == undefined || messageToEmail == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}
	// var sqlCheckSender = mysql.format("select password, salt, id from users where email=?", [messageFromEmail]);
	const sqlCheckSender = {
		// give the query a unique name
		name: 'fetch-sender',
		text: 'select password, salt, id from users where email=$1',
		values: [messageFromEmail],
	}
	con.query(sqlCheckSender, function (errorCheckSender, resultsCheckSender, fields) {
		if (errorCheckSender) throw errorCheckSender;
		if (resultsCheckSender.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckSender.rows[0].salt);
		if (resultsCheckSender.rows[0].password === key) {
			// var sqlCheckReceiverExist = mysql.format("select password, salt, id from users where email=?", [messageToEmail]);
			const sqlCheckReceiverExist = {
				// give the query a unique name
				name: 'fetch-receiver',
				text: 'select password, salt, id from users where email=$1',
				values: [messageToEmail],
			}
			con.query(sqlCheckReceiverExist, function (errorCheckReceiverExist, resultsCheckReceiverExist, fields) {
				if (errorCheckReceiverExist) throw errorCheckReceiverExist;
				if (resultsCheckReceiverExist.rows.length == 0) {
					return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
				}
				// Sender credentials are correct and both sender and receiver exists
				saveMessage(message, resultsCheckSender.rows[0].id, resultsCheckReceiverExist.rows[0].id);
				return res.send({'status': 200, 'error': false, 'message': 'Message sent', 'time': new Date(Date.now())} );
			});
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
		}
	});
})

// messages-user-and-friend
app.post('/messages-user-and-friend', urlencodedParser, function (req, res) {
	const userEmail = req.body.userEmail;
	const friendEmail = req.body.friendEmail;
	const password = req.body.password;

	if (userEmail == undefined || friendEmail == undefined  || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}

	// var sqlCheckUser = mysql.format("select password, salt, id from users where email=?", [userEmail]);
	const sqlCheckUser = {
		// give the query a unique name
		name: 'fetch-user5',
		text: 'select password, salt, id from users where email=$1',
		values: [userEmail],
	}
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.rows.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckUser.rows[0].salt);
		if (resultsCheckUser.rows[0].password === key) {
			// var sqlCheckReceiverExist = mysql.format("select password, salt, id from users where email=?", [friendEmail]);
			const sqlCheckReceiverExist = {
				// give the query a unique name
				name: 'fetch-friends2',
				text: 'select password, salt, id from users where email=$1',
				values: [friendEmail],
			}
			con.query(sqlCheckReceiverExist, function (errorCheckReceiverExist, resultsCheckReceiverExist, fields) {
				if (errorCheckReceiverExist) throw errorCheckReceiverExist;
				if (resultsCheckReceiverExist.rows.length == 0) {
					return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
				}
				// Sender credentials are correct and both sender and receiver exists
				// var sqlFriendMessages = mysql.format("SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = ? and m.message_to_id = ?) or (m.message_from_id = ? and m.message_to_id = ?) order by m.sent_at", [resultsCheckUser[0].id, resultsCheckReceiverExist[0].id, resultsCheckReceiverExist[0].id, resultsCheckUser[0].id]);
				const sqlFriendMessages = {
					// give the query a unique name
					name: 'fetch-friend-messages',
					text: 'SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = $1 and m.message_to_id = $2) or (m.message_from_id = $3 and m.message_to_id = $4) order by m.sent_at',
					values: [resultsCheckUser.rows[0].id, resultsCheckReceiverExist.rows[0].id, resultsCheckReceiverExist.rows[0].id, resultsCheckUser.rows[0].id],
				}
				con.query(sqlFriendMessages, function (errorFriendMessages, resultsFriendMessages, fields) {
					if (errorFriendMessages) throw errorFriendMessages;
					var messages = [];
					for (i=0; i<resultsFriendMessages.rows.length; i++) {
						if (resultsFriendMessages.rows[i].message_from_id == resultsCheckUser.rows[0].id && resultsFriendMessages.rows[i].message_to_id == resultsCheckReceiverExist.rows[0].id) {
							messages.push({"message" : resultsFriendMessages.rows[i].message, "messageFromEmail" : userEmail, "messageToEmail" : friendEmail, "sentAt": resultsFriendMessages.rows[i].sent_at});
						}
						else {
							messages.push({"message" : resultsFriendMessages.rows[i].message, "messageFromEmail" : friendEmail, "messageToEmail" : userEmail, "sentAt": resultsFriendMessages.rows[i].sent_at});
						}
					}
					return res.send({'status': 200, 'error': false, 'message': 'Messages attached', 'time': new Date(Date.now()), 'msgs' : messages} );
				});
			});
		} else {
			return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
		}
	});
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

const sendNewUserEmail  = (email, accountCreated, token) => {
	var mailOptions = {
		from: 'example@gmail.com',
		to: email,
		subject: accountCreated ? "Welcome to Friends" : "Password Reset Request",
		text: accountCreated ? "New Account Created." : "To reset your password, click the link below:\n http://localhost:3000/#/reset-password?token=" + token
	  };
	  
	  transporter.sendMail(mailOptions, function(error, info){
		if (error) {
		  console.log(error);
		}
	  });
}

const saveMessage = (message, senderID, receiverID) => {

	// var sqlGetID = mysql.format("SELECT id FROM message ORDER BY id DESC LIMIT 1");
	const sqlGetID = {
		// give the query a unique name
		name: 'fetch-id2',
		text: 'SELECT id FROM message ORDER BY id DESC LIMIT 1',
	}
	con.query(sqlGetID, function (errorGetID, resultsGetID, fields) {
		if (errorGetID) throw errorGetID;
		var newMessageID = 1;

		if (resultsGetID.rows.length > 0) {
			newMessageID = resultsGetID.rows[0].id + 1;
		}
		// var sqlSendMessage = mysql.format("insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES (?, ?, ?, ?, ?)", [newMessageID, message, new Date(Date.now()), senderID, receiverID]);
		const sqlSendMessage = {
			// give the query a unique name
			name: 'fetch-insert-message',
			text: 'insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES ($1, $2, $3, $4, $5)',
			values: [newMessageID, message, new Date(Date.now()), senderID, receiverID],
		}
		con.query(sqlSendMessage, function (errorSendMessage, resultsSendMessage, fields) {
			if (errorSendMessage) throw errorSendMessage;
		});
	});

}



// start the app
app.listen(8080)
