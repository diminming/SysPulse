package com.syspulse.common;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;

public class HttpUtil {

    static final HttpClient client = HttpClient.newHttpClient();

    public static int sendPostRequest(String url, String jsonStr) throws InterruptedException, ExecutionException {

        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .POST(HttpRequest.BodyPublishers.ofString(jsonStr))
                .build();
        CompletableFuture<HttpResponse<String>> futureResponse = client.sendAsync(request,
                HttpResponse.BodyHandlers.ofString());

        HttpResponse<String> response = futureResponse.get();
        assert response.statusCode() == 200 : "response code is not 200";
        assert response.body().equals("{\"message\":\"ok\"}") : "response body is not {\"message\":\"ok\"}";
        return -1;
    }

}
