using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using StakeholdersService.Data;
using StakeholdersService.Entities;
using StakeholdersService.Repositories;
using StakeholdersService.Services;

var builder = WebApplication.CreateBuilder(args);
var allowAngularClientPolicy = "AllowAngularClient";

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();
builder.Services.AddDbContext<StakeholdersDbContext>(options =>
    options.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));
builder.Services.AddCors(options =>
{
    options.AddPolicy(allowAngularClientPolicy, policy =>
    {
        policy.WithOrigins("http://localhost:4200")
            .AllowAnyHeader()
            .AllowAnyMethod();
    });
});

builder.Services.AddScoped<IHealthRepository, HealthRepository>();
builder.Services.AddScoped<IHealthService, HealthService>();
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IAuthService, AuthService>();
builder.Services.AddScoped<IUserService, UserService>();

var app = builder.Build();

using (var scope = app.Services.CreateScope())
{
    var dbContext = scope.ServiceProvider.GetRequiredService<StakeholdersDbContext>();
    dbContext.Database.Migrate();

    if (!dbContext.Users.Any(user => user.Username == "admin"))
    {
        var adminUser = new User
        {
            Username = "admin",
            Email = "admin@stakeholders.local",
            Role = "Administrator",
            IsBlocked = false,
            FirstName = "System",
            LastName = "Admin"
        };

        var passwordHasher = new PasswordHasher<User>();
        adminUser.PasswordHash = passwordHasher.HashPassword(adminUser, "admin123");

        dbContext.Users.Add(adminUser);
        dbContext.SaveChanges();
    }
}

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseCors(allowAngularClientPolicy);
app.UseAuthorization();
app.MapControllers();

app.Run();
