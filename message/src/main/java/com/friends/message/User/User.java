package com.friends.message.User;

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
    private Integer id; 
    
    @NotNull
    private String firstName; 
    @NotNull
    private String lastName; 

    private String password; 
    private String token; 
    
    @Column(unique = true) 
    @NotNull
    @Email
    private String email;
    
    @NotNull
    private AuthType authType; 

    User(String fName, String lName, String newEmail , String pwd, String auth, String tkn) {

        if (auth.equals(AuthType.regular.toString())) {
            firstName = fName; 
            lastName = lName; 
            password = pwd; 
            email = newEmail; 
            authType = AuthType.regular; 
        } else if (auth.equals(AuthType.special.toString())) {
            firstName = fName; 
            lastName = lName; 
            token = tkn; 
            email = newEmail; 
            authType = AuthType.special; 
        }
    }

    User() {
        
    }

}