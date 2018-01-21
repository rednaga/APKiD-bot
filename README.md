# APKiD-bot

A simple Slack bot which wraps APKiD (the docker script at least) and scans for uploads which might be interesting to scan.

This is deployed to https://slack.rednaga.io - unsure if anyone else would want to use it, but here is the code!

## Installing

Pull this code recursively, then build the `APKiD` docker container. Then modify anything specific for your env (I have not abstracted all configs to environment variables) and compile the golang code.

## Usage

Run the `run.sh` with your own Slack token.