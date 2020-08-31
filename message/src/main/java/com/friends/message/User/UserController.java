package com.friends.message.User;

import org.springframework.web.bind.annotation.RestController;

import java.sql.Timestamp;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;


@RestController
public class UserController {

    public enum AuthType {regular, special};

    private UserResponse userResponse; 

    @Autowired 
    private UserRepository userRepository;

    @PostMapping("/add-user")
    public ResponseEntity<UserResponse> addUser(
        @RequestParam String firstName, 
        @RequestParam String lastName,
        @RequestParam String email, 
        String password, 
        @RequestParam String authType, 
        String token
    ) {

        userResponse = new UserResponse(); 
        userResponse.setMessage("");
        userResponse.setError(false);

        // validation
        if (firstName.equals("")) {
            userResponse.setError(true);
            userResponse.setMessage("First Name can not be empty. ");
        } 

        if (lastName.equals("")) {
            userResponse.setError(true);
            String temp = userResponse.getMessage(); 
            userResponse.setMessage(temp + "Last Name can not be empty. ");
        } 

        if (email.equals("")) {
            userResponse.setError(true);
            String temp = userResponse.getMessage(); 
            userResponse.setMessage(temp + "Email can not be empty. ");
        } 

        if (!authType.equals(AuthType.regular.toString()) && !authType.equals(AuthType.special.toString())) {
            userResponse.setError(true);
            String temp = userResponse.getMessage(); 
            userResponse.setMessage(temp + "Authentication can only be regular or special. ");
        }

        if (authType.equals(AuthType.regular.toString()) && password.equals("")) {
            userResponse.setError(true);
            String temp = userResponse.getMessage(); 
            userResponse.setMessage(temp + "Password can not be empty. ");
        }

        if (authType.equals(AuthType.special.toString()) && token.equals("")) {
            userResponse.setError(true);
            String temp = userResponse.getMessage(); 
            userResponse.setMessage(temp + "Token can not be empty. ");
        }

        if (userResponse.getError()) {
            userResponse.setStatus(400);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }

        // input is valid 

        if (userRepository.findByEmail(email) == null) {

            User newUser = new User(firstName, lastName, email, password, authType, token); 
            userRepository.save(newUser); 

            userResponse.setMessage("User Created");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }

        else {
            userResponse.setMessage("User Already Exists");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }
    }



    @PostMapping("/login")
    public ResponseEntity<UserResponse> login(
        @RequestParam String email, 
        String password, 
        @RequestParam String authType, 
        String token
    ) {

        try {

            userResponse = new UserResponse(); 

            User currentUser = userRepository.findByEmail(email);

            if (currentUser == null) {

                userResponse.setMessage("User Does Not Exist");
                userResponse.setStatus(400);
                userResponse.setError(true);
                userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(userResponse);
            }

            else {

                if (authType.equals(AuthType.regular.toString())) {

                    String salt = currentUser.getSalt(); 
                    String securedPassword = currentUser.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        userResponse.setMessage("Logged In");
                        userResponse.setStatus(200);
                        userResponse.setError(false);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    }
            
                } else if (authType.equals(AuthType.special.toString())) {
                    // login using OAuth
                    boolean passwordMatch = token.equals(currentUser.getToken());
            
                    if(passwordMatch) 
                    {
                        userResponse.setMessage("Logged In");
                        userResponse.setStatus(200);
                        userResponse.setError(false);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    }
                } else {
                    userResponse.setMessage("Authentication Type not right");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                    return ResponseEntity.ok(userResponse);
                }

            }
        } catch(Exception e) {
            userResponse = new UserResponse(); 
            userResponse.setMessage(e.toString());
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }
    }

    @PostMapping("/delete")
    public ResponseEntity<UserResponse> delete(
        @RequestParam String email, 
        String password, 
        @RequestParam String authType, 
        String token
    ) {

        try {

            userResponse = new UserResponse(); 

            User currentUser = userRepository.findByEmail(email);

            if (currentUser == null) {

                userResponse.setMessage("User Does Not Exist");
                userResponse.setStatus(400);
                userResponse.setError(true);
                userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                return ResponseEntity.ok(userResponse);
            }

            else {

                if (authType.equals(AuthType.regular.toString())) {

                    String salt = currentUser.getSalt(); 
                    String securedPassword = currentUser.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        userRepository.delete(currentUser);
                        userResponse.setMessage("User Deleted");
                        userResponse.setStatus(200);
                        userResponse.setError(false);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    }
            
                } else if (authType.equals(AuthType.special.toString())) {
                    // login using OAuth
                    boolean passwordMatch = token.equals(currentUser.getToken());
            
                    if(passwordMatch) 
                    {
                        userRepository.delete(currentUser);
                        userResponse.setMessage("User Deleted");
                        userResponse.setStatus(200);
                        userResponse.setError(false);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    }
                } else {
                    userResponse.setMessage("Authentication Type not right");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                    return ResponseEntity.ok(userResponse);
                }

            }
        } catch(Exception e) {
            userResponse = new UserResponse(); 
            userResponse.setMessage(e.toString());
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }
    }


    @GetMapping(path="/all")
    public @ResponseBody Iterable<User> getAllUsers() {
      // This returns a JSON or XML with the users
      return userRepository.findAll();
    }
}