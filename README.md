# dynamic-dns

This is an ultra-minimalist dynamic dns self-update service. My CenturyLink
router cannot speak Google's protocol for unknown reasons, so I rolled my own.

When run, the program fetches its own public IP, and if the name is not already
associated with that IP, it sends it to the DNS service before exiting.

## Installation

I've included a Kubernetes CronJob specification to ensure DNS is properly
configured as frequently as desired.

Ensure you have a working cluster and [`ko`](github.com/google/ko) installation.

1. Edit `conf/secrets.default.yaml` as needed. The request URL is set to
   `https://domains.google.com/nic/update`; all other values are placeholders.
2. Edit the cron schedule in `conf/cronjob.yaml` as desired.
3. Run `ko apply -f conf/`.

Disclaimer: I wrote this service using Google's [example HTTP
request](https://support.google.com/domains/answer/6147083?hl=en#zippy=%2Cuse-the-api-to-update-your-dynamic-dns-record).
I can't guarantee that this is a well behaved client.
