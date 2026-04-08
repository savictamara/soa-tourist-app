FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY services/blog-service/BlogService.csproj ./
RUN dotnet restore BlogService.csproj

COPY services/blog-service/. ./
RUN dotnet publish BlogService.csproj -c Release -o /app/publish /p:UseAppHost=false

FROM build AS migration
WORKDIR /src
RUN dotnet tool install --global dotnet-ef --version 8.*
ENV PATH="${PATH}:/root/.dotnet/tools"
ENTRYPOINT ["dotnet", "ef", "database", "update", "--project", "BlogService.csproj", "--startup-project", "BlogService.csproj"]

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
COPY --from=build /app/publish ./

ENV ASPNETCORE_URLS=http://+:8080
EXPOSE 8080

ENTRYPOINT ["dotnet", "BlogService.dll"]
