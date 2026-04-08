FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src

COPY services/stakeholders-service/StakeholdersService.csproj ./
RUN dotnet restore StakeholdersService.csproj

COPY services/stakeholders-service/. ./
RUN dotnet publish StakeholdersService.csproj -c Release -o /app/publish /p:UseAppHost=false

FROM build AS migration
WORKDIR /src
RUN dotnet tool install --global dotnet-ef --version 8.*
ENV PATH="${PATH}:/root/.dotnet/tools"
ENTRYPOINT ["dotnet", "ef", "database", "update", "--project", "StakeholdersService.csproj", "--startup-project", "StakeholdersService.csproj"]

FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS final
WORKDIR /app
COPY --from=build /app/publish ./

ENV ASPNETCORE_URLS=http://+:8080
EXPOSE 8080

ENTRYPOINT ["dotnet", "StakeholdersService.dll"]
