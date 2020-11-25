require 'test_helper'

class FriendsControllerTest < ActionDispatch::IntegrationTest
  test "should get login" do
    get friends_login_url
    assert_response :success
  end

end
