require 'date'
require 'openssl'
require "base64"
require 'securerandom'
require 'net/smtp'

class UsersController < ApplicationController
    skip_before_action :verify_authenticity_token

    
    def addUser
        email = params[:email]
        password = params[:password]
        firstName = params[:firstName]
        lastName = params[:lastName]

        if email == nil || password == nil || firstName == nil || lastName == nil
            return render json: {
                status: 400,
                error: true,
                message: "Required data missing",
                time: DateTime.now() 
            }
        end

        stmt = "select password, salt from user where email=\"" + email + "\""
        data =  ActiveRecord::Base.connection.exec_query(stmt)

        if data.length() != 0 
            return render json: {
                status: 400,
                error: true,
                message: "User Already Exists",
                time: DateTime.now() 
            }
        end

        salt = generateSalt()
        dbPassword = Base64.encode64(generateKey(password, salt)).gsub("\n",'')

        # Get the last id to create id for user as our table is not auto increment on id column
        stmt = ("select id from user order by id desc limit 1")
        data =  ActiveRecord::Base.connection.exec_query(stmt)
        
        userID = 0
        if data.length() != 0
            userID = data.entries[0].fetch("id") + 1
        else 
            userID = 1
        end

        #Create new user
        stmt = "insert into user (id, auth_type, created_at, email, first_name, last_name, password, salt) VALUES (" + userID.to_s + ", 0, \"" + DateTime.now().to_s + "\", \"" + email + "\", \"" + firstName + "\", \"" + lastName + "\", \"" + dbPassword + "\", \"" + salt + "\")"
        data =  ActiveRecord::Base.connection.exec_query(stmt)

        UserMailer.with(email: email, accountCreated: true, token: "").welcome_email.deliver_now
        return render json: {
                status: 200,
                error: false,
                message: "User Created",
                time: DateTime.now() 
        }
    end

    def login
        email = params[:email]
        password = params[:password]

        if email == nil || password == nil
            return render json: {
                status: 400,
                error: true,
                message: "Email or Password missing",
                time: DateTime.now() 
            }
        end
        
        stmt = "select password, salt from user where email=\"" + email + "\""
        result =  ActiveRecord::Base.connection.exec_query(stmt)

        if result.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, result.entries[0].fetch("salt"))

        if result.entries[0].fetch("password") == Base64.encode64(key).gsub("\n",'')
            return render json: {
                status: 200,
                error: false,
                message: "Logged In",
                time: DateTime.now() 
             }
        else 
            return render json: {
                status: 400,
                error: true,
                message: "Login Failed",
                time: DateTime.now() 
             }
        end
    end

    def changePassword
        email = params[:email]
        password = params[:password]
        newPassword = params[:newPassword]

        if email == nil || password == nil || newPassword == nil
            return render json: {
                status: 400,
                error: true,
                message: "Required data missing",
                time: DateTime.now() 
            }
        end

        stmt = "select password, salt from user where email=\"" + email + "\""
        data =  ActiveRecord::Base.connection.exec_query(stmt)

        if data.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, data.entries[0].fetch("salt"))

        if data.entries[0].fetch("password") == Base64.encode64(key).gsub("\n",'')
            salt = generateSalt()
            dbPassword = Base64.encode64(generateKey(newPassword, salt)).gsub("\n",'')
            stmt = "UPDATE user SET salt=\"" + salt + "\", password=\"" + dbPassword + "\"  where email=\"" + email + "\""
            data =  ActiveRecord::Base.connection.exec_query(stmt)
            return render json: {
                status: 200,
                error: false,
                message: "Password changed",
                time: DateTime.now() 
            }   
        else 
            return render json: {
                status: 400,
                error: true,
                message: "Original password not right",
                time: DateTime.now() 
             }
        end
    end
    
    def forgotPassword
        email = params[:email]

        if email == nil 
            return render json: {
                status: 400,
                error: true,
                message: "Email is required",
                time: DateTime.now() 
            }
        end

        stmt = "select password, salt from user where email=\"" + email + "\""
        result =  ActiveRecord::Base.connection.exec_query(stmt)

        if result.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        uuidNumber = SecureRandom.uuid 
        stmt = "UPDATE user SET reset_token=\"" + uuidNumber.to_s + "\" where email=\"" + email + "\""
        data =  ActiveRecord::Base.connection.exec_query(stmt)
        
        UserMailer.with(email: email, accountCreated: false, token: uuidNumber.to_s).welcome_email.deliver_now

        return render json: {
            status: 200,
            error: false,
            message: "Password changed",
            time: DateTime.now() 
        }  

    end

    def resetPassword
        token = params[:token]
        password = params[:password]

        if token == nil || password == nil
            return render json: {
                status: 400,
                error: true,
                message: "Required data missing",
                time: DateTime.now() 
            }
        end

        stmt = "select id from user where reset_token=\"" + token + "\""
        data =  ActiveRecord::Base.connection.exec_query(stmt)

        if data.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "Token is not valid",
                time: DateTime.now() 
            }
        end

        salt = generateSalt()
        dbPassword = Base64.encode64(generateKey(password, salt)).gsub("\n",'')
        stmt = "UPDATE user SET salt=\"" + salt + "\", password=\"" + dbPassword + "\", reset_token=null where id=" + data.entries[0].fetch("id").to_s
        data =  ActiveRecord::Base.connection.exec_query(stmt)
        return render json: {
            status: 200,
            error: false,
            message: "Password successfully reset",
            time: DateTime.now() 
        } 

    end
    
    def allFriends
        email = params[:email]
        password = params[:password]

        if email == nil || password == nil
            return render json: {
                status: 400,
                error: true,
                message: "Email or Password missing",
                time: DateTime.now() 
            }
        end

        stmt = "select password, salt from user where email=\"" + email + "\""
        result =  ActiveRecord::Base.connection.exec_query(stmt)

        if result.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end
        key = generateKey(password, result.entries[0].fetch("salt"))

        if result.entries[0].fetch("password") == Base64.encode64(key).gsub("\n",'')

            stmt = "select first_name, last_name, email FROM user where email!=\"" + email + "\""
            result =  ActiveRecord::Base.connection.exec_query(stmt)

            users = []

            result.each { |item|
                users.append({"firstName" => item.fetch("first_name"), "lastName" => item.fetch("last_name"), "email" => item.fetch("email")})
            }


            return render json: {
                status: 200,
                error: false,
                message: "Friends attached",
                time: DateTime.now(),
                usersAll: users
             }
        else 
            return render json: {
                status: 400,
                error: true,
                message: "Login Failed",
                time: DateTime.now() 
             }
        end

    end


    def generateSalt ()
        $ALPHABET = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
        returnValue = ""
        for x in 0..29
            returnValue = returnValue + $ALPHABET[rand(0..29)]
        end
        return returnValue
    end

end
