package com.syspulse.nmon.httphandler;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;

public abstract class BaseHttpHandler implements HttpHandler {
    protected void jsonResponse(HttpExchange httpExchange, int respCode, String response) {
        try {
            byte[] responseContentByte = response == null ? new byte[0] : response.getBytes("utf-8");

            // 设置响应头，必须在sendResponseHeaders方法之前设置！
            httpExchange.getResponseHeaders().add("Content-Type:", "application/json;charset=utf-8");
    
            // 设置响应码和响应体长度，必须在getResponseBody方法之前调用！
            httpExchange.sendResponseHeaders(respCode, responseContentByte.length);
    
            OutputStream out = httpExchange.getResponseBody();
            out.write(responseContentByte);
            out.flush();
            out.close();
        } catch(IOException exp) {
            exp.printStackTrace();
        }
        
    }

    protected String getRequestBody(HttpExchange httpExchange) throws IOException {
        String reqBody = "";

        // 非GET请求读请求体
        BufferedReader bufferedReader = new BufferedReader(
                new InputStreamReader(httpExchange.getRequestBody(), "utf-8"));
        StringBuilder requestBodyContent = new StringBuilder();
        String line = null;
        while ((line = bufferedReader.readLine()) != null) {
            requestBodyContent.append(line);
        }
        reqBody = requestBodyContent.toString();

        return reqBody;
    }
}
