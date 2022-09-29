FROM golang:1.18.1 as go-build
WORKDIR /app
COPY ./ ./
RUN apt-get update -q && apt-get install -y zip
RUN cd ./fs && go mod tidy && make build
RUN cd ./iam && go mod tidy && make build


FROM node:16
# install aws cdk
RUN npm install -g aws-cdk@1.171.0

#install aws cli
RUN apt-get update -q
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-aarch64.zip" -o "awscliv2.zip"
RUN unzip awscliv2.zip
RUN ./aws/install -i /usr/local/aws-cli -b /usr/local/bin
RUN aws --version
RUN aws configure set aws_access_key_id xxxxx
RUN aws configure set aws_secret_access_key xxxxx
RUN aws configure set region us-east-1
WORKDIR /app
COPY --from=go-build ./app ./app
RUN cdk bootstrap aws://xxxxx/us-east-1

