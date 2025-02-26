package com.syspulse.nmon.httphandler;

import com.google.common.collect.ImmutableList;
import com.google.common.collect.ImmutableMap;
import com.google.gson.Gson;
import com.google.gson.reflect.TypeToken;
import com.ibm.nmon.data.DataRecord;
import com.ibm.nmon.data.DataType;
import com.ibm.nmon.data.NMONDataSet;
import com.ibm.nmon.data.ProcessDataType;
import com.ibm.nmon.parser.NMONParser;

import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpHandler;
import java.io.File;
import java.io.IOException;
import java.lang.reflect.Type;
import java.util.*;
import java.util.stream.Collectors;

public class NMONHttpHandler {

    // @GetMapping("/categories")
    public static HttpHandler getCategoriesHandler() {
        return new BaseHttpHandler() {

            private String getFilePath(HttpExchange httpExchange) throws IOException {
                String requestBody = getRequestBody(httpExchange);

                Type empMapType = new TypeToken<Map<String, String>>() {
                }.getType();

                Map<String, String> args = new Gson().fromJson(requestBody, empMapType);
                if (args.containsKey("filePath")) {
                    return args.get("filePath");
                }
                throw new RuntimeException("There is no key named filePath in request body.");
            }

            @Override
            public void handle(HttpExchange httpExchange) {

                try {
                    Map<String, Map<String, Object>> categories = new HashMap<>();
                    Map<String, Map<String, Object>> processLst = new HashMap<>();

                    String filePath = getFilePath(httpExchange);
                    File file = new File(filePath);
                    NMONParser parser = new NMONParser();
                    NMONDataSet dataSet = parser.parse(file, TimeZone.getDefault(), true);
                    Iterable<DataType> types = dataSet.getTypes();
                    types.forEach(item -> {

                        String key = String.format("%s_%s", dataSet.getHostname(), item.getName());
                        Map<String, Object> map = ImmutableMap.<String, Object>builder()
                                .put("id", item.getId())
                                .put("name", item.getName())
                                .put("fields", item.getFields())
                                .build();

                        if (item instanceof ProcessDataType) {
                            if (!processLst.containsKey(key)) {
                                processLst.put(key, map);
                            }
                        } else {
                            if (!categories.containsKey(key)) {
                                categories.put(key, map);
                            }
                        }
                    });
                    String reponseBody = new Gson().toJson(ImmutableList.of(
                            categories.values().stream().sorted(Comparator.comparing(o -> ((String) o.get("name"))))
                                    .collect(Collectors.toList()),
                            processLst.values().stream().sorted(Comparator.comparing(o -> ((String) o.get("name"))))
                                    .collect(Collectors.toList())));
                    jsonResponse(httpExchange, 200, reponseBody);
                } catch (Exception err) {
                    jsonResponse(httpExchange, 500, err.getMessage());
                }

            }
        };
    }

    // @GetMapping("/data")
    public static HttpHandler getDataHandler() {

        return new BaseHttpHandler() {

            private Map<String, String> getParameters(HttpExchange httpExchange) throws IOException {
                String requestBody = getRequestBody(httpExchange);

                Type empMapType = new TypeToken<Map<String, String>>() {
                }.getType();

                return new Gson().fromJson(requestBody, empMapType);
            }

            @Override
            public void handle(HttpExchange httpExchange) throws IOException {

                Map<String, String> params = getParameters(httpExchange);
                String filePath = params.get("filePath");
                String category = params.get("category");
                String field = params.get("field");

                Map<String, List<Object[]>> result = new LinkedHashMap<>();

                File file = new File(filePath);
                NMONParser parser = new NMONParser();
                NMONDataSet dataSet = parser.parse(file, TimeZone.getDefault(), true);
                DataType type = dataSet.getType(category);
                List<String> fieldLst;
                if (field == null) {
                    fieldLst = type.getFields();
                } else {
                    fieldLst = ImmutableList.of(field);
                }
                for (DataRecord dr : dataSet.getRecords()) {
                    long timestamp = dr.getTime();
                    for (String field1 : fieldLst) {
                        List<Object[]> series = result.computeIfAbsent(field1, k -> new LinkedList<>());
                        series.add(new Object[] { timestamp, dr.getData(type, field1) });
                    }
                }

                for (Map.Entry<String, List<Object[]>> entry : result.entrySet()) {
                    List<Object[]> values = entry.getValue();
                    values.sort((Object[] array1, Object[] array2) -> {
                        Long t1 = (Long) array1[0];
                        Long t2 = (Long) array2[0];
                        return t1.compareTo(t2);
                    });
                }
                String responseString = new Gson().toJson(result);
                jsonResponse(httpExchange, 200, responseString);
            }

        };

    }

    // @DeleteMapping
    public static HttpHandler getDeleteHandler() {
        return new BaseHttpHandler() {

            private String getFilePath(HttpExchange httpExchange) throws IOException {
                String requestBody = getRequestBody(httpExchange);

                Type empMapType = new TypeToken<Map<String, String>>() {
                }.getType();

                Map<String, String> args = new Gson().fromJson(requestBody, empMapType);
                if (args.containsKey("filePath")) {
                    return args.get("filePath");
                }
                throw new RuntimeException("There is no key named filePath in request body.");
            }

            @Override
            public void handle(HttpExchange httpExchage) throws IOException {
                String filePath = getFilePath(httpExchage);
                new File(filePath).delete();
                jsonResponse(httpExchage, 200, null);
            }
            
        };
    }

}
