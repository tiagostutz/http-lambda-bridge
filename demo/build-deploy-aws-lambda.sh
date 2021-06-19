#!/bin/sh

echo "Create ECR repository using Cloudformation..."
aws cloudformation deploy --template-file cf-ecr.yml --stack-name http-lambda-bridge-demo-ecr --capabilities CAPABILITY_NAMED_IAM

repoUri=$(aws cloudformation describe-stacks --stack-name http-lambda-bridge-demo-ecr --query Stacks[].[Outputs[].OutputValue] --output text)

echo "ECR repo is $repoUri"

echo "Building container image and tagging with $repoUri..."
docker build -t $repoUri .

echo "Logging docker cli to ECR repo"
repoDomain=$(echo $repoUri | cut -d/ -f1)
aws ecr --no-verify-ssl get-login-password | docker login --username AWS --password-stdin $repoDomain

echo "Pushing container image to repo..."
version=$(docker inspect --format='{{index .Id}}' $repoUri | cut -d\: -f2)
repoUriTag=$repoUri\:$version
docker tag $repoUri $repoUriTag
docker push $repoUriTag

echo "Deploying AWS API Gateway and Lambda Functions with Cloudformation..."
aws cloudformation deploy \
    --template-file cf-api-lambda.yml \
    --stack-name http-lambda-bridge-demo-service \
    --parameter-overrides LambdaContainerImageUri=$repoUriTag \
    --capabilities CAPABILITY_NAMED_IAM

sleep 3
out=$(aws cloudformation describe-stacks --stack-name http-lambda-bridge-demo-service --query Stacks[].[Outputs[].OutputValue] --output text)
apiUri=$(echo $out | cut -f1)

echo "HTTP Lambda Bridge Demo API is running at $apiUri"
set +x
# curl $apiUri/repo