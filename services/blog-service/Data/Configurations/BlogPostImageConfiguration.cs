using BlogService.Entities;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Metadata.Builders;

namespace BlogService.Data.Configurations;

public class BlogPostImageConfiguration : IEntityTypeConfiguration<BlogPostImage>
{
    public void Configure(EntityTypeBuilder<BlogPostImage> builder)
    {
        builder.ToTable("blog_post_images");

        builder.HasKey(image => image.Id);

        builder.Property(image => image.Id)
            .ValueGeneratedOnAdd();

        builder.Property(image => image.ImageUrl)
            .HasColumnType("text")
            .IsRequired();
    }
}
