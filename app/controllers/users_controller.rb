require 'date'
require 'openssl'
require "base64"
require 'securerandom'
require 'net/smtp'
require 'mongo'

$db = Mongo::Client.new([ '127.0.0.1:27017' ], :database => 'friends_mongo')
$userCollection = $db[:user]
$messageCollection = $db[:message]

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

        data = $userCollection.find( { email: email } ).first

        if data != nil 
            return render json: {
                status: 400,
                error: true,
                message: "User Already Exists",
                time: DateTime.now() 
            }
        end

        salt = generateSalt()
        dbPassword = Base64.encode64(generateKey(password, salt)).gsub("\n",'')

        #Create new user
        newUser = { createdAt: DateTime.now().to_s, email: email, firstName: firstName, lastName: lastName, password: dbPassword, salt: salt}
        result = $userCollection.insert_one(newUser)

        # UserMailer.with(email: email, accountCreated: true, token: "").welcome_email.deliver_now
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

        result = $userCollection.find( { email: email } ).first

        if result == nil
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, result["salt"])

        if result["password"] == Base64.encode64(key).gsub("\n",'')
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

        data = $userCollection.find( { email: email } ).first

        if data == nil 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, data["salt"])

        if data["password"] == Base64.encode64(key).gsub("\n",'')
            salt = generateSalt()
            dbPassword = Base64.encode64(generateKey(newPassword, salt)).gsub("\n",'')

            data = $userCollection.update_one( { 'email' => email }, { '$set' => { 'salt' => salt,  'password' => dbPassword} } )
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

        result = $userCollection.find( { email: email } ).first

        if result == nil
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end

        uuidNumber = SecureRandom.uuid 
        data = $userCollection.update_one( { 'email' => email }, { '$set' => { 'resetToken' => uuidNumber.to_s} } )

        # UserMailer.with(email: email, accountCreated: false, token: uuidNumber.to_s).welcome_email.deliver_now

        return render json: {
            status: 200,
            error: false,
            message: "Reset password is sent",
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

        data = $userCollection.find( { resetToken: token } ).first

        if data == nil 
            return render json: {
                status: 400,
                error: true,
                message: "Token is not valid",
                time: DateTime.now() 
            }
        end

        salt = generateSalt()
        dbPassword = Base64.encode64(generateKey(password, salt)).gsub("\n",'')

        $userCollection.update_one( { 'email' => data['email'] }, { '$set' => { 'salt' => salt,  'password' => dbPassword, 'resetToken' => nil} } )

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
        
        result = $userCollection.find( { email: email } ).first


        if result == nil 
            return render json: {
                status: 400,
                error: true,
                message: "User Does Not Exist",
                time: DateTime.now() 
            }
        end
        key = generateKey(password, result["salt"])

        if result["password"] == Base64.encode64(key).gsub("\n",'')

            friends = $userCollection.find( { 'email' => {'$ne' => email}  } )

            users = []

            friends.each { |item|
                users.append({"firstName" => item["firstName"], "lastName" => item["lastName"], "email" => item["email"]})
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
