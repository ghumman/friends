package com.friends.message.Users;

import java.util.List;
import java.util.Optional;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;

public interface UsersRepository extends CrudRepository<Users, Integer>{
    Users findByEmail(String email);
    Optional<Users> findByResetToken(String resetToken);

    @Query("SELECT u FROM Users u WHERE u != ?1 order by u.firstName, u.lastName")
    List<Users> findFriendsWithCustomQuery(Users user);
}