package com.friends.message.User;

import org.springframework.web.bind.annotation.RestController;

import java.sql.Timestamp;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;


@RestController
public class UserController {

    private UserResponse userResponse; 
    @Autowired 
    private UserRepository userRepository;

    @PostMapping("/add-user")
    public ResponseEntity<UserResponse> addUser(
        @RequestParam String firstName, 
        @RequestParam String lastName,
        @RequestParam String email, 
        @RequestParam String password, 
        @RequestParam String authType
    ) {


        User newUser = new User(); 
        newUser.setFirstName(firstName);
        newUser.setLastName(lastName);
        newUser.setEmail(email);
        newUser.setPassword(password);
        newUser.setAuthType(authType);


        userRepository.save(newUser); 

        userResponse = new UserResponse(); 
        userResponse.setMessage("User Created");
        userResponse.setStatus(200);
        userResponse.setError(null);
        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
        return ResponseEntity.ok(userResponse);
    }

    @GetMapping(path="/all")
    public @ResponseBody Iterable<User> getAllUsers() {
      // This returns a JSON or XML with the users
      return userRepository.findAll();
    }
}