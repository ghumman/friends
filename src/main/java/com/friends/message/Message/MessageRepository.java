package com.friends.message.Message;

import java.util.List;

import com.friends.message.User.User;

import org.springframework.data.mongodb.repository.MongoRepository;
// import org.springframework.data.mongodb.repository.Query;

import org.springframework.stereotype.Repository;

@Repository
public interface MessageRepository extends MongoRepository<Message, Integer>{
    List<Message> findAllByMessageFrom (User messageFrom);
    List<Message> findAllByMessageTo (User messageTo);

    // @Query("SELECT m FROM Message m WHERE (m.messageFrom = ?1 and m.messageTo = ?2) or (m.messageFrom = ?2 and m.messageTo = ?1) order by m.sentAt")
    // List<Message> findAllWithCustomQuery(User user, User friend);
    // List<Message> findAll();

    // @Query("SELECT m FROM Message m WHERE (m.messageFrom = ?1 and m.messageTo = ?2) or (m.messageFrom = ?2 and m.messageTo = ?1) order by m.sentAt")
    // List<Message> findAllWithCustomQuery(User user, User friend);
    List<Message> findAllByMessageToAndMessageFromOrMessageToAndMessageFrom(User to, User from, User to2, User from2);


}