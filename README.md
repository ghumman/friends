This branch contains java spring boot back end application
### Push code using following command
git push friends master:java-spring-boot
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
