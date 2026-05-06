using Friends.Api.Models;

namespace Friends.Api.Services;

public interface IUserService
{
    Task<User?> GetUserByEmail(string email);
    Task<List<User>> GetAllUsers();
    Task<bool> CreateUser(AddUserRequest request);
    Task<bool> ValidateUser(string email, string password);
    Task<bool> ChangePassword(string email, string oldPassword, string newPassword);
    Task<bool> CreatePasswordResetToken(string email, out string token);
    Task<bool> ResetPassword(string token, string newPassword);
    Task<List<User>> GetAllFriends(string email);
    Task<bool> AddFriend(string userEmail, string friendEmail);
    bool IsAuthenticated(string email, string token);
    string GenerateAuthToken(string email);
}