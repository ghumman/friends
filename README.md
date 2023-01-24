# About

This branch shares experience of using AWS ECS Elastic Container Service for hosting docker images. For setting up aws credentials and associating them with docker, creating and uploading docker image to public and private ECR Elastic Container Registry, use `README` from branch `aws-lambda-docker`. 

In order to setup Java Docker Container on ECS, I used following guide. And it works really good. 
```
https://mydeveloperplanet.com/2021/09/07/how-to-deploy-a-spring-boot-app-on-aws-ecs-cluster/
```

I tried multiple things and guide but it was not working initially, and the thing I was donig wrong was that when creating the cluster, under networking, I was not selecting all the subnets in a vpc and that's why the whole process used to go smooth but I was not able to hit the endpoints. For other branches I was creating the docker image by running something like `docker build -t ghumman/java-backend .`. And then after that I was testing the image locally and pushing it the repositories. That should work too but following the above guide. I added following code to the `pom.xml`. 
```
<plugin>
    <groupId>com.spotify</groupId>
    <artifactId>dockerfile-maven-plugin</artifactId>
    <version>1.4.12</version>
    <executions>
        <execution>
            <id>default</id>
            <goals>
                <goal>build</goal>
                <goal>push</goal>
            </goals>
        </execution>
    </executions>

    <configuration>
        <repository>ghumman/java-backend</repository>
        <tag>${project.version}</tag>
        <buildArgs>
            <JAR_FILE>target/message-0.0.1-SNAPSHOT.jar</JAR_FILE>
        </buildArgs>
    </configuration>
</plugin>

```

And then run `mvn clean verify`. This will automatically create docker image named `ghumman/java-backend` for you. And then use the same procedure mentioned in the README of branch `aws-lambda-docker` to upload it to ECR. 
