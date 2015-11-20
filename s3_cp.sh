#!/bin/bash

export AWS_ACCESS_KEY_ID=AKIAOLKVAPHKSE7ATAQQ
export AWS_SECRET_ACCESS_KEY=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
export AWS_DEFAULT_REGION=cn-north-1
S3_BUCKET="asiainfoldp-file-backup"
print `ls $1`
for host in `ls $1` ;
do aws s3 cp $host s3://$S3_BUCKET/$host
done

