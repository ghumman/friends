Rails.application.routes.draw do
  # For details on the DSL available within this file, see https://guides.rubyonrails.org/routing.html
  get '/add-user', to: 'users#addUser'
  get '/login', to: 'users#login'
  get '/change-password', to: 'users#changePassword'
  get '/forgot-password', to: 'users#forgotPassword'
  get '/reset-password', to: 'users#resetPassword'
  get '/all-friends', to: 'users#allFriends'

  get '/send-message', to: 'messages#sendMessage'
  get '/messages-user-and-friend', to: 'messages#messagesUserAndFriend'
end
