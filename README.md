# APKiD-bot

A simple Slack bot which wraps APKiD (the docker script at least) and scans for uploads which might be interesting to scan.

This is deployed to https://slack.rednaga.io - unsure if anyone else would want to use it, but here is the code!

## Installing

Pull this code recursively, then build the `APKiD` docker container. Then modify anything specific for your env (I have not abstracted all configs to environment variables) and compile the golang code.

```
diff@milo:~/repo/apkid-bot$ cd APKiD/
diff@milo:~/repo/apkid-bot/APKiD$ docker-compose build
Building apkid
Step 1/15 : FROM python:2.7-slim
 ---> e9adbdab327d
Step 2/15 : RUN apt-get update -qq && apt-get install -y git build-essential gcc pandoc
 ---> Using cache
 ---> c81dea894ee9
Step 3/15 : RUN pip install --upgrade pip
 ---> Using cache
...
...
Successfully built 525f01e5d5a6
Successfully tagged rednaga/apkid:v1
diff@milo:~/repo/apkid-bot/APKiD$ cd ..
diff@milo:~/repo/apkid-bot$ go build .
```

## Usage

Run the `run.sh` with your own Slack token.

```
diff@milo:~/repo/apkid-bot$ ./run.sh
```