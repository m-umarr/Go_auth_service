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
    - Create a `.env` file in the root directory and add the following:
    ```env

    # /.env
    TWILIO_ACCOUNT_SID=YOUR_ACCOUNT_SID
    TWILIO_AUTH_TOKEN=YOUR_AUTH_TOKEN
    TWILIO_FROM_PHONE=YOUR_PHONE_NUMBER
    TWILIO_SERVICES_ID=YOUR_SERVICE_ID
    AMQP_URL=YOUR_AMPQ_URL
    DB_HOST=localhost
    DB_PORT=5432
    DB_USER=USERNAME
    DB_PASSWORD=PASSWORD
    DB_NAME=DB
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


