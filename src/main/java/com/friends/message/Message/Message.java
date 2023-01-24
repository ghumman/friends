package com.friends.message.Message;

import java.sql.Timestamp;
import java.util.Date;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.ManyToOne;
import javax.validation.constraints.NotNull;

import com.friends.message.User.User;

import lombok.Data;

@Entity
@Data
public class Message {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Long id; 
    
    @NotNull
    private String message; 
    @NotNull
    @ManyToOne
    private User messageFrom;
    @NotNull
    @ManyToOne
    private User messageTo;

    @NotNull
    private Timestamp sentAt;

    Message(String message, User messageFrom, User messageTo) {

        this.message = message; 
        this.messageFrom = messageFrom; 
        this.messageTo = messageTo; 
        Date date = new Date();
        sentAt= new Timestamp(date.getTime()); 
    }

    Message() {}
}