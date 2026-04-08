using BlogService.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace BlogService.Data.Configurations;

public class BlogLikeConfiguration : IEntityTypeConfiguration<BlogLike>
{
    public void Configure(EntityTypeBuilder<BlogLike> builder)
    {
        builder.ToTable("blog_likes");

        builder.HasKey(like => like.Id);

        builder.Property(like => like.Id)
            .ValueGeneratedOnAdd();

        builder.Property(like => like.Username)
            .HasMaxLength(100)
            .IsRequired();

        builder.Property(like => like.CreatedAtUtc)
            .IsRequired();

        builder.HasIndex(like => new { like.BlogPostId, like.Username })
            .IsUnique();
    }
}
