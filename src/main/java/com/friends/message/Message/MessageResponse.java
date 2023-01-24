package com.friends.message.Message;

import java.sql.Timestamp;
import java.util.List;

import lombok.Data;

@Data
public class MessageResponse {
    private Timestamp time; 
    private int status;
    private Boolean error; 
    private String message; 
    private List<MessagesAll> msgs; 
}