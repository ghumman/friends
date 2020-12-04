class ApplicationController < ActionController::Base

    def generateKey (password, salt)
        return OpenSSL::PKCS5.pbkdf2_hmac(password.encode(Encoding::ASCII_8BIT), salt.encode(Encoding::ASCII_8BIT), 10000, 32, "sha1")
    end

end
