using BlogService.Repositories;
using BlogService.Services;

var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

builder.Services.AddScoped<IHealthRepository, HealthRepository>();
builder.Services.AddScoped<IHealthService, HealthService>();

// Prepared for future PostgreSQL wiring through ConnectionStrings:DefaultConnection.
// Register DbContext and persistence services here when data access is introduced.

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();
app.UseAuthorization();
app.MapControllers();

app.Run();
