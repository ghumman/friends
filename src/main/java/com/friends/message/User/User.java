package com.friends.message.User;

import java.util.Date;
import javax.validation.constraints.Email;
import javax.validation.constraints.NotNull;

import org.springframework.data.annotation.Id;
import org.springframework.data.mongodb.core.mapping.Document;

import lombok.Data;

@Data
@Document
public class User {

	public enum AuthType {regular, special};

	@Id
	private String userId;
    
    @NotNull
    private String firstName; 
    @NotNull
    private String lastName; 

    @NotNull
    private Date createdAt;

    private String salt;
    private String password; 
    private String token; 

    private String resetToken; 
    
    // @Column(unique = true) 
    @NotNull
    @Email
    private String email;
    
    @NotNull
	private Integer authType; 

	    User(String fName, String lName, String newEmail , String pwd, String auth, String tkn) {

        if (auth.equals(AuthType.regular.toString())) {
            salt = PasswordUtils.getSalt(30);
			password = PasswordUtils.generateSecurePassword(pwd, salt);
			

            authType = 0; 
        } else if (auth.equals(AuthType.special.toString())) {
            token = tkn; 
            authType = 1;
        }

        firstName = fName; 
        lastName = lName; 
        email = newEmail; 
        Date date = new Date();
        createdAt = new Date(date.getTime()); 

    }
	
	User() {}
}
