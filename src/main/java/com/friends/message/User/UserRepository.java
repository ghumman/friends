package com.friends.message.User;

import java.util.List;
import java.util.Optional;


import org.springframework.data.mongodb.repository.MongoRepository;
import org.springframework.data.mongodb.repository.Query;

import org.springframework.stereotype.Repository;


@Repository
public interface UserRepository extends MongoRepository<User, String> {
    User findByEmail(String email);
    Optional<User> findByResetToken(String resetToken);


    // List<User> findByEmailNotOrderByFirstNameAndLastName(String email);
    @Query("{email : {$ne : ?0}}")
    public List<User> findByNameNotQuery(String email);
    
}
