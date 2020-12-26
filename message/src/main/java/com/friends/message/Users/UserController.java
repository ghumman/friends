package com.friends.message.Users;

import org.springframework.web.bind.annotation.RestController;

import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

import javax.servlet.http.HttpServletRequest;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.mail.SimpleMailMessage;
import org.springframework.mail.javamail.JavaMailSender;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;


@RestController
@CrossOrigin(origins = "http://localhost:3000")
public class UserController {

    public enum AuthType {regular, special};

    private UserResponse userResponse; 

    @Autowired 
    private UsersRepository userRepository;

    @Autowired
    private JavaMailSender emailSender;

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
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        }

        // input is valid 

        if (userRepository.findByEmail(email) == null) {

            Users newUser = new Users(firstName, lastName, email, password, authType, token); 
            userRepository.save(newUser); 

            // send confirmation email
            SimpleMailMessage message = new SimpleMailMessage(); 
            message.setFrom("server-email@gmail.com");
            message.setTo(email); 
            message.setSubject("Welcome to Friends"); 
            message.setText("New Account Created");
            emailSender.send(message);

            userResponse.setMessage("User Created");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        }

        else {
            userResponse.setMessage("User Already Exists");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
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

            Users currentUser = userRepository.findByEmail(email);

            if (currentUser == null) {

                userResponse.setMessage("User Does Not Exist");
                userResponse.setStatus(400);
                userResponse.setError(true);
                userResponse.setTime(new Date());
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
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Date());
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
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    }
                } else {
                    userResponse.setMessage("Authentication Type not right");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Date());
                    return ResponseEntity.ok(userResponse);
                }

            }
        } catch(Exception e) {
            userResponse = new UserResponse(); 
            userResponse.setMessage(e.toString());
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
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

            Users currentUser = userRepository.findByEmail(email);

            if (currentUser == null) {

                userResponse.setMessage("User Does Not Exist");
                userResponse.setStatus(400);
                userResponse.setError(true);
                userResponse.setTime(new Date());
                return ResponseEntity.ok(userResponse);
            }

            else {

                if (authType.equals(AuthType.special.toString())) {

                    userResponse.setMessage("Logged in using OAuth");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Date());
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
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Original password not right");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    }
            
                } 

                else {
                    userResponse.setMessage("Authentication Type not right");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Date());
                    return ResponseEntity.ok(userResponse);
                }

            }
        } catch(Exception e) {
            userResponse = new UserResponse(); 
            userResponse.setMessage(e.toString());
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        }
    }

    @PostMapping("/forgot-password")
    public ResponseEntity<UserResponse> forgotPassword(
        @RequestParam String email, 
        HttpServletRequest request
    ) {
        userResponse = new UserResponse(); 
        Users currentUser = userRepository.findByEmail(email); 

        if (currentUser == null || currentUser.getAuthType() == Users.AuthType.special) {

            userResponse.setMessage("No account or logged in using social media");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        } else {

            currentUser.setResetToken(UUID.randomUUID().toString());
            userRepository.save(currentUser); 

			String appUrl = request.getScheme() + "://" + request.getServerName();
			
			// Email message
			SimpleMailMessage passwordResetEmail = new SimpleMailMessage();
			passwordResetEmail.setFrom("support@demo.com");
			passwordResetEmail.setTo(currentUser.getEmail());
			passwordResetEmail.setSubject("Password Reset Request");
			passwordResetEmail.setText("To reset your password, click the link below:\n" + appUrl
					+ ":3000/#/reset-password?token=" + currentUser.getResetToken());
			
            emailSender.send(passwordResetEmail);
            
            userResponse.setMessage("Reset password is sent");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);

        }

    }

    @PostMapping(path="/reset-password")
    public ResponseEntity<UserResponse> resetPassword(
        @RequestParam String token,
        @RequestParam String password
    ) {
        Optional<Users> currentUserOpt = userRepository.findByResetToken(token);
        Users currentUser = currentUserOpt.get(); 

		if (currentUser == null) {
            userResponse.setMessage("Token is not valid");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        } else {
            String salt = PasswordUtils.getSalt(30);
            String newPassword = PasswordUtils.generateSecurePassword(password, salt);

            currentUser.setSalt(salt);
            currentUser.setPassword(newPassword);

            userRepository.save(currentUser); 

            userResponse.setMessage("Password successfully reset");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        }
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

            Users currentUser = userRepository.findByEmail(email);

            if (currentUser == null) {

                userResponse.setMessage("User Does Not Exist");
                userResponse.setStatus(400);
                userResponse.setError(true);
                userResponse.setTime(new Date());
                return ResponseEntity.ok(userResponse);
            }

            else {

                if (authType.equals(AuthType.regular.toString())) {

                    String salt = currentUser.getSalt(); 
                    String securedPassword = currentUser.getPassword(); 
                    boolean passwordMatch = PasswordUtils.verifyUserPassword(password, securedPassword, salt);
            
                    if(passwordMatch) 
                    {
                        List<Users> users = userRepository.findFriendsWithCustomQuery(currentUser);
                        List<UsersAll> customUsers = new ArrayList<>(); 

                        for (Users user : users) {
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
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    }
            
                } else if (authType.equals(AuthType.special.toString())) {
                    // login using OAuth
                    boolean passwordMatch = token.equals(currentUser.getToken());
            
                    if(passwordMatch) 
                    {
                        List<Users> users = userRepository.findFriendsWithCustomQuery(currentUser);
                        List<UsersAll> customUsers = new ArrayList<>(); 

                        for (Users user : users) {
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
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    } else {
                        userResponse.setMessage("Login Failed");
                        userResponse.setStatus(400);
                        userResponse.setError(true);
                        userResponse.setTime(new Date());
                        return ResponseEntity.ok(userResponse);
                    }
                } else {
                    userResponse.setMessage("Authentication Type not right");
                    userResponse.setStatus(400);
                    userResponse.setError(true);
                    userResponse.setTime(new Date());
                    return ResponseEntity.ok(userResponse);
                }

            }
        } catch(Exception e) {
            userResponse = new UserResponse(); 
            userResponse.setMessage(e.toString());
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Date());
            return ResponseEntity.ok(userResponse);
        }
    }
}
