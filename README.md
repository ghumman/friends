# About 
This repo is for Java + Spring Boot + Mysql + Docker Images<br/>
There are some feature changes in this repo compared to other java projects without docker. They are mainly following.<br/>
- No email is sent in this repo.
- Sending email is removed when creating an account.
- Sending email is removed when using feature: forgot-password.
- You can just change the password when using reset-password without sending `token`. Just need to send email and password.

## If you want to run Java application only
Run following
```
cd backend
./mvnw spring-boot:run
```
Use Postman to test the application which is running on `localhost:8080`. 

## Running this compose project
In order to run the compose project, run following command.
``` 
docker compose up -d 
``` 
`-d` is to run everything in background<br/> 
 
In order to destroy everything, use following command.<br/> 
``` 
docker compose down 
``` 
If you want to remove the volumes, you will need to add the --volumes flag<br/> 


# Following documentation is from project java spring boot branch

This branch contains java spring boot back end application
### Push code using following command
git push origin java-spring-boot
### In order to allow less secure apps to send emails using gmail use following link
https://myaccount.google.com/lesssecureapps

### Starting on a new Computer
1. Make sure your VS Code has following extentions
- Java Extenstion Pack
- Spring Boot Extenstion Pack
- Lombok Annotations Support for VS Code

2. Used following guide to install Java and setup java_home
- [Install Java](https://linuxhint.com/install_jdk_14_ubuntu/)

3. Used following guide to install mysql locally and create database named "friends_mysql" and create user ghumman with password ghumman
- [Install Mysql Locally](https://www.digitalocean.com/community/tutorials/how-to-install-mysql-on-ubuntu-20-04)

4. Based on your JAVA installed version update pom.xml at the following line
```
<java.version>14</java.version>
```

### Notes
- In order go get 200 response when calling endpoints like add-user and others you need to provide gmail or other email provider username and password in application.properties file and inside application.properties in JavaMailSender method. Otherwise you'll be getting Internal Server Error as endpoints try to send emails to newly created user emails. 

- To check if swagger is working use following url 
```
http://localhost:8080/v2/api-docs
```

- To use swagger ui use following url
```
http://localhost:8080/swagger-ui/
```

- In order to run the application using command line. 
```
cd friends/message
./mvnw spring-boot:run
```
- In order to run it on VS Code, which is ready and have the prerequisites, just go to Run and Debug and run the 
Start Debugger.<br/>

- In order to create jar file/files.
```
cd message
./mvnw package
```
