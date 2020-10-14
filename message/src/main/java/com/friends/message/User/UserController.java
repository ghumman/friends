package com.friends.message.User;

import org.springframework.web.bind.annotation.RestController;

import java.sql.Timestamp;
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
    private UserRepository userRepository;

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
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);
        }

        // input is valid 

        if (userRepository.findByEmail(email) == null) {

            User newUser = new User(firstName, lastName, email, password, authType, token); 
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
        User currentUser = userRepository.findByEmail(email); 

        if (currentUser == null || currentUser.getAuthType() == User.AuthType.special) {
            userResponse.setMessage("No account or logged in using social media");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
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
					+ "/reset?token=" + currentUser.getResetToken());
			
            emailSender.send(passwordResetEmail);
            
            userResponse.setMessage("Reset password is sent");
            userResponse.setStatus(200);
            userResponse.setError(false);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
            return ResponseEntity.ok(userResponse);

        }

    }

    @PostMapping(path="/reset-password")
    public ResponseEntity<UserResponse> resetPassword(
        @RequestParam String token,
        @RequestParam String password
    ) {
        Optional<User> currentUserOpt = userRepository.findByResetToken(token);
        User currentUser = currentUserOpt.get(); 

		if (currentUser == null) {
            userResponse.setMessage("Token is not valid");
            userResponse.setStatus(400);
            userResponse.setError(true);
            userResponse.setTime(new Timestamp(System.currentTimeMillis()));
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