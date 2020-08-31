package com.friends.message.User;

import java.util.Optional;

import org.springframework.data.repository.CrudRepository;

public interface UserRepository extends CrudRepository<User, Integer>{
    User findByEmail(String email);
    Optional<User> findByResetToken(String resetToken);
}