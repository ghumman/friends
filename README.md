This branch contains java spring boot back end application
### Push code using following command
```
git init
git remote add java-mongo git@github.com:ghumman/friends.git
touch README.md
git add README.md
git commit -m "Initial commit"
git push java-mongo master:java-spring-boot-mongo
```
### In order to allow less secure apps to send emails using gmail use following link
https://myaccount.google.com/lesssecureapps



### Starting on a new Computer
1. Make sure your VS Code has following extentions
- Java Extenstion Pack
- Spring Boot Extenstion Pack
- Lombok Annotations Support for VS Code

2. Used following guide to install Java and setup java_home
- [Install Java](https://linuxhint.com/install_jdk_14_ubuntu/)

3. Used following guide to install mongo locally and create database named "friends_mongo" with no credentials
- [Install MongoDB Locally](https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/)

4. Based on your JAVA installed version update pom.xml at the following line
```
<java.version>14</java.version>
```

### Notes
- In order go get 200 response when calling endpoints like add-user and others you need to provide gmail or other email provider username and password in application.properties file and inside application.properties in JavaMailSender method. Otherwise you'll be getting Internal Server Error as endpoints try to send emails to newly created user emails. 

- To check if swagger is working use following url 
```
http://localhost:8102/v2/api-docs
```

- To use swagger ui use following url
```
http://localhost:8102/swagger-ui/
```
