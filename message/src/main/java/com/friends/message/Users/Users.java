package com.friends.message.Users;

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
public class Users {

    public enum AuthType {regular, special};

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id; 
    
    @NotNull
    private String firstName; 
    @NotNull
    private String lastName; 

    @NotNull
    private Date createdAt;

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

    Users(String fName, String lName, String newEmail , String pwd, String auth, String tkn) {

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
        createdAt = new Date();;

    }

    Users() {}

}