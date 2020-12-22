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
var mongo = require('mongodb');

var con = mysql.createConnection({
  host: "localhost",
  user: "ghumman",
  password: "ghumman",
  database : "friends_mysql"
});

var MongoClient = mongo.MongoClient;
var url = "mongodb://localhost:27017/";

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

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result == null) {
				// User does not exist, creat a new one
				// generate salt
				const salt = generateSalt();
				// generate database password
				const dbPassword = generateKey(password, salt);

				var newUser = { firstName: firstName, lastName: lastName, createdAt:  new Date(Date.now()), salt: salt, password: dbPassword, email: email };
				dbo.collection("user").insertOne(newUser, function(err, result) {
					if (err) throw err;
					// sendNewUserEmail(email, true, "")
					return res.send({'status': 200, 'error': false, 'message': 'User Created', 'time': new Date(Date.now())} );
				})

			} else {
				return res.send({'status': 400, 'error': true, 'message': 'User Already Exists', 'time': new Date(Date.now())} );
			}
			  
			db.close();
		});
	});
})

// login
app.post('/login', urlencodedParser, function (req, res) {

	const email = req.body.email;
	const password = req.body.password;

	if (email == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Email or Password missing', 'time': new Date(Date.now())} );
	}
	// var sql = mysql.format("select password, salt from user where email=?", [email]);
	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result.length == 0) {
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}

			const key = generateKey(password, result.salt);
			if (result.password === key) {
				return res.send({'status': 200, 'error': false, 'message': 'Logged In', 'time': new Date(Date.now())} );
			} else {
				return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
			}
			  
			db.close();
		});
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

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result != null) {

				const key = generateKey(password, result.salt);
				if (result.password === key) {

					// generate salt
					const salt = generateSalt();
					// generate database password
					const dbPassword = generateKey(newPassword, salt);

					var query = { email:  email };
					var newValues = { $set: {salt: salt, password:  dbPassword } };
					dbo.collection("user").updateOne(query, newValues, function(err, result) {
						if (err) throw err;
						return res.send({'status': 200, 'error': false, 'message': 'Password changed', 'time': new Date(Date.now())} );
					})

				} else {
					return res.send({'status': 400, 'error': true, 'message': 'Original password not right', 'time': new Date(Date.now())} );
				}

			} else {
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}
			  
			db.close();
		});
	});
})

// forgot-password
app.post('/forgot-password', urlencodedParser, function (req, res) {
	const email = req.body.email;

	if (email == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Email is required', 'time': new Date(Date.now())} );
	}

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result != null) {

				var uuidNumber = uuidv4();

				var query = { email:  email };
				var newValues = { $set: {resetToken: uuidNumber } };
				dbo.collection("user").updateOne(query, newValues, function(err, result) {
					if (err) throw err;
					// sendNewUserEmail(email, false, uuidNumber)
					return res.send({'status': 200, 'error': false, 'message': 'Reset password is sent', 'time': new Date(Date.now())} );
				})

			} else {
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}
			  
			db.close();
		});
	});

	// var sqlLogin = mysql.format("select password, salt from user where email=?", [email]);
	// con.query(sqlLogin, function (errorLogin, resultsLogin, fields) {
	// 	if (errorLogin) throw errorLogin;
	// 	if (resultsLogin.length == 0) {
	// 		return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
	// 	}

	// 	var uuidNumber = uuidv4();
	// 	var sqlUpdateUser = mysql.format("UPDATE user SET reset_token=? where email=?", [uuidNumber, email]);
	// 	con.query(sqlUpdateUser, function (errorUpdateUser, resultsUpdateUser, fields) {
	// 		if (errorUpdateUser) throw errorUpdateUser;

	// 		sendNewUserEmail(email, false, uuidNumber)
	// 		return res.send({'status': 200, 'error': false, 'message': 'Reset password is sent', 'time': new Date(Date.now())} );
	// 	});
	// });

})

