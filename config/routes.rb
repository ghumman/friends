Rails.application.routes.draw do
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
  post '/add-user', to: 'users#addUser'
  post '/login', to: 'users#login'
  post '/change-password', to: 'users#changePassword'
  post '/forgot-password', to: 'users#forgotPassword'
  post '/reset-password', to: 'users#resetPassword'
  post '/all-friends', to: 'users#allFriends'

  post '/send-message', to: 'messages#sendMessage'
  post '/messages-user-and-friend', to: 'messages#messagesUserAndFriend'
end
