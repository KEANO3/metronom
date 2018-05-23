API Routes:
/v1
/v1/{password count}
to request 50 passwords at once, you would GET /v1/50

Every route accepts the following parameters to modify the password generation:
max: Limits the passwords length (Default: 12)
min: The passwords length will at least be this many characters (Default: 8)
num: Specifies how many numbers should be included (Default: random)
spec: Specifies how many special characters should be included (Default: random)
mode: Randomly swaps vowels with numbers, if set to 1337 (swapped vowels don't count against the number limitation)

There's a Dockerfile and docker-compose definition for easier usability. To get started just execute "docker-compose up" in the project's directory.
Afterwords you could GET localhost:8080/v1/20?max=20&min=10&mode=1337

The project is split into two parts for better reusability:
- a password generating library
- and a REST server


Examples:
curl -X GET "localhost:8080/v1"
will get you a single password with default parameters (8-12 characters, random amount of numbers and special characters)

curl -X GET "localhost:8080/v1/200"
will get you 200 passwords with default parameters

curl -X GET "localhost:8080/v1?min=200&max=0"
will get you a single password with exactly 200 characters

curl -X GET "localhost:8080/v1?min=10&max=20&spec=30"
will get you an error because you requested more special characters than your maximum limit allows