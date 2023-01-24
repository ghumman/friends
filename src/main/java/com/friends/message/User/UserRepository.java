package com.friends.message.User;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

public interface UserRepository extends CrudRepository<User, Integer>{
    User findByEmail(String email);
    Optional<User> findByResetToken(String resetToken);

    @Query("SELECT u FROM User u WHERE u != ?1 order by u.firstName, u.lastName")
    List<User> findFriendsWithCustomQuery(User user);
}