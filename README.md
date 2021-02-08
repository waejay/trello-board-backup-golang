## trello-board-backup-golang

I created this script because my Trello board (which is my primary life/finance/hobby/whatever management tool) contains a lot of personal and valuable data. If Trello suddenly disappeared, a small chunk of my life would disappear too.
  - In order to preserve that data, this script fetches the Trello board's
    content in JSON form using Trello's RESTful API and uploads it into my
personal free-tier S3 bucket.
  - I uploaded this script on a Lambda function that is triggered by CloudWatch
    to invoke every month, so my data should be backed up every month for the
rest of my life).
