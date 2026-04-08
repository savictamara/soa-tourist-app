using BlogService.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace BlogService.Data.Configurations;

public class BlogCommentConfiguration : IEntityTypeConfiguration<BlogComment>
{
    public void Configure(EntityTypeBuilder<BlogComment> builder)
    {
        builder.ToTable("blog_comments");

        builder.HasKey(comment => comment.Id);

        builder.Property(comment => comment.Id)
            .ValueGeneratedOnAdd();

        builder.Property(comment => comment.AuthorUsername)
            .HasMaxLength(100)
            .IsRequired();

        builder.Property(comment => comment.AuthorRole)
            .HasMaxLength(30)
            .IsRequired();

        builder.Property(comment => comment.Text)
            .HasMaxLength(2000)
            .IsRequired();

        builder.Property(comment => comment.CreatedAtUtc)
            .IsRequired();

        builder.Property(comment => comment.LastModifiedAtUtc)
            .IsRequired();
    }
}
