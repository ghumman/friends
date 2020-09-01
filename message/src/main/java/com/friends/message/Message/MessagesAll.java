package com.friends.message.Message;

import java.sql.Timestamp;

import lombok.Data;

@Data
public class MessagesAll {
    private String message; 
    private String messageFromEmail;
    private String messageToEmail;
    private Timestamp sentAt;
}