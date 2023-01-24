package com.friends.message;

import com.amazonaws.serverless.exceptions.ContainerInitializationException;
import com.amazonaws.serverless.proxy.internal.testutils.Timer;
import com.amazonaws.serverless.proxy.model.AwsProxyRequest;
import com.amazonaws.serverless.proxy.model.AwsProxyResponse;
import com.amazonaws.serverless.proxy.spring.SpringBootLambdaContainerHandler;
import com.amazonaws.services.lambda.runtime.Context;
import com.amazonaws.services.lambda.runtime.RequestStreamHandler;

import javax.servlet.DispatcherType;
import javax.servlet.FilterRegistration;

import java.io.IOException;
import java.io.InputStream;
import java.io.OutputStream;
import java.util.EnumSet;


import com.friends.message.MessageApplication;

public class StreamLambdaHandler implements RequestStreamHandler {

    public static final SpringBootLambdaContainerHandler<AwsProxyRequest, AwsProxyResponse> handler;
 
    static { 
       try { 
           handler = SpringBootLambdaContainerHandler.getAwsProxyHandler(MessageApplication.class);
           handler.activateSpringProfiles("lambda");
      } catch (ContainerInitializationException e) { 
           throw new RuntimeException("Could not initialize Spring Boot application", e); 
       } 
    }

    @Override 
    public void handleRequest(InputStream inputStream, OutputStream outputStream, Context context) throws IOException {
        handler.proxyStream(inputStream, outputStream, context);
        // just in case it wasn't closed 
        outputStream.close(); 
    }
}
