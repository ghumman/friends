# Create the project
```
dotnet new aspire -n Friends
cd Friends
# Create the ASP.NET Core Web API project
dotnet new webapi -n Friends.Api

# Add the API project to your solution
dotnet sln add Friends.Api/Friends.Api.csproj

# Add project reference
dotnet add Friends.AppHost reference Friends.Api
```

# Run the application
```
dotnet run --project Friends.AppHost/Friends.AppHost.csproj
```