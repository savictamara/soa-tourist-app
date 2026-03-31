using BlogService.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace BlogService.Data.Configurations;

public class BlogPostConfiguration : IEntityTypeConfiguration<BlogPost>
{
    public void Configure(EntityTypeBuilder<BlogPost> builder)
    {
        builder.ToTable("blog_posts");

        builder.HasKey(post => post.Id);

        builder.Property(post => post.Id)
            .ValueGeneratedOnAdd();

        builder.Property(post => post.Title)
            .HasMaxLength(200)
            .IsRequired();

        builder.Property(post => post.DescriptionMarkdown)
            .HasMaxLength(12000)
            .IsRequired();

        builder.Property(post => post.CreatedAtUtc)
            .IsRequired();

        builder.Property(post => post.AuthorUsername)
            .HasMaxLength(100)
            .IsRequired();

        builder.HasMany(post => post.Images)
            .WithOne(image => image.BlogPost)
            .HasForeignKey(image => image.BlogPostId)
            .OnDelete(DeleteBehavior.Cascade);
    }
}
