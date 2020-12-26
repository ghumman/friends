package com.friends.message.Users;

import java.util.Date;

import java.util.List;

import lombok.Data;

@Data
public class UserResponse {
    private Date time;
    private int status;
    private Boolean error; 
    private String message; 
    private List<UsersAll> usersAll; 
}