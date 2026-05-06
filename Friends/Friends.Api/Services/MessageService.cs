using System.Collections.Concurrent;
using Friends.Api.Models;

namespace Friends.Api.Services;

public class MessageService : IMessageService
{
    private static readonly ConcurrentBag<Message> _messages = new();

    public Task<bool> SendMessage(SendMessageRequest request)
    {
        var message = new Message
        {
            Content = request.Message,
            FromEmail = request.MessageFromEmail.ToLower(),
            ToEmail = request.MessageToEmail.ToLower()
        };

        _messages.Add(message);
        return Task.FromResult(true);
    }

    public Task<List<Message>> GetConversation(string userEmail, string friendEmail)
    {
        var conversation = _messages
            .Where(m => (m.FromEmail == userEmail.ToLower() && m.ToEmail == friendEmail.ToLower()) ||
                       (m.FromEmail == friendEmail.ToLower() && m.ToEmail == userEmail.ToLower()))
            .OrderBy(m => m.SentAt)
            .ToList();

        return Task.FromResult(conversation);
    }
}