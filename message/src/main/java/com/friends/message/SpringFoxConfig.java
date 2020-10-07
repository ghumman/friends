package com.friends.message;

// import org.springframework.context.annotation.Bean;
// import org.springframework.context.annotation.Configuration;

// import springfox.documentation.builders.PathSelectors;
// import springfox.documentation.builders.RequestHandlerSelectors;
// import springfox.documentation.spi.DocumentationType;
// import springfox.documentation.spring.web.plugins.Docket;

// @Configuration
// public class SpringFoxConfig {
//     @Bean
//     public Docket api() { 
//         return new Docket(DocumentationType.SWAGGER_2)  
//           .select()                                  
//           .apis(RequestHandlerSelectors.any())              
//           .paths(PathSelectors.any())                          
//           .build();                                           
//     }
// }

import com.fasterxml.classmate.TypeResolver;
import com.friends.message.Message.MessageController;
import com.friends.message.User.UserController;

// import org.joda.time.LocalDate;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;
import org.springframework.http.HttpMethod;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.context.request.async.DeferredResult;
import springfox.documentation.builders.PathSelectors;
import springfox.documentation.builders.RequestHandlerSelectors;
import springfox.documentation.builders.ResponseBuilder;
import springfox.documentation.schema.ScalarType;
import springfox.documentation.schema.WildcardType;
import springfox.documentation.service.ApiKey;
import springfox.documentation.service.AuthorizationScope;
import springfox.documentation.service.ParameterType;
import springfox.documentation.service.SecurityReference;
import springfox.documentation.service.Tag;
import springfox.documentation.spi.DocumentationType;
import springfox.documentation.spi.service.contexts.SecurityContext;
import springfox.documentation.spring.web.plugins.Docket;
import springfox.documentation.swagger.web.DocExpansion;
import springfox.documentation.swagger.web.ModelRendering;
import springfox.documentation.swagger.web.OperationsSorter;
import springfox.documentation.swagger.web.SecurityConfiguration;
import springfox.documentation.swagger.web.SecurityConfigurationBuilder;
import springfox.documentation.swagger.web.TagsSorter;
import springfox.documentation.swagger.web.UiConfiguration;
import springfox.documentation.swagger.web.UiConfigurationBuilder;
import springfox.documentation.swagger2.annotations.EnableSwagger2;
// import springfox.petstore.controller.PetController;

import java.time.LocalDate;
import java.util.List;

import static java.util.Collections.*;
import static springfox.documentation.schema.AlternateTypeRules.*;

// Enables Springfox swagger 2
@EnableSwagger2 

// Instructs spring where to scan for API controllers
@ComponentScan(basePackageClasses = {
    UserController.class, MessageController.class
})
public class SpringFoxConfig {

  public static void main(String[] args) {
    SpringApplication.run(MessageApplication.class, args);
  }


