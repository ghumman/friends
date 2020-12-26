This branch contains java spring boot back end application
### Push code using following command
```
git init
git remote add java-spring-boot-psql git@github.com:ghumman/friends.git
git add .
git commit -m "Initial commit"
git push java-spring-boot-psql master:java-spring-boot-psql
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


3. Based on your JAVA installed version update pom.xml at the following line
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
