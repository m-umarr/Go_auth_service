# Backend Engineer Skills Test

This project demonstrates a basic setup of two microservices (`auth-service` and `otp-service`) using Golang, gRPC (Connect), PostgreSQL, and Twilio.

## Services

- **auth-service**: Manages user authentication and profile.
- **otp-service**: Generates OTP using Twilio's API.

## Setup Instructions

1. **Clone the repository**:
    ```sh
    git clone <repository_url>
    cd <repository_name>
    ```

2. **Environment Variables**:
    - Create a `.env` file in the `auth-service` and `otp-service` directories and add the following:
    ```env
    # auth-service/.env
    JWT_SECRET=your_jwt_secret

    # otp-service/.env
    TWILIO_ACCOUNT_SID=your_twilio_account_sid
    TWILIO_AUTH_TOKEN=your_twilio_auth_token
    ```

3. **Run Docker Compose**:
    ```sh
    docker-compose up --build
    ```

4. **Endpoints**:

    - `auth-service`:
        - `POST /auth.SignupWithPhoneNumber`
        - `POST /auth.VerifyPhoneNumber`
        - `POST /auth.LoginWithPhoneNumber`
        - `POST /auth.ValidatePhoneNumberLogin`
        - `GET /auth.GetProfile`

    - `otp-service`:
        - `POST /otp.GenerateOtp`

## Video Walkthrough

[Provide a link to the video walkthrough explaining your thought process and implementation details.]