  @Bean
  public Docket petApi() {
    // Docket, Springfox’s, primary api configuration mechanism is initialized for swagger specification 2.0
    return new Docket(DocumentationType.SWAGGER_2)
        // select() returns an instance of ApiSelectorBuilder to give fine grained control over the endpoints exposed via swagger.
        .select() 
        // apis() allows selection of RequestHandler's using a predicate. The example here uses an any predicate (default). Out of the box predicates provided are any, none, withClassAnnotation, withMethodAnnotation and basePackage.
        .apis(RequestHandlerSelectors.any()) 
        // paths() allows selection of Path's using a predicate. The example here uses an any predicate (default). Out of the box we provide predicates for regex, ant, any, none.
        .paths(PathSelectors.any()) 
        // 	The selector needs to be built after configuring the api and path selectors.
        .build() 
        // Adds a servlet path mapping, when the servlet has a path mapping. This prefixes paths with the provided path mapping.
        .pathMapping("/") 
        // Convenience rule builder that substitutes LocalDate with String when rendering model properties
        .directModelSubstitute(LocalDate.class, String.class)
        .genericModelSubstitutes(ResponseEntity.class)
        // 	Convenience rule builder that substitutes a generic type with one type parameter with the type parameter. In this example ResponseEntity<T> with T. alternateTypeRules allows custom rules that are a bit more involved. The example substitutes DeferredResult<ResponseEntity<T>> with T generically.
        .alternateTypeRules(
            newRule(typeResolver.resolve(DeferredResult.class,
                typeResolver.resolve(ResponseEntity.class, WildcardType.class)),
                typeResolver.resolve(WildcardType.class))) 
        // 	Flag to indicate if default http response codes need to be used or not.
        .useDefaultResponseMessages(false) 
        // Allows globally overriding response messages for different http methods. In this example we override the 500 error code for all GET requests …​
        .globalResponses(HttpMethod.GET, 
            singletonList(new ResponseBuilder()
                .code("500")
                .description("500 message")
                .representation(MediaType.TEXT_XML)
                .apply(r ->
                    r.model(m ->
                        m.referenceModel(ref ->
                            ref.key(k ->
                                k.qualifiedModelName(q ->
                                    q.namespace("some:namespace")
                                        // 	…​ and indicate that it will use the response model Error (which will be defined elsewhere)
                                        .name("ERROR")))))) 
                .build()))
        // 	Sets up the security schemes used to protect the apis. Supported schemes are ApiKey, BasicAuth and OAuth
        .securitySchemes(singletonList(apiKey())) 
        // Provides a way to globally set up security contexts for operation. The idea here is that we provide a way to select operations to be protected by one of the specified security schemes.
        .securityContexts(singletonList(securityContext())) 
        /*
        * Incubating * setting this flag signals to the processor that the paths generated should try and use form style query expansion. 
        As a result we could distinguish paths that have the same path stem but different query string combinations. 
        An example of this would be two apis: First, http://example.org/findCustomersBy?name=Test to find customers by name. Per RFC 6570, 
        this would be represented as http://example.org/findCustomersBy{?name}. Second, http://example.org/findCustomersBy?zip=76051 to find customers by zip. 
        Per RFC 6570, this would be represented as http://example.org/findCustomersBy{?zip}.
        */
        .enableUrlTemplating(true) 
        /*
        Allows globally configuration of default path-/request-/headerparameters which are common for every rest operation of the api, 
        but aren`t needed in spring controller method signature (for example authenticaton information). 
        Parameters added here will be part of every API Operation in the generated swagger specification. on 
        how the security is setup the name of the header used may need to be different. 
        Overriding this value is a way to override the default behavior.
        */
        .globalRequestParameters(
            singletonList(new springfox.documentation.builders.RequestParameterBuilder()
                .name("someGlobalParameter")
                .description("Description of someGlobalParameter")
                .in(ParameterType.QUERY)
                .required(true)
                .query(q -> q.model(m -> m.scalarModel(ScalarType.STRING)))
                .build()))
        // 	Adding tags is a way to define all the available tags that services/operations can opt into. Currently this only has name and description.
        .tags(new Tag("Pet Service", "All apis relating to pets"))
        /* Are there models in the application that are not "reachable"? 
        Not reachable is when we have models that we would like to be described but 
        aren’t explicitly used in any operation. 
        An example of this is an operation that returns a model serialized as a string. 
        We do want to communicate the expectation of the schema for the string. 
        This is a way to do exactly that.
        */
        //.additionalModels(typeResolver.resolve(AdditionalModel.class)); 
        ;
  }

  @Autowired
  private TypeResolver typeResolver;

  // Here we use ApiKey as the security schema that is identified by the name mykey
  private ApiKey apiKey() {
    return new ApiKey("mykey", "api_key", "header"); 
  }

  private SecurityContext securityContext() {
    return SecurityContext.builder()
        .securityReferences(defaultAuth())
        // Selector for the paths this security context applies to.
        // A.Ghumman: it is deprecated
        // .forPaths(PathSelectors.regex("/anyPath.*"))
        .build();
  }

  List<SecurityReference> defaultAuth() {
    AuthorizationScope authorizationScope
        = new AuthorizationScope("global", "accessEverything");
    AuthorizationScope[] authorizationScopes = new AuthorizationScope[1];
    authorizationScopes[0] = authorizationScope;
    return singletonList(
        // Here we use the same key defined in the security scheme mykey
        new SecurityReference("mykey", authorizationScopes)); 
  }

  @Bean
  SecurityConfiguration security() {
    // Optional swagger-ui security configuration for oauth and apiKey settings
    return SecurityConfigurationBuilder.builder() 
        .clientId("test-app-client-id")
        .clientSecret("test-app-client-secret")
        .realm("test-app-realm")
        .appName("test-app")
        .scopeSeparator(",")
        .additionalQueryStringParams(null)
        .useBasicAuthenticationWithAccessCodeGrant(false)
        .enableCsrfSupport(false)
        .build();
  }

  @Bean
  UiConfiguration uiConfig() {
    // Optional swagger-ui ui configuration currently only supports the validation url
    return UiConfigurationBuilder.builder() 
        .deepLinking(true)
        .displayOperationId(false)
        .defaultModelsExpandDepth(1)
        .defaultModelExpandDepth(1)
        .defaultModelRendering(ModelRendering.EXAMPLE)
        .displayRequestDuration(false)
        .docExpansion(DocExpansion.NONE)
        .filter(false)
        .maxDisplayedTags(null)
        .operationsSorter(OperationsSorter.ALPHA)
        .showExtensions(false)
        .showCommonExtensions(false)
        .tagsSorter(TagsSorter.ALPHA)
        .supportedSubmitMethods(UiConfiguration.Constants.DEFAULT_SUBMIT_METHODS)
        .validatorUrl(null)
        .build();
  }

}
