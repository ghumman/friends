package com.friends.message.User;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

import lombok.Data;

@Entity
@Data
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer id; 
    
    private String firstName; 
    private String lastName; 

    private String password; 
    private String email; 
    
    private String authType; 

}