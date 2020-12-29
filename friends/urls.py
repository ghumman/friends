from django.urls import path

from . import views

urlpatterns = [
    path('add-user/', views.addUser, name='addUser'),
    path('login/', views.login, name='login'),
    path('change-password/', views.changePassword, name='changePassword'),
    path('forgot-password/', views.forgotPassword, name='forgotPassword'),
    path('reset-password/', views.resetPassword, name='resetPassword'),
    path('all-friends/', views.allFriends, name='allFriends'),
    path('send-message/', views.sendMessage, name='sendMessage'),
    path('messages-user-and-friend/', views.messagesUserAndFriend, name='messagesUserAndFriend'),
]