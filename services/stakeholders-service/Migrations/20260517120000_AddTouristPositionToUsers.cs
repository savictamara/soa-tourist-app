using Microsoft.EntityFrameworkCore.Infrastructure;
using Microsoft.EntityFrameworkCore.Migrations;
using StakeholdersService.Data;

#nullable disable

namespace StakeholdersService.Migrations
{
    /// <inheritdoc />
    [DbContext(typeof(StakeholdersDbContext))]
    [Migration("20260517120000_AddTouristPositionToUsers")]
    public partial class AddTouristPositionToUsers : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<double>(
                name: "Latitude",
                table: "users",
                type: "double precision",
                nullable: true);

            migrationBuilder.AddColumn<double>(
                name: "Longitude",
                table: "users",
                type: "double precision",
                nullable: true);
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "Latitude",
                table: "users");

            migrationBuilder.DropColumn(
                name: "Longitude",
                table: "users");
        }
    }
}
