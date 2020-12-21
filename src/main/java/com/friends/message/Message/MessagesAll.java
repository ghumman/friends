package com.friends.message.Message;

import java.util.Date;

import lombok.Data;

@Data
public class MessagesAll {
    private String message; 
    private String messageFromEmail;
    private String messageToEmail;
    private Date sentAt;
}