package com.friends.message.Message;


import java.util.Date;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.ManyToOne;
import javax.validation.constraints.NotNull;

import com.friends.message.Users.Users;

import lombok.Data;

@Entity
@Data
public class Message {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id; 
    
    @NotNull
    private String message; 
    @NotNull
    @ManyToOne
    private Users messageFrom;
    @NotNull
    @ManyToOne
    private Users messageTo;

    @NotNull
    private Date sentAt;

    Message(String message, Users messageFrom, Users messageTo) {

        this.message = message; 
        this.messageFrom = messageFrom; 
        this.messageTo = messageTo;
        sentAt= new Date();
    }

    Message() {}
}