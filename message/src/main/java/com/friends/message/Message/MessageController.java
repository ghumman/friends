package com.friends.message.Message;

import com.friends.message.User.PasswordUtils;
import com.friends.message.User.User;
import com.friends.message.User.UserRepository;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@CrossOrigin(origins = "http://localhost:3000")
public class MessageController {
    @Autowired
    private UserRepository userRepository; 

    @Autowired
    private MessageRepository messageRepository; 

    public enum AuthType {regular, special};

    private MessageResponse messageResponse; 

    @PostMapping("/send-message")
    public  ResponseEntity<MessageResponse> sendMessage(
        @RequestParam String message, 
        @RequestParam String messageFromEmail,
        @RequestParam String messageToEmail,
        @RequestParam String authType,
        String password, 
        String token
    ) {
        try {

            messageResponse = new MessageResponse(); 

            User sendUser = userRepository.findByEmail(messageFromEmail);
            User receiveUser = userRepository.findByEmail(messageToEmail);

            if (sendUser == null) {

                messageResponse.setMessage("Sender Does Not Exist");
                messageResponse.setStatus(400);
                messageResponse.setError(true);
                messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(messageResponse);
            } 
            else if (receiveUser == null) {
                messageResponse.setMessage("Receiver Does Not Exist");
                messageResponse.setStatus(400);
                messageResponse.setError(true);
                messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(messageResponse);
            }
            else {

                if (authType.equals(AuthType.regular.toString())) {

                    String salt = sendUser.getSalt(); 
                    String securedPassword = sendUser.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        saveMessage(message, sendUser, receiveUser); 
                        messageResponse.setMessage("Message sent");
                        messageResponse.setStatus(200);
                        messageResponse.setError(false);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    } else {
                        messageResponse.setMessage("Login Failed");
                        messageResponse.setStatus(400);
                        messageResponse.setError(true);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    }
            
                } else if (authType.equals(AuthType.special.toString())) {
                    // login using OAuth
                    boolean passwordMatch = token.equals(sendUser.getToken());
            
                    if(passwordMatch) 
                    {
                        saveMessage(message, sendUser, receiveUser); 
                        messageResponse.setMessage("Message sent");
                        messageResponse.setStatus(200);
                        messageResponse.setError(false);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    } else {
                        messageResponse.setMessage("Login Failed");
                        messageResponse.setStatus(400);
                        messageResponse.setError(true);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    }
                } else {
                    messageResponse.setMessage("Authentication Type not right");
                    messageResponse.setStatus(400);
                    messageResponse.setError(true);
                    messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                    return ResponseEntity.ok(messageResponse);
                }

            }
        } catch(Exception e) {
            messageResponse = new MessageResponse(); 
            messageResponse.setMessage(e.toString());
            messageResponse.setStatus(400);
            messageResponse.setError(true);
            messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(messageResponse);
        }

    }

    public void saveMessage (String message, User sendUser, User receiveUser) {
        Message msg = new Message(message, sendUser, receiveUser) ; 
        messageRepository.save(msg);        
    }

    @PostMapping("/messages-user-and-friend")
    public  ResponseEntity<MessageResponse> messagesBySenderToReceiver(
        @RequestParam String userEmail,
        @RequestParam String friendEmail,
        @RequestParam String authType,
        String password, 
        String token
    ) {
        try {

            messageResponse = new MessageResponse(); 

            User user = userRepository.findByEmail(userEmail);
            User friend = userRepository.findByEmail(friendEmail);

            if (user == null) {

                messageResponse.setMessage("Sender Does Not Exist");
                messageResponse.setStatus(400);
                messageResponse.setError(true);
                messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(messageResponse);
            } 
            else if (friend == null) {
                messageResponse.setMessage("Receiver Does Not Exist");
                messageResponse.setStatus(400);
                messageResponse.setError(true);
                messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(messageResponse);
            }
            else {

                if (authType.equals(AuthType.regular.toString())) {

                    String salt = user.getSalt(); 
                    String securedPassword = user.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        List<Message> msgs = messageRepository.findAllWithCustomQuery(user, friend);
                        List<MessagesAll> customMessages = new ArrayList<>(); 

                        for (Message msg : msgs) {
                            MessagesAll customMessage = new MessagesAll(); 
                            customMessage.setMessage(msg.getMessage());
                            customMessage.setMessageFromEmail(msg.getMessageFrom().getEmail());
                            customMessage.setMessageToEmail(msg.getMessageTo().getEmail());
                            customMessage.setSentAt(msg.getSentAt());

                            customMessages.add(customMessage); 
                        }

                        messageResponse.setMsgs(customMessages); 
                        messageResponse.setMessage("Messages attached");
                        messageResponse.setStatus(200);
                        messageResponse.setError(false);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    } else {
                        messageResponse.setMessage("Login Failed");
                        messageResponse.setStatus(400);
                        messageResponse.setError(true);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    }
            
                } else if (authType.equals(AuthType.special.toString())) {
                    // login using OAuth
                    boolean passwordMatch = token.equals(user.getToken());
            
                    if(passwordMatch) 
                    {
                        List<Message> msgs = messageRepository.findAllWithCustomQuery(user, friend);
                        List<MessagesAll> customMessages = new ArrayList<>(); 

                        for (Message msg : msgs) {
                            MessagesAll customMessage = new MessagesAll(); 
                            customMessage.setMessage(msg.getMessage());
                            customMessage.setMessageFromEmail(msg.getMessageFrom().getEmail());
                            customMessage.setMessageToEmail(msg.getMessageTo().getEmail());
                            customMessage.setSentAt(msg.getSentAt());

                            customMessages.add(customMessage); 
                        }

                        messageResponse.setMsgs(customMessages); 
                        messageResponse.setMessage("Messages attached");
                        messageResponse.setStatus(200);
                        messageResponse.setError(false);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    } else {
                        messageResponse.setMessage("Login Failed");
                        messageResponse.setStatus(400);
                        messageResponse.setError(true);
                        messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(messageResponse);
                    }
                } else {
                    messageResponse.setMessage("Authentication Type not right");
                    messageResponse.setStatus(400);
                    messageResponse.setError(true);
                    messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
                    return ResponseEntity.ok(messageResponse);
                }

            }
        } catch(Exception e) {
            messageResponse = new MessageResponse(); 
            messageResponse.setMessage(e.toString());
            messageResponse.setStatus(400);
            messageResponse.setError(true);
            messageResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(messageResponse);
        }
    }


}