using Microsoft.AspNetCore.Mvc;
using Friends.Api.Models;
using Friends.Api.Services;

namespace Friends.Api.Controllers;

[ApiController]
[Route("/")]
public class MessageController : ControllerBase
{
    private readonly IUserService _userService;
    private readonly IMessageService _messageService;
    private readonly ILogger<MessageController> _logger;

    public MessageController(IUserService userService, IMessageService messageService, ILogger<MessageController> logger)
    {
        _userService = userService;
        _messageService = messageService;
        _logger = logger;
    }

    [HttpPost("send-message")]
    public async Task<IActionResult> SendMessage([FromBody] SendMessageRequest request)
    {
        _logger.LogInformation("SendMessage from {FromEmail} to {ToEmail}", request.MessageFromEmail, request.MessageToEmail);
        
        // Validate required fields
        if (string.IsNullOrEmpty(request.Message) || string.IsNullOrEmpty(request.MessageFromEmail) ||
            string.IsNullOrEmpty(request.MessageToEmail) || string.IsNullOrEmpty(request.Password))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Missing required fields" 
            });
        }

        // Validate authentication
        if (!_userService.IsAuthenticated(request.MessageFromEmail, request.Token))
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid or missing authentication token" 
            });
        }

        // Validate user credentials
        var isValid = await _userService.ValidateUser(request.MessageFromEmail, request.Password);
        
        if (!isValid)
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid credentials" 
            });
        }

        // Check if sender and receiver exist
        var sender = await _userService.GetUserByEmail(request.MessageFromEmail);
        var receiver = await _userService.GetUserByEmail(request.MessageToEmail);
        
        if (sender == null || receiver == null)
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Sender or receiver not found" 
            });
        }

        // Add as friends
        await _userService.AddFriend(request.MessageFromEmail, request.MessageToEmail);
        await _userService.AddFriend(request.MessageToEmail, request.MessageFromEmail);

        // Send message
        var success = await _messageService.SendMessage(request);
        
        if (!success)
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Failed to send message" 
            });
        }

        _logger.LogInformation("Message sent successfully");
        
        return Ok(new ApiResponse<object> 
        { 
            Success = true, 
            Message = "Message sent successfully" 
        });
    }

    [HttpPost("messages-user-and-friend")]
    public async Task<IActionResult> GetConversation([FromBody] MessagesUserAndFriendRequest request)
    {
        _logger.LogInformation("Get conversation between {UserEmail} and {FriendEmail}", request.UserEmail, request.FriendEmail);
        
        // Validate required fields
        if (string.IsNullOrEmpty(request.UserEmail) || string.IsNullOrEmpty(request.FriendEmail) ||
            string.IsNullOrEmpty(request.Password))
        {
            return BadRequest(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Missing required fields" 
            });
        }

        // Validate authentication
        if (!_userService.IsAuthenticated(request.UserEmail, request.Token))
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid or missing authentication token" 
            });
        }

        // Validate user credentials
        var isValid = await _userService.ValidateUser(request.UserEmail, request.Password);
        
        if (!isValid)
        {
            return Unauthorized(new ApiResponse<object> 
            { 
                Success = false, 
                Message = "Invalid credentials" 
            });
        }

        // Get conversation
        var conversation = await _messageService.GetConversation(request.UserEmail, request.FriendEmail);
        
        var messageResponses = conversation.Select(m => new MessageResponse
        {
            Id = m.Id,
            Content = m.Content,
            FromEmail = m.FromEmail,
            ToEmail = m.ToEmail,
            SentAt = m.SentAt,
            IsRead = m.IsRead
        }).ToList();

        return Ok(new ApiResponse<List<MessageResponse>> 
        { 
            Success = true, 
            Message = "Conversation retrieved successfully",
            Data = messageResponses
        });
    }
}