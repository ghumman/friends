## About Friends
Friends is a collection of fullstack projects created in different programming languages which includes backend applications, frontend applications, different databases and frameworks. Currently it has following projects. 

### Backend Applications
- [Java Spring Boot with MySQL, Hibernate, JPA, Swagger](https://github.com/ghumman/friends/tree/java-spring-boot).
- [Java Spring Boot with MongoDB, JPA, Swagger](https://github.com/ghumman/friends/tree/java-spring-boot-mongo).
- [Java Spring Boot with PostgreSQL, JPA, Swagger](https://github.com/ghumman/friends/tree/java-spring-boot-psql).
- [Docker Java(Spring Boot, SQL Lite) + React](https://github.com/ghumman/friends/tree/java-react-sqllite-docker-compose).
- [Golang with MySQL](https://github.com/ghumman/friends/tree/go).
- [Golang with MongoDB](https://github.com/ghumman/friends/tree/go-mongo).
- [Golang with PostgreSQL](https://github.com/ghumman/friends/tree/go-psql).
- [Python Flask with MySQL](https://github.com/ghumman/friends/tree/python).
- [Python Flask with MongoDB](https://github.com/ghumman/friends/tree/python-mongo).
- [Python Flask with PostgreSQL](https://github.com/ghumman/friends/tree/python-psql).
- [Python Django with MySQL](https://github.com/ghumman/friends/tree/python-django).
- [Node.js with MySQL](https://github.com/ghumman/friends/tree/node-js).
- [Node.js with MongoDB](https://github.com/ghumman/friends/tree/node-js-mongo).
- [Node.js with PostgreSQL](https://github.com/ghumman/friends/tree/node-js-psql).
- [Ruby with MySQL](https://github.com/ghumman/friends/tree/ruby).
- [Ruby with MongoDB](https://github.com/ghumman/friends/tree/ruby-mongo).
- [Ruby with PostgreSQL](https://github.com/ghumman/friends/tree/ruby-psql).
- [PHP with MySQL](https://github.com/ghumman/friends/tree/php).
- [PHP with MongoDB](https://github.com/ghumman/friends/tree/php-mongo).
- [PHP with PostgreSQL](https://github.com/ghumman/friends/tree/php-psql).
- [PHP Laravel with MySQL](https://github.com/ghumman/friends/tree/php-laravel).

### Frontend Applications
#### React 
- [React.js with Functional Components, Redux](https://github.com/ghumman/friends/tree/react-js).

- [React Website](https://ghumman.github.io/friends)

#### Angular
- [Angular.js 11](https://github.com/ghumman/friends/tree/angular).

- [Angular Website](https://ghumman.github.io/friends-angular-ui)

#### Vue
- [Vue.js with Vite](https://github.com/ghumman/friends/tree/vue-js).

## Features
Following are the main features which are provided in every project and can also be found here [FEATURES](https://github.com/ghumman/friends/blob/project-documentation/features.txt).

```
Authentication: 

	Create User and Sending Email
		add-user: 
			authType
			firstName
			lastName
			password
			email
			token

	Login
		login	
			authType
			email
			password
			token

	Change Password
		change-password
			authType
			email
			password
			newPassword
			token
			
	Forgot Password and Sending Email
		forgot-password
			authType
			email

	Reset Password using Reset Email Link
		reset-password
			authType
			token
			password

	Show All Friends
		all-friends
			authType
			email
			password
			token
			


Messages: 

	Send Message
		send-message
			authType
			message	
			password
			messageFromEmail
			messageToEmail
			token

	Show Conversation Between User and Friend
		messages-user-and-friend
			authType
			password
			userEmail
			friendEmail
			token
			
```

## Github Project
Project timeline and tickets' statuses can be found here. 
[Github Project](https://github.com/users/ghumman/projects/1).

## Request -> Body -> Variables
### add-user (addUser)
```
firstName, lastName, email, password, authType, token
```
### login (login)
```
email, password, authType, token
```
### change-password (changePassword)
```
email, password, newPassword, authType
```
### forgot-password (forgotPassword)
```
email
```
### reset-password (resetPassword)
```
token, password
```
### all-friends (allFriends)
```
email, password, authType, token
```
### send-message (sendMessage)
```
message, messageFromEmail, messageToEmail, authType, password, token
```
### messages-user-and-friend (messagesUserAndFriend or messagesBySenderToReceiver)
```
userEmail, friendEmail, authType, password, token
```
## Curling all the endpoints
### add-user (addUser)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"firstName":"John","lastName":"Smith", "email":"jsmith@gmail.com", "password": "smIth@2020", "authType": "regular", "token": ""}' \
  http://localhost:8080/add-user
```

### login (login)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"jsmith@gmail.com", "password": "smIth@2020", "authType": "regular", "token": ""}' \
  http://localhost:8080/login
```

### change-password (changePassword)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"jsmith@gmail.com", "password": "smIth@2020", "newPassword": "smIth@2021", "authType": "regular"}' \
  http://localhost:8080/change-password
```

### forgot-password (forgotPassword)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"jsmith@gmail.com"}' \
  http://localhost:8080/forgot-password
```

### reset-password (resetPassword)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"token":"xxxxxxxxxxx", "password": "smIth@2019"}' \
  http://localhost:8080/reset-password
```

### all-friends (allFriends)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"email":"jsmith@gmail.com", "password": "smIth@2020", "authType": "regular", "token": ""}' \
  http://localhost:8080/all-friends
```

### send-message (sendMessage)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"message":"Hello Alex, How are you?", "messageFromEmail": "jsmith@gmail.com", "messageToEmail": "alexgreen@gmail.com", "authType": "regular", "password": "smIth@2020", "token": ""}' \
  http://localhost:8080/send-message
```

### messages-user-and-friend (messageUserAndFriend or messagesBySenderToReceiver)
```
curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"userEmail":"jsmith@gmail.com", "friendEmail": "alexgreen@gmail.com", "authType": "regular", "password": "smIth@2020", "token": ""}' \
  http://localhost:8080/messages-user-and-friend
```

## Contribution - PRs welcome
Feel free to customize and refactor current code or add new languages/frameworks. 
