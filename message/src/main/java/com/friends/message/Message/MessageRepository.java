package com.friends.message.Message;

import java.util.List;

import com.friends.message.User.User;

import org.springframework.data.repository.CrudRepository;
public interface MessageRepository extends CrudRepository<Message, Integer>{
    List<Message> findByMessageFrom (User messageFrom);
    List<Message> findByMessageTo (User messageTo);
}