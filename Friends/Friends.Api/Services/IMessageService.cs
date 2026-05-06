using Friends.Api.Models;

namespace Friends.Api.Services;

public interface IMessageService
{
    Task<bool> SendMessage(SendMessageRequest request);
    Task<List<Message>> GetConversation(string userEmail, string friendEmail);
}