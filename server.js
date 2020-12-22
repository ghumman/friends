var express = require('express');
var app = express();
var bodyParser = require('body-parser');
var fs = require("fs");
const crypto = require('crypto');
const { v4: uuidv4 } = require('uuid');
var nodemailer = require('nodemailer');

// create application/x-www-form-urlencoded parser
var urlencodedParser = bodyParser.urlencoded({ extended: false });

var mongo = require('mongodb');

var MongoClient = mongo.MongoClient;
var url = "mongodb://localhost:27017/";

var transporter = nodemailer.createTransport({
	service: 'gmail',
	auth: {
	  user: 'example@gmail.com',
	  pass: 'xxxxxx'
	}
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
					db.close();
					return res.send({'status': 200, 'error': false, 'message': 'User Created', 'time': new Date(Date.now())} );
				})

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'User Already Exists', 'time': new Date(Date.now())} );
			}
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
	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result.length == 0) {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}

			const key = generateKey(password, result.salt);
			if (result.password === key) {
				db.close();
				return res.send({'status': 200, 'error': false, 'message': 'Logged In', 'time': new Date(Date.now())} );
			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
			}
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
						db.close();
						return res.send({'status': 200, 'error': false, 'message': 'Password changed', 'time': new Date(Date.now())} );
					})

				} else {
					db.close();
					return res.send({'status': 400, 'error': true, 'message': 'Original password not right', 'time': new Date(Date.now())} );
				}

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}
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
					db.close();
					return res.send({'status': 200, 'error': false, 'message': 'Reset password is sent', 'time': new Date(Date.now())} );
				})

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}
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
					db.close();
					return res.send({'status': 200, 'error': false, 'message': 'Password successfully reset', 'time': new Date(Date.now())} );
				})

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'Token is not valid', 'time': new Date(Date.now())} );
			}
		});
	});
})

// all-friends
app.post('/all-friends', urlencodedParser, function (req, res) {
	const email = req.body.email;
	const password = req.body.password;

	if (email == undefined || password == undefined ) {
		return res.send({'status': 400, 'error': true, 'message': 'Email or Password missing', 'time': new Date(Date.now())} );
	}

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : email}, function(err, result) {
			if (err) throw err;
			if (result != null) {

				const key = generateKey(password, result.salt);
				if (result.password === key) {

					dbo.collection("user").find({email : {$ne : email}}).toArray(function(err, friends) {
						if (err) throw err;
						var users = [];
						for(let friend of friends) {
							users.push({"firstName" : friend.firstName, "lastName" : friend.lastName, "email" : friend.email})
						}
						db.close();
						return res.send({'status': 200, 'error': false, 'message': 'Friends attached', 'time': new Date(Date.now()), 'usersAll' : users} );
					})
				} else {
					db.close();
					return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
				}

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'User Does Not Exist', 'time': new Date(Date.now())} );
			}
		});
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

	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : messageFromEmail}, function(err, resultSender) {
			if (err) throw err;
			if (resultSender != null) {

				const key = generateKey(password, resultSender.salt);
				if (resultSender.password === key) {


					return dbo.collection("user").findOne({email : messageToEmail}, function(err, resultReceiver) {
						if (err) throw err;
						if (resultReceiver != null) {
							var newMessage = { message: message, sentAt: new Date(Date.now()), messageFrom:  resultSender, messageTo: resultReceiver };

							return dbo.collection("message").insertOne(newMessage, function(err, resultNotUsed) {
								if (err) throw err;
								db.close();
								return res.send({'status': 200, 'error': false, 'message': 'Message sent', 'time': new Date(Date.now())} );
							});
						} else {
							db.close();
							return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
						}
					});
					

				} else {
					db.close();
					return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
				}

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
			}
		});
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


	MongoClient.connect(url, function(err, db) {
		if (err) throw err;
		var dbo = db.db("friends_mongo");
		dbo.collection("user").findOne({email : userEmail}, function(err, resultSender) {
			if (err) throw err;
			if (resultSender != null) {

				const key = generateKey(password, resultSender.salt);
				if (resultSender.password === key) {

					return dbo.collection("user").findOne({email : friendEmail}, function(err, resultReceiver) {
						if (err) throw err;
						if (resultReceiver != null) {
							return dbo.collection("message").find({$or : [{messageFrom : resultSender, messageTo : resultReceiver}, {messageFrom : resultReceiver, messageTo : resultSender}]}).toArray(function(err, conversations) {
								if (err) throw err;
								var messages = [];
								for(let conversation of conversations) {
									messages.push({"message" : conversation.message, "messageFromEmail" : conversation.messageFrom.email, "messageToEmail" : conversation.messageTo.email, "sentAt": conversation.sentAt})
								}
								db.close();
								return res.send({'status': 200, 'error': false, 'message': 'Messages attached', 'time': new Date(Date.now()), 'msgs' : messages} );
							})
						} else {
							db.close();
							return res.send({'status': 400, 'error': true, 'message': 'Receiver Does Not Exist', 'time': new Date(Date.now())} );
						}
					});

				} else {
					db.close();
					return res.send({'status': 400, 'error': true, 'message': 'Login Failed', 'time': new Date(Date.now())} );
				}

			} else {
				db.close();
				return res.send({'status': 400, 'error': true, 'message': 'Sender Does Not Exist', 'time': new Date(Date.now())} );
			}
		});
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

// start the app
app.listen(8080)
