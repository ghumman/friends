This branch includes project documentation.
### Push code using following command
```
git push friends project-documentation
or 
git push friends master:project-documentation
```

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

