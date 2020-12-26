package com.friends.message.Message;

import java.util.Date;
import java.util.List;

import lombok.Data;

@Data
public class MessageResponse {
    private Date time; 
    private int status;
    private Boolean error; 
    private String message; 
    private List<MessagesAll> msgs; 
}