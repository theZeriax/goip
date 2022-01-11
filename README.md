# goip

A simple code to get visitor's IP address with Go.
Every time someone visits the URL, the code will send an email to you.
**DO NOT USE IT WITHOUT THE VISITOR PERMISSION, IT IS ONLY FOR TESTING AND EDUCATIONAL PURPOSE!**

## Installation

Install all the packages required:

```bash
go get github.com/joho/godotenv
go get gopkg.in/gomail.v2
```

Then add the following to your `.env` file:

```dotenv
EMAIL_FROM=<SENDER (YOUR EMAIL)>
EMAIL_TO=<RECEIVER>
PASSWORD=<YOUR PASSWORD (EMAIL_FROM)>
```

## Usage

Now you run the code through terminal:

```bash
go run main.go
```

or build the binary:

```bash
go build main.go
```

Then it should be done! The default port is `:3000`, but you can always change it by changing the `port` variable in the code, line `47`.
