class UserMailer < ApplicationMailer
    default from: 'username@example.com'
 
    def welcome_email
      @email = params[:email]
      @accountCreated = params[:accountCreated]
      @token = params[:token]

      title = @accountCreated ? "Welcome to Friends" : "Password Reset Request"
      text = @accountCreated ? "New Account Created." : "To reset your password, click the link below:\n http://localhost:3000/#/reset-password?token=" + @token

      print @user
      mail(to: @email, subject: title,  body: text)
    end
end



