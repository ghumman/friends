package com.friends.message.User;

import java.sql.Timestamp;

import lombok.Data;

@Data
public class UserResponse {
    private Timestamp time; 
    private int status;
    private String error; 
    private String message; 
}