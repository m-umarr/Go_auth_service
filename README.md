# Backend Engineer Skills Test

This project demonstrates a basic setup of two microservices (`auth_service` and `otp_service`) using Golang, gRPC (Connect), PostgreSQL, and Twilio.

## Services

- **auth_service**: Manages user authentication and profile.
- **otp_service**: Generates OTP using Twilio's API.

## Setup Instructions

1. **Clone the repository**:
    ```sh
    git clone <repository_url>
    cd <repository_name>
    ```

2. **Environment Variables**:
    - Create a `.env` file in the `auth_service` and `otp_service` directories and add the following:
    ```env

    # /.env
    TWILIO_ACCOUNT_SID=your_twilio_account_sid
    TWILIO_AUTH_TOKEN=your_twilio_auth_token
    ```

3. **Run Docker Compose**:
    ```sh
    docker-compose up --build
    ```

4. **Endpoints**:

    - `auth-service`:
        - `/auth.SignupWithPhoneNumber`
        - `/auth.VerifyPhoneNumber`
        - `/auth.LoginWithPhoneNumber`
        - `/auth.ValidatePhoneNumberLogin`
        - `/auth.GetProfile`

    - `otp-service`:
        - `/otp.GenerateOtp`


