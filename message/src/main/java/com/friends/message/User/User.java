package com.friends.message.User;


import java.sql.Timestamp;
import java.util.Date;

import javax.persistence.Column;
import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotNull;



import lombok.Data;

@Entity
@Data
public class User {

    public enum AuthType {regular, special};

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id; 
    
    @NotNull
    private String firstName; 
    @NotNull
    private String lastName; 

    @NotNull
    private Timestamp createdAt;

    private String salt;
    private String password; 
    private String token; 

    private String resetToken; 
    
    @Column(unique = true) 
    @NotNull
    @Email
    private String email;
    
    @NotNull
    private AuthType authType; 

    User(String fName, String lName, String newEmail , String pwd, String auth, String tkn) {

        if (auth.equals(AuthType.regular.toString())) {
            salt = PasswordUtils.getSalt(30);
            password = PasswordUtils.generateSecurePassword(pwd, salt);
            authType = AuthType.regular; 
        } else if (auth.equals(AuthType.special.toString())) {
            token = tkn; 
            authType = AuthType.special; 
        }

        firstName = fName; 
        lastName = lName; 
        email = newEmail; 
        Date date = new Date();
        createdAt = new Timestamp(date.getTime()); 

    }

    User() {}

}