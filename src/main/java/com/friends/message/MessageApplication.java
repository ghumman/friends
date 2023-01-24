package com.friends.message;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.web.servlet.support.SpringBootServletInitializer;

@SpringBootApplication

// @ComponentScan("com")
// @SpringBootApplication
// @EnableWebMvc
// @Profile("lambda")
public class MessageApplication extends SpringBootServletInitializer {	public static void main(String[] args) {
		SpringApplication.run(MessageApplication.class, args);
	}
}


// public class MessageApplication {

// 	public static void main(String[] args) {

// 		SpringApplication.run(MessageApplication.class, args);

// 	}


// }
