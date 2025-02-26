package com.syspulse.nmon;

import com.sun.net.httpserver.HttpServer;
import com.syspulse.nmon.httphandler.NMONHttpHandler;

import java.io.IOException;
import java.net.InetSocketAddress;
import java.util.concurrent.Executors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


public class NMONServer {

    private static final Logger LOGGER = LoggerFactory.getLogger(NMONServer.class);

    public static void main(String[] args) throws IOException {

        String nmonDir = "/data/syspulse/nmon";
        int serverPort = 8000;

        new Thread(new Runnable() {

            @Override
            public void run() {
                try {
                    new NMONFileHandler(nmonDir).start();
                } catch (Exception e) {
                    LOGGER.error("error init file listener thread.", e);
                }
            }
            
        }).start();

        //创建一个HttpServer实例，并绑定到指定的IP地址和端口号
        HttpServer httpServer = HttpServer.create(new InetSocketAddress("0.0.0.0", serverPort), 0);

        //创建一个HttpContext，将请求映射到MyHttpHandler处理器
        httpServer.createContext("/nmon/categories", NMONHttpHandler.getCategoriesHandler());

        httpServer.createContext("/nmon/data", NMONHttpHandler.getDataHandler());

        httpServer.createContext("/nmon/del", NMONHttpHandler.getDeleteHandler());

        //设置服务器的线程池对象
        httpServer.setExecutor(Executors.newFixedThreadPool(10));

        //启动服务器
        httpServer.start();
    }
}
