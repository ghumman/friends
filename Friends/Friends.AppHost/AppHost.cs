var builder = DistributedApplication.CreateBuilder(args);
builder.AddProject<Projects.Friends_Api>("friends-api");

builder.Build().Run();
