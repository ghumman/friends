namespace Friends.Api.Models
{
    public class User
    {
        public string Id { get; set; } = Guid.NewGuid().ToString();
        public string FirstName { get; set; } = string.Empty;
        public string LastName { get; set; } = string.Empty;
        public string Email { get; set; } = string.Empty;
        public string PasswordHash { get; set; } = string.Empty;
        public string AuthType { get; set; } = string.Empty; 
        public string? ResetToken { get; set; }
        public DateTime? ResetTokenExpiry { get; set; }
        public List<string> Friends { get; set; } = new();
        public DateTime CreatedAt { get; set; } = DateTime.UtcNow;

    }
}