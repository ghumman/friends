Rails.application.routes.draw do

  get '/add-user', to: 'friends#addUser'
  get '/login', to: 'friends#login'
  get '/change-password', to: 'friends#changePassword'
  get '/forgot-password', to: 'friends#forgotPassword'
  get '/reset-password', to: 'friends#resetPassword'
  get '/all-friends', to: 'friends#allFriends'

  get '/send-message', to: 'messages#sendMessage'
  get '/messages-user-and-friend', to: 'messages#messagesUserAndFriend'

  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
end
