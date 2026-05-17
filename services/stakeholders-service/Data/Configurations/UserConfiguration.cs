using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;
using StakeholdersService.Entities;

namespace StakeholdersService.Data.Configurations;

public class UserConfiguration : IEntityTypeConfiguration<User>
{
    public void Configure(EntityTypeBuilder<User> builder)
    {
        builder.ToTable("users");

        builder.HasKey(user => user.Id);

        builder.Property(user => user.Id)
            .ValueGeneratedOnAdd();

        builder.Property(user => user.Username)
            .HasMaxLength(100)
            .IsRequired();

        builder.Property(user => user.PasswordHash)
            .HasMaxLength(512)
            .IsRequired();

        builder.Property(user => user.Email)
            .HasMaxLength(255)
            .IsRequired();

        builder.Property(user => user.Role)
            .HasMaxLength(50)
            .IsRequired();

        builder.Property(user => user.FirstName)
            .HasMaxLength(100);

        builder.Property(user => user.LastName)
            .HasMaxLength(100);

        builder.Property(user => user.ProfileImage)
            .HasColumnType("text");

        builder.Property(user => user.Motto)
            .HasMaxLength(250);

        builder.Property(user => user.Biography)
            .HasMaxLength(4000);

        builder.Property(user => user.Latitude);

        builder.Property(user => user.Longitude);

        builder.HasIndex(user => user.Username)
            .IsUnique();

        builder.HasIndex(user => user.Email)
            .IsUnique();
    }
}
