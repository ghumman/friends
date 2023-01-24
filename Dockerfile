#FROM public.ecr.aws/lambda/java:11.2023.01.11.07
# ARG JAR_FILE=target/*.jar
# COPY ${JAR_FILE} app.jar
# ENTRYPOINT ["java","-jar","/app.jar"]
# COPY . ${LAMBDA_TASK_ROOT}/
#CMD ["sh", "script.sh"]

#FROM public.ecr.aws/lambda/java:11

# Copy function code and runtime dependencies from Maven layout
#COPY target/classes ${LAMBDA_TASK_ROOT}
#COPY target/dependency/* ${LAMBDA_TASK_ROOT}/lib/

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
#CMD ["sh", "script.sh"]
# CMD [ "com.example.LambdaHandler::handleRequest" ]

#FROM public.ecr.aws/lambda/java:11

# Copy function code and runtime dependencies from Maven layout
#COPY target/classes ${LAMBDA_TASK_ROOT}
#COPY target/dependency/* ${LAMBDA_TASK_ROOT}/lib/

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
#CMD [ "com.friends.message::handleRequest" ]

FROM public.ecr.aws/lambda/java:11

# Copy function code and runtime dependencies from Maven layout
COPY target/classes ${LAMBDA_TASK_ROOT}
COPY target/dependency/* ${LAMBDA_TASK_ROOT}/lib/

# Set the CMD to your handler (could also be done as a parameter override outside of the Dockerfile)
CMD [ "com.friends.message.StreamLambdaHandler::handleRequest" ]

