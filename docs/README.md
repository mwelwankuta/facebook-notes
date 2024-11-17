<h1 align="center">Facebook Notes Documentation</h1>
<p align="center">
OpenAI Powered X (formerly Twitter) Community Notes "Equivalent" for Facebook
</p>

Note: _Docs are not completed_

## About
Facebook Notes is an AI-powered community fact-checking system inspired by X's Community Notes. It allows users to collaboratively add context and fact-checking notes to Facebook posts.

## Features
- AI-powered fact verification
- Community-driven note creation and validation
- Integration with Facebook posts
- Authentication system for verified contributors
- Automated summary generation

## Getting Started

### Prerequisites
- Go 1.19 or higher
- Facebook Developer Account
- OpenAI API key

### Installation
1. Clone the repository:
    ```sh
    git clone https://github.com/mwelwankuta/facebook-notes.git
    ```
2. Navigate to the project directory:
    ```sh
    cd facebook-notes
    ```
3. Install dependencies:
    ```sh
    go mod tidy
    ```

### Usage
1. Log in with your Facebook Developer Account.
2. Connect your OpenAI API key.
3. Start adding and reviewing notes on Facebook posts.

## API Documentation
API documentation is available via Swagger UI at `/swagger/index.html` when running the server locally. You can also find the OpenAPI specification at `/swagger/doc.json`.

## Configuration

### Environment Variables
- `FACEBOOK_APP_ID`: Your Facebook App ID.
- `FACEBOOK_APP_SECRET`: Your Facebook App Secret.
- `OPENAI_API_KEY`: Your OpenAI API key.
- `DATABASE_URL`: URL for your database connection.

### Authentication Service
For detailed information on the authentication service, refer to the [Authentication Service Documentation](./AUTHENTICATION_SERVICE.md).

### Summaries Service
For detailed information on the summaries service, refer to the [Summaries Service Documentation](./SUMMARIES_SERVICE.md).

## Contributing
We welcome contributions! Please read our [Contributing Guidelines](./CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License
This project is licensed under the MIT License - see the [LICENSE](./LICENSE) file for details.

## Acknowledgments
- Inspired by X's Community Notes
- Powered by OpenAI

## Changelog
All notable changes to this project will be documented in the [CHANGELOG](./CHANGELOG.md) file.

## Frequently Asked Questions (FAQ)
For common questions and answers, please refer to the [FAQ](./FAQ.md) file.
## Environment Variables

The application requires the following environment variables to be set:

### General
- `APP_ENV`: The application environment (e.g., `development`, `production`).
- `PORT`: The port on which the application will run (default: `8080`).

### Summaries Service
- `SUMMARIES_SERVICE_URL`: The URL for the summaries service.
- `SUMMARIES_SERVICE_API_KEY`: The API key used for accessing the summaries service.

### Database
- `DB_HOST`: The database host.
- `DB_PORT`: The database port.
- `DB_USER`: The database user.
- `DB_PASSWORD`: The database password.
- `DB_NAME`: The database name.

### Facebook API
- `FACEBOOK_APP_ID`: The Facebook application ID.
- `FACEBOOK_APP_SECRET`: The Facebook application secret.
- `FACEBOOK_ACCESS_TOKEN`: The Facebook access token.

### OpenAI API
- `OPENAI_API_KEY`: The API key for accessing OpenAI services.

### Example `.env` file
```env
APP_ENV=development
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
FACEBOOK_APP_ID=your_facebook_app_id
FACEBOOK_APP_SECRET=your_facebook_app_secret
FACEBOOK_ACCESS_TOKEN=your_facebook_access_token
OPENAI_API_KEY=your_openai_api_key
```

## Running Services

### Authentication Service
To run the authentication service, navigate to the `auth` directory and execute:
```sh
cd internal/auth
go run main.go
```

### Summaries Service
To run the summaries service, navigate to the `summaries` directory and execute:
```sh
cd internal/summaries
go run main.go
```

### Main Application
To run the main application, navigate to the project root directory and execute:
```sh
go run main.go
```

## Project Structure
The project structure is organized as follows:
```
facebook-notes/
├── cmd/                    # Main applications of the project
│   └── facebook-notes/     # The Facebook Notes application
├── internal/               # Private application and library code
│   ├── auth/               # Authentication service
│   ├── notes/              # Notes management service
│   ├── summaries/          # Summaries generation service
│   └── ...                 # Other internal packages
├── pkg/                    # Public library code
│   └── ...                 # Shared packages
├── docs/                   # Documentation files
├── .env.example            # Example environment variables file
├── go.mod                  # Go module file
├── go.sum                  # Go dependencies file
└── main.go                 # Main application entry point
```

## Requirements
- Go 1.19 or higher
- Facebook Developer Account
- OpenAI API key
- PostgreSQL or any other SQL database

## Testing
To run tests, use the following command:
```sh
go test ./...
```

## Deployment
### Docker
You can deploy the application using Docker. Follow these steps:
1. Build the Docker image:
    ```sh
    docker build -t facebook-notes .
    ```
2. Run the Docker container:
    ```sh
    docker run -d -p 8080:8080 --env-file .env facebook-notes
    ```

## Support
If you encounter any issues or have questions, please open an issue on the [GitHub repository](https://github.com/mwelwankuta/facebook-notes/issues).