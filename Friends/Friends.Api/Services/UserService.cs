using System.Collections.Concurrent;
using Friends.Api.Helpers;
using Friends.Api.Models;

namespace Friends.Api.Services;

public class UserService : IUserService
{
    private static readonly ConcurrentDictionary<string, User> _users = new();
    private static readonly ConcurrentDictionary<string, string> _authTokens = new();

    public Task<User?> GetUserByEmail(string email)
    {
        _users.TryGetValue(email.ToLower(), out var user);
        return Task.FromResult(user);
    }

    public Task<bool> CreateUser(AddUserRequest request)
    {
        if (_users.ContainsKey(request.Email.ToLower()))
            return Task.FromResult(false);

        var user = new User
        {
            FirstName = request.FirstName,
            LastName = request.LastName,
            Email = request.Email.ToLower(),
            PasswordHash = PasswordHelper.HashPassword(request.Password),
            AuthType = request.AuthType
        };

        return Task.FromResult(_users.TryAdd(request.Email.ToLower(), user));
    }

    public Task<bool> ValidateUser(string email, string password)
    {
        if (!_users.TryGetValue(email.ToLower(), out var user))
            return Task.FromResult(false);

        return Task.FromResult(PasswordHelper.VerifyPassword(password, user.PasswordHash));
    }

    public Task<bool> ChangePassword(string email, string oldPassword, string newPassword)
    {
        if (!_users.TryGetValue(email.ToLower(), out var user))
            return Task.FromResult(false);

        if (!PasswordHelper.VerifyPassword(oldPassword, user.PasswordHash))
            return Task.FromResult(false);

        user.PasswordHash = PasswordHelper.HashPassword(newPassword);
        return Task.FromResult(true);
    }

    public Task<bool> CreatePasswordResetToken(string email, out string token)
    {
        token = string.Empty;
        
        if (!_users.TryGetValue(email.ToLower(), out var user))
            return Task.FromResult(false);

        token = Convert.ToBase64String(Guid.NewGuid().ToByteArray());
        user.ResetToken = token;
        user.ResetTokenExpiry = DateTime.UtcNow.AddHours(1);
        
        return Task.FromResult(true);
    }

    public Task<bool> ResetPassword(string token, string newPassword)
    {
        var user = _users.Values.FirstOrDefault(u => u.ResetToken == token);
        
        if (user == null || user.ResetTokenExpiry < DateTime.UtcNow)
            return Task.FromResult(false);

        user.PasswordHash = PasswordHelper.HashPassword(newPassword);
        user.ResetToken = null;
        user.ResetTokenExpiry = null;
        
        return Task.FromResult(true);
    }

    public Task<List<User>> GetAllFriends(string email)
    {
        var friends = new List<User>();
        
        if (_users.TryGetValue(email.ToLower(), out var user))
        {
            foreach (var friendEmail in user.Friends)
            {
                if (_users.TryGetValue(friendEmail.ToLower(), out var friend))
                    friends.Add(friend);
            }
        }
        
        return Task.FromResult(friends);
    }

    public Task<bool> AddFriend(string userEmail, string friendEmail)
    {
        if (!_users.ContainsKey(userEmail.ToLower()) || !_users.ContainsKey(friendEmail.ToLower()))
            return Task.FromResult(false);

        if (!_users.TryGetValue(userEmail.ToLower(), out var user))
            return Task.FromResult(false);

        if (!user.Friends.Contains(friendEmail.ToLower()))
            user.Friends.Add(friendEmail.ToLower());

        return Task.FromResult(true);
    }

    public bool IsAuthenticated(string email, string token)
    {
        if (string.IsNullOrEmpty(token))
            return false;
            
        _authTokens.TryGetValue(email.ToLower(), out var storedToken);
        return storedToken == token;
    }

    public string GenerateAuthToken(string email)
    {
        var token = Convert.ToBase64String(Guid.NewGuid().ToByteArray());
        _authTokens[email.ToLower()] = token;
        return token;
    }

    public Task<List<User>> GetAllUsers()
    {
        return Task.FromResult(_users.Values.ToList());
    }
}