// reset-password
app.post('/reset-password', urlencodedParser, function (req, res) {
	const token = req.body.token;
	const password = req.body.password;

	if (token == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({resetToken : token}, function(err, result) {
			if (err) throw err;
			if (result != null) {

				// generate salt
				const salt = generateSalt();
				// generate database password
				const dbPassword = generateKey(password, salt);

				var query = { email:  result.email };
				var newValues = { $set: {salt: salt, password:  dbPassword, resetToken: null } };
				dbo.collection("user").updateOne(query, newValues, function(err, result2) {
					if (err) throw err;
					return res.send({'status': 200, 'error': false, 'message': 'Password successfully reset', 'time': new Date(Date.now())} );
				})

			} else {
				return res.send({'status': 400, 'error': true, 'message': 'Token is not valid', 'time': new Date(Date.now())} );
			}
			  
			db.close();
		});
	});


	// var sqlCheckToken = mysql.format("select id from user where reset_token=?", [token]);
	// con.query(sqlCheckToken, function (errorCheckToken, resultsCheckToken, fields) {
	// 	if (errorCheckToken) throw errorCheckToken;
	// 	if (resultsCheckToken.length == 0) {
	// 		return res.send({'status': 400, 'error': true, 'message': 'Token is not valid', 'time': new Date(Date.now())} );
	// 	} else {
	// 		// generate salt
	// 		const salt = generateSalt();
	// 		// generate database password
	// 		const dbPassword = generateKey(password, salt);
	// 		// update user using salt, password, token using id
	// 		var sqlUpdateUser = mysql.format("UPDATE user SET salt=?, password=?, reset_token=? where id=?", [salt, dbPassword, null, resultsCheckToken[0].id]);
	// 		con.query(sqlUpdateUser, function (errorCreateUser, resultsCreateUser, fields) {
	// 			if (errorCreateUser) throw errorCreateUser;

	// 			return res.send({'status': 200, 'error': false, 'message': 'Password successfully reset', 'time': new Date(Date.now())} );
	// 		});
	// 	}
	// });
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
	const message = req.body.message;
	const messageFromEmail = req.body.messageFromEmail;
	const messageToEmail = req.body.messageToEmail;
	const password = req.body.password;

	if (message == undefined || messageFromEmail == undefined || messageToEmail == undefined || password == undefined) {
		return res.send({'status': 400, 'error': true, 'message': 'Required data missing', 'time': new Date(Date.now())} );
	}
	var sqlCheckSender = mysql.format("select password, salt, id from user where email=?", [messageFromEmail]);
	con.query(sqlCheckSender, function (errorCheckSender, resultsCheckSender, fields) {
		if (errorCheckSender) throw errorCheckSender;
		if (resultsCheckSender.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckSender[0].salt);
		if (resultsCheckSender[0].password === key) {
			var sqlCheckReceiverExist = mysql.format("select password, salt, id from user where email=?", [messageToEmail]);
			con.query(sqlCheckReceiverExist, function (errorCheckReceiverExist, resultsCheckReceiverExist, fields) {
				if (errorCheckReceiverExist) throw errorCheckReceiverExist;
				if (resultsCheckReceiverExist.length == 0) {
					return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
				}
				// Sender credentials are correct and both sender and receiver exists
				saveMessage(message, resultsCheckSender[0].id, resultsCheckReceiverExist[0].id);
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

	var sqlCheckUser = mysql.format("select password, salt, id from user where email=?", [userEmail]);
	con.query(sqlCheckUser, function (errorCheckUser, resultsCheckUser, fields) {
		if (errorCheckUser) throw errorCheckUser;
		if (resultsCheckUser.length == 0) {
			return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
		}

		const key = generateKey(password, resultsCheckUser[0].salt);
		if (resultsCheckUser[0].password === key) {
			var sqlCheckReceiverExist = mysql.format("select password, salt, id from user where email=?", [friendEmail]);
			con.query(sqlCheckReceiverExist, function (errorCheckReceiverExist, resultsCheckReceiverExist, fields) {
				if (errorCheckReceiverExist) throw errorCheckReceiverExist;
				if (resultsCheckReceiverExist.length == 0) {
					return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
				}
				// Sender credentials are correct and both sender and receiver exists
				var sqlFriendMessages = mysql.format("SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = ? and m.message_to_id = ?) or (m.message_from_id = ? and m.message_to_id = ?) order by m.sent_at", [resultsCheckUser[0].id, resultsCheckReceiverExist[0].id, resultsCheckReceiverExist[0].id, resultsCheckUser[0].id]);
				con.query(sqlFriendMessages, function (errorFriendMessages, resultsFriendMessages, fields) {
					if (errorFriendMessages) throw errorFriendMessages;
					var messages = [];
					for (i=0; i<resultsFriendMessages.length; i++) {
						if (resultsFriendMessages[i].message_from_id == resultsCheckUser[0].id && resultsFriendMessages[i].message_to_id == resultsCheckReceiverExist[0].id) {
							messages.push({"message" : resultsFriendMessages[i].message, "messageFromEmail" : userEmail, "messageToEmail" : friendEmail, "sentAt": resultsFriendMessages[i].sent_at});
						}
						else {
							messages.push({"message" : resultsFriendMessages[i].message, "messageFromEmail" : friendEmail, "messageToEmail" : userEmail, "sentAt": resultsFriendMessages[i].sent_at});
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

	var sqlGetID = mysql.format("SELECT id FROM message ORDER BY id DESC LIMIT 1");
	con.query(sqlGetID, function (errorGetID, resultsGetID, fields) {
		if (errorGetID) throw errorGetID;
		var newMessageID = 1;

		if (resultsGetID.length > 0) {
			newMessageID = resultsGetID[0].id + 1;
		}
		var sqlSendMessage = mysql.format("insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES (?, ?, ?, ?, ?)", [newMessageID, message, new Date(Date.now()), senderID, receiverID]);
		con.query(sqlSendMessage, function (errorSendMessage, resultsSendMessage, fields) {
			if (errorSendMessage) throw errorSendMessage;
		});
	});

}



// start the app
app.listen(8080)
