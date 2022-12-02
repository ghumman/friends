package com.friends.message.User;

import org.springframework.web.bind.annotation.RestController;

import java.sql.Timestamp;
import java.util.ArrayList;
import java.util.List;

import javax.servlet.http.HttpServletRequest;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;


@RestController
@CrossOrigin(origins = "http://localhost:3000")
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

    @PostMapping("/change-password")
    public ResponseEntity<UserResponse> changePassword(
        @RequestParam String email, 
        String password, 
        String newPassword,
        @RequestParam String authType
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

                if (authType.equals(AuthType.special.toString())) {

                    userResponse.setMessage("Logged in using OAuth");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                    return ResponseEntity.ok(userResponse);
                } 

                else if (authType.equals(AuthType.regular.toString())) {

                    String salt = currentUser.getSalt(); 
                    String securedPassword = currentUser.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        String newSalt = PasswordUtils.getSalt(30);
                        String newPwd = PasswordUtils.generateSecurePassword(newPassword, newSalt);

                        currentUser.setSalt(newSalt); 
                        currentUser.setPassword(newPwd);

                        userRepository.save(currentUser);

                        userResponse.setMessage("Password changed");
                        userResponse.setStatus(200);
                        userResponse.setError(false);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Original password not right");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
                        return ResponseEntity.ok(userResponse);
                    }
            
                } 

                else {
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

    @PostMapping("/forgot-password")
    public ResponseEntity<UserResponse> forgotPassword(
        @RequestParam String email, 
        HttpServletRequest request
    ) {
        userResponse = new UserResponse(); 
        User currentUser = userRepository.findByEmail(email); 

        if (currentUser == null || currentUser.getAuthType() == User.AuthType.special) {

            userResponse.setMessage("No account or logged in using social media");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        } else {
            
            userResponse.setMessage("This feature has been disabled. Just reset your password");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);

        }

    }

    @PostMapping(path="/reset-password")
    public ResponseEntity<UserResponse> resetPassword(
        @RequestParam String email,
        @RequestParam String password
    ) {
        User currentUser = userRepository.findByEmail(email);
        String salt = PasswordUtils.getSalt(30);
        String newPassword = PasswordUtils.generateSecurePassword(password, salt);

        currentUser.setSalt(salt);
        currentUser.setPassword(newPassword);

        userRepository.save(currentUser); 

        userResponse.setMessage("Password successfully reset");
        userResponse.setStatus(200);
        userResponse.setError(false);
        userResponse.setTime(new Timestamp(System.currentTimeMillis()));
        return ResponseEntity.ok(userResponse);
    }

    @PostMapping(path="/all-friends")
    public ResponseEntity<UserResponse> allFriends(
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
                        List<User> users = userRepository.findFriendsWithCustomQuery(currentUser);
                        List<UsersAll> customUsers = new ArrayList<>(); 

                        for (User user : users) {
                            UsersAll customUser = new UsersAll(); 
                            customUser.setEmail(user.getEmail());
                            customUser.setFirstName(user.getFirstName());
                            customUser.setLastName(user.getLastName());

                            customUsers.add(customUser); 
                        }

                        userResponse.setUsersAll(customUsers);
                        userResponse.setMessage("Friends attached");
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
                        List<User> users = userRepository.findFriendsWithCustomQuery(currentUser);
                        List<UsersAll> customUsers = new ArrayList<>(); 

                        for (User user : users) {
                            UsersAll customUser = new UsersAll(); 
                            customUser.setEmail(user.getEmail());
                            customUser.setFirstName(user.getFirstName());
                            customUser.setLastName(user.getLastName());

                            customUsers.add(customUser); 
                        }

                        userResponse.setUsersAll(customUsers);
                        userResponse.setMessage("Friends attached");
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
}