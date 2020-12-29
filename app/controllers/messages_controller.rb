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

        stmt = "select password, salt, id from users where email=\'" + messageFromEmail + "\'"
        dataSender =  ActiveRecord::Base.connection.exec_query(stmt)

        if dataSender.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "Sender Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, dataSender.entries[0].fetch("salt"))

        if dataSender.entries[0].fetch("password") == Base64.encode64(key).gsub("\n",'')
            
            stmt = "select password, salt, id from users where email=\'" + messageToEmail + "\'"
            dataReceiver =  ActiveRecord::Base.connection.exec_query(stmt)

            if dataReceiver.length() == 0 
                return render json: {
                    status: 400,
                    error: true,
                    message: "Receiver Does Not Exist",
                    time: DateTime.now() 
                }
            end

            # Sender credentials are correct and both sender and receiver exists
            saveMessage(message, dataSender.entries[0].fetch("id") , dataReceiver.entries[0].fetch("id") )

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
        stmt = "select password, salt, id from users where email=\'" + userEmail + "\'"
        dataSender =  ActiveRecord::Base.connection.exec_query(stmt)

        if dataSender.length() == 0 
            return render json: {
                status: 400,
                error: true,
                message: "Sender Does Not Exist",
                time: DateTime.now() 
            }
        end

        key = generateKey(password, dataSender.entries[0].fetch("salt"))

        if dataSender.entries[0].fetch("password") == Base64.encode64(key).gsub("\n",'')
            
            stmt = "select password, salt, id from users where email=\'" + friendEmail + "\'"
            dataReceiver =  ActiveRecord::Base.connection.exec_query(stmt)

            if dataReceiver.length() == 0 
                return render json: {
                    status: 400,
                    error: true,
                    message: "Receiver Does Not Exist",
                    time: DateTime.now() 
                }
            end

            stmt = "SELECT message, message_from_id, message_to_id, sent_at  FROM message m WHERE (m.message_from_id = " + dataSender.entries[0].fetch("id").to_s + " and m.message_to_id = " + dataReceiver.entries[0].fetch("id").to_s + ") or (m.message_from_id = " + dataReceiver.entries[0].fetch("id").to_s + " and m.message_to_id = " + dataSender.entries[0].fetch("id").to_s + ") order by m.sent_at"
            result =  ActiveRecord::Base.connection.exec_query(stmt)

            messages = []

            result.each { |item|
                if item.fetch("message_from_id") == dataSender.entries[0].fetch("id") && item.fetch("message_to_id") == dataReceiver.entries[0].fetch("id") 
                    messages.append({"message" => item.fetch("message"), "messageFromEmail" => userEmail, "messageToEmail" => friendEmail, "sentAt" => item.fetch("sent_at")})
                else
                    messages.append({"message" => item.fetch("message"), "messageFromEmail" => friendEmail, "messageToEmail" => userEmail, "sentAt" => item.fetch("sent_at")})
                end
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

    def saveMessage(message, senderID, receiverID)

        stmt = ("SELECT id FROM message ORDER BY id DESC LIMIT 1")
        data =  ActiveRecord::Base.connection.exec_query(stmt)
        
        messageID = 0
        if data.length() != 0
            messageID = data.entries[0].fetch("id") + 1
        else 
            messageID = 1
        end

        #Create new message
        stmt = "insert into message (id, message, sent_at, message_from_id, message_to_id) VALUES (" + messageID.to_s + ", \'" + message + "\', \'" + DateTime.now().to_s + "\', " + senderID.to_s + ", " + receiverID.to_s + ")"
        data =  ActiveRecord::Base.connection.exec_query(stmt)

    end
 
end
