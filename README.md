# About

This branch explains the steps to host a docker image to AWS ECR - Elastic Container Registry. Also how to setup lambda function to use that ECR image. 

## Pre-requisites

You need to have aws cli and docker cli installed on your computer. 

Guide to install/upgrade aws cli can be found here. 
https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html

The complete guide for creating lambda container images and deploying lambda function as container image can be found here. 

https://docs.aws.amazon.com/lambda/latest/dg/images-create.html

https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-images.html

## Creating ECR repository
Go to ECR on AWS and create new repository or you can use command line to create new repository. Public repos within limits do not cost anything. 

You can have either public repository or private repository. Following are the complete guides.

For Amazon ECR Public Registries

https://docs.aws.amazon.com/AmazonECR/latest/public/public-registries.html#public-registry-auth

For Amazon ECR Private Registries

https://docs.aws.amazon.com/lambda/latest/dg/images-create.html

- General Setup for both private and public registries

Go to AWS IAM. Create a User. If you have already created an group, add this user to the group otherwise create a new group with admin previliges. 

Go to that user IAM > Users. Under security credentials, create access key. While creating access key you will get access key and secret and also be able to download csv which has these 2 fields. 

Run following to create profile named `personal-ghummantech`. This profile will be save at `~/.aws/credentials`. We are using a profile name as we have other profiles to access company aws accounts. 
``` 
aws configure --profile personal-ghummantech
```

Following is an example of `aws configure` taken from aws guide. 
```
AWS Access Key ID [None]: AKIAIOSFODNN7EXAMPLE
AWS Secret Access Key [None]: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
Default region name [None]: us-east-1
Default output format [None]: json
```

- Private Repository Setup

Note down your aws account id. You can get it when you click your username in aws. 

Before using following replace `123456789012` with your aws account id. 
```
aws ecr get-login-password --region us-east-1 --profile personal-ghummantech | docker login --username AWS --password-stdin 123456789012.dkr.ecr.us-east-1.amazonaws.com    
```

Create the repository like following: 
```
aws ecr create-repository --repository-name java-backend --image-scanning-configuration scanOnPush=true --image-tag-mutability MUTABLE
```
Here `java-backend` is the name of repository. 

Tag and push local repository to private registry. Change `123456789012` with your aws account id. 
```
docker tag  ghumman/java-backend:latest 123456789012.dkr.ecr.us-east-1.amazonaws.com/java-backend:latest
docker push 123456789012.dkr.ecr.us-east-1.amazonaws.com/java-backend:latest        
```

- Public Repository Setup

One time setup to authenticate docker to public ecr account. 
```
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
```

Tag the local repo. Currenly I have `m6a5y9u9` as username. I requested `ghumman` which is in review. 
```
docker tag ghumman/java-backend:latest public.ecr.aws/m6a5y9u9/java-backend
```

Push
```
docker push public.ecr.aws/m6a5y9u9/java-backend
```