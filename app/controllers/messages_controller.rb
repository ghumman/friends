require 'mongo'


$db = Mongo::Client.new([ '127.0.0.1:27017' ], :database => 'friends_mongo')
$userCollection = $db[:user]
$messageCollection = $db[:message]

class MessagesController < ApplicationController
    skip_before_action :verify_authenticity_token

    def sendMessage
        message = params[:message]
        messageFromEmail = params[:messageFromEmail]
        messageToEmail = params[:messageToEmail]
        password = params[:password]

        if message == nil || messageFromEmail == nil || messageToEmail == nil || password == nil
            return render json: {
                status: 400,
                error: true,
                message: "Required data missing",
                time: DateTime.now() 
            }
        end

        dataSender = $userCollection.find( { email: messageFromEmail } ).first

        if dataSender == nil
            return render json: {
                status: 400,
                error: true,
                message: "Sender Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, dataSender["salt"])

        if dataSender["password"] == Base64.encode64(key).gsub("\n",'')

            dataReceiver = $userCollection.find( { email: messageToEmail } ).first

            if dataReceiver == nil
                return render json: {
                    status: 400,
                    error: true,
                    message: "Receiver Does Not Exist",
                    time: DateTime.now() 
                }
            end

            # Sender credentials are correct and both sender and receiver exists
            saveMessage(message, dataSender , dataReceiver )

            return render json: {
                status: 200,
                error: false,
                message: "Message sent",
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
    
    def messagesUserAndFriend
        userEmail = params[:userEmail]
        friendEmail = params[:friendEmail]
        password = params[:password]

        if userEmail == nil || friendEmail == nil || password == nil
            return render json: {
                status: 400,
                error: true,
                message: "Required data missing",
                time: DateTime.now() 
            }
        end

        dataSender = $userCollection.find( { email: userEmail } ).first

        if dataSender == nil
            return render json: {
                status: 400,
                error: true,
                message: "Sender Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, dataSender["salt"])

        if dataSender["password"] == Base64.encode64(key).gsub("\n",'')

            dataReceiver = $userCollection.find( { email: friendEmail } ).first

            if dataReceiver == nil
                return render json: {
                    status: 400,
                    error: true,
                    message: "Receiver Does Not Exist",
                    time: DateTime.now() 
                }
            end

            result = $messageCollection.find( { '$or' => [{'messageFrom' => dataSender, 'messageTo' => dataReceiver},  {'messageFrom' => dataReceiver, 'messageTo' => dataSender} ] })

            messages = []

            result.each { |item|
                messages.append({"message" => item["message"], "messageFromEmail" => item['messageFrom']['email'], "messageToEmail" => item['messageTo']['email'], "sentAt" => item["sent_at"]})
            }

            return render json: {
                status: 200,
                error: false,
                message: "Messages attached",
                time: DateTime.now(),
                msgs: messages
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

    def saveMessage(message, sender, receiver)

        #Create new message
        newMessage = { message: message, sentAt: DateTime.now().to_s, messageFrom: sender, messageTo: receiver}
        $messageCollection.insert_one(newMessage)

    end
 
end
