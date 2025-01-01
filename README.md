# GoHighLevel Integration Example

This is an example Go application demonstrating the use of OAuth2 for authentication and authorization with the GoHighLevel API.

> **Disclaimer:** This application is intended for educational purposes only and is not suitable for production use. It has not been fully tested and may contain bugs or security vulnerabilities.

## Prerequisites

 > **Important:** Before anything else, you need to follow the steps [here](https://highlevel.stoplight.io/docs/integrations/a04191c0fabf9-authorization) to setup a Marketplace App and a Test Account. Don't play with your live data unless you know what you're doing!

- Go 1.16 or later
- A [Marketplace](https://marketplace.gohighlevel.com/) dev account 
- A GoHighLevel account (you can create a test account through Marketplace)
- A `.env` file with the following environment variables:
    - `OAUTH2_CLIENT_ID`
    - `OAUTH2_CLIENT_SECRET`
    - `OAUTH2_SCOPES`
    - `OAUTH2_BASE_URL`

## Installation

1. Clone the repository:

```
git clone https://github.com/dgomespt/oauth2-example-app.git
cd oauth2-example-app
```

2. Create a `.env` file in the root directory and add your OAuth2 credentials:

```
OAUTH2_CLIENT_ID=your_client_id
OAUTH2_CLIENT_SECRET=your_client_secret
OAUTH2_SCOPES=scope1,scope2
OAUTH2_BASE_URL=https://services.leadconnectorhq.com
```

3. Install dependencies:

```
go mod tidy
```

## Usage

1. Run the application:

```sh
go run main.go
```

2. Open your browser and navigate to `http://localhost:8080`.

3. Click the "Login" link to authenticate with GoHighLevel.

4. After successful authentication, you can use the following endpoints:
        - `/me`: Fetch user information
        - `/contacts?locationId=your_location_id`: Fetch contacts for a specific location

## Endpoints

- `GET /`: Root endpoint with a login link.
- `GET /login`: Redirects to the OAuth2 authorization page.
- `GET /callback`: Handles the OAuth2 callback and exchanges the authorization code for an access token.
- `POST /webhook`: Handles incoming webhooks.
- `GET /me`: Fetches user information from GoHighLevel.
- `GET /contacts`: Fetches contacts from GoHighLevel.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Acknowledgements

- [Go OAuth2](https://github.com/golang/oauth2)
- [GoHighLevel API](https://developers.gohighlevel.com/)
