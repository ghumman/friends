class UsersController < ApplicationController
    def addUser
    puts "inside add-user endpoint"
    end

    def login
        render json: {error: 'User not Found'}
    end

    def changePassword
    puts "inside change-password endpoint"
    end
    
    def forgotPassword
    puts "inside forgot-password endpoint"
    end

    def resetPassword
    puts "inside reset-password endpoint"
    end
    
    def allFriends
    puts "inside all-friends endpoint"
    end
end
