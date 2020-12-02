class UsersController < ApplicationController
    def addUser
    puts "inside add-user endpoint"
    end

    def login
        $name = "ghummantech@gmail.com"
        $sql = "select * from user where email=\"#$name\""
        # render json: {values: $sql}
        # sql = "select * from user where email=#$name"
        render json: {values: ActiveRecord::Base.connection.execute($sql)}
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
