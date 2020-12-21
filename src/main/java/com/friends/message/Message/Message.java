package com.friends.message.Message;

import java.util.Date;
import javax.validation.constraints.NotNull;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import lombok.Data;
import com.friends.message.User.User;


@Data
@Document
public class Message {

    @Id
    private String messageId;
    
    @NotNull
    private String message; 
    @NotNull
    private User messageFrom;
    @NotNull
    private User messageTo;

    @NotNull
    private Date sentAt;

    Message(String message, User messageFrom, User messageTo) {

        this.message = message; 
        this.messageFrom = messageFrom; 
        this.messageTo = messageTo; 
        Date date = new Date();
        sentAt= new Date(date.getTime()); 
    }

    Message() {}
}