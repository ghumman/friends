require 'test_helper'

class MessagesControllerTest < ActionDispatch::IntegrationTest
  test "should get sendMessage" do
    get messages_sendMessage_url
    assert_response :success
  end

end
