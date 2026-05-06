using Microsoft.AspNetCore.Mvc;
using Friends.Api.Models;
using Friends.Api.Services;

namespace Friends.Api.Controllers;

[ApiController]
[Route("/")]
public class AuthController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly ILogger<AuthController> _logger;

    public AuthController(IUserService userService, ILogger<AuthController> logger)
    {
        _userService = userService;
        _logger = logger;
    }

    [HttpGet("test")]
    public IActionResult Test()
    {
        return Ok(new { message = "API is working!" });
    }

    [HttpPost("add-user")]
    public async Task<IActionResult> AddUser([FromBody] AddUserRequest request)
    {
        _logger.LogInformation("AddUser called with email: {Email}", request.Email);
        
        try
        {
            // Validate required fields
            if (string.IsNullOrEmpty(request.FirstName) || string.IsNullOrEmpty(request.LastName) ||
                string.IsNullOrEmpty(request.Email) || string.IsNullOrEmpty(request.Password))
            {
                _logger.LogWarning("Missing required fields");
                return BadRequest(new ApiResponse<object> 
                { 
                    Success = false, 
                    Message = "Missing required fields" 
                });
            }

            var existingUser = await _userService.GetUserByEmail(request.Email);
            if (existingUser != null)
            {
                _logger.LogWarning("User already exists: {Email}", request.Email);
                return BadRequest(new ApiResponse<object> 
                { 
                    Success = false, 
                    Message = "User already exists" 
                });
            }

            var success = await _userService.CreateUser(request);
            
            if (!success)
            {
                return BadRequest(new ApiResponse<object> 
                { 
                    Success = false, 
                    Message = "Failed to create user" 
                });
            }

            _logger.LogInformation("User created successfully: {Email}", request.Email);
            
            return Ok(new ApiResponse<object> 
            { 
                Success = true, 
                Message = "User created successfully. Please check your email." 
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating user");
            return StatusCode(500, new ApiResponse<object> 
            { 
                Success = false, 
                Message = "An error occurred" 
            });
        }
    }

    [HttpPost("login")]
    public async Task<IActionResult> Login([FromBody] LoginRequest request)
    {
        _logger.LogInformation("Login called for email: {Email}", request.Email);
        
        if (string.IsNullOrEmpty(request.Email) || string.IsNullOrEmpty(request.Password))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Email and password are required" 
            });
        }

        var isValid = await _userService.ValidateUser(request.Email, request.Password);
        
        if (!isValid)
        {
            _logger.LogWarning("Invalid login attempt for: {Email}", request.Email);
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid email or password" 
            });
        }

        var token = _userService.GenerateAuthToken(request.Email);
        
        _logger.LogInformation("User logged in: {Email}", request.Email);
        
        return Ok(new ApiResponse<LoginResponse> 
        { 
            Success = true, 
            Message = "Login successful",
            Data = new LoginResponse 
            { 
                Email = request.Email, 
                Token = token 
            }
        });
    }

    [HttpPost("change-password")]
    public async Task<IActionResult> ChangePassword([FromBody] ChangePasswordRequest request)
    {
        if (string.IsNullOrEmpty(request.Email) || string.IsNullOrEmpty(request.Password) ||
            string.IsNullOrEmpty(request.NewPassword))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Missing required fields" 
            });
        }

        var success = await _userService.ChangePassword(request.Email, request.Password, request.NewPassword);
        
        if (!success)
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Failed to change password. Please check your current password." 
            });
        }

        return Ok(new ApiResponse<object> 
        { 
            Success = true, 
            Message = "Password changed successfully" 
        });
    }

    [HttpPost("forgot-password")]
    public async Task<IActionResult> ForgotPassword([FromBody] ForgotPasswordRequest request)
    {
        if (string.IsNullOrEmpty(request.Email))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Email is required" 
            });
        }

        var success = await _userService.CreatePasswordResetToken(request.Email, out var token);
        
        if (success)
        {
            // Simulate sending email
            Console.WriteLine($"Password reset token for {request.Email}: {token}");
        }
        
        // Always return success for security (don't reveal if email exists)
        return Ok(new ApiResponse<object> 
        { 
            Success = true, 
            Message = "If the email exists, a password reset link has been sent." 
        });
    }

    [HttpPost("reset-password")]
    public async Task<IActionResult> ResetPassword([FromBody] ResetPasswordRequest request)
    {
        if (string.IsNullOrEmpty(request.Token) || string.IsNullOrEmpty(request.Password))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Token and new password are required" 
            });
        }

        var success = await _userService.ResetPassword(request.Token, request.Password);
        
        if (!success)
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid or expired token" 
            });
        }

        return Ok(new ApiResponse<object> 
        { 
            Success = true, 
            Message = "Password has been reset successfully" 
        });
    }

    [HttpPost("all-friends")]
    public async Task<IActionResult> GetAllFriends([FromBody] AllFriendsRequest request)
    {
        // Validate authentication
        if (!_userService.IsAuthenticated(request.Email, request.Token))
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid or missing authentication token" 
            });
        }

        var isValid = await _userService.ValidateUser(request.Email, request.Password);
        
        if (!isValid)
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid credentials" 
            });
        }

        var friends = await _userService.GetAllFriends(request.Email);
        
        var friendResponses = friends.Select(f => new FriendResponse
        {
            Email = f.Email,
            FirstName = f.FirstName,
            LastName = f.LastName
        }).ToList();

        return Ok(new ApiResponse<List<FriendResponse>> 
        { 
            Success = true, 
            Message = "Friends retrieved successfully",
            Data = friendResponses
        });
    }

    [HttpGet("all-users")]
    public async Task<IActionResult> GetAllUsers()
    {
        var users = await _userService.GetAllUsers();
        
        var userResponses = users.Select(u => new FriendResponse
        {
            Email = u.Email,
            FirstName = u.FirstName,
            LastName = u.LastName
        }).ToList();

        return Ok(new ApiResponse<List<FriendResponse>> 
        { 
            Success = true, 
            Message = "Users retrieved successfully",
            Data = userResponses
        });
    }

}