package com.syspulse.nmon;

import com.google.gson.Gson;
import com.ibm.nmon.data.NMONDataSet;
import com.ibm.nmon.parser.NMONParser;
import com.syspulse.common.HttpUtil;

import org.apache.commons.io.monitor.FileAlterationListenerAdaptor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.HashMap;
import java.util.Map;
import java.util.TimeZone;

public class NMONFileListener extends FileAlterationListenerAdaptor {

    private static final Logger LOGGER = LoggerFactory.getLogger(NMONFileListener.class);

    private void wait4Complete(File file) throws InterruptedException {
        long fileLength0;
        do {
            fileLength0 = file.length();
            Thread.sleep(300);
        } while (fileLength0 != file.length());
    }

    private String getMD5(File file) {
        BigInteger bi = null;
        try {
            byte[] buffer = new byte[8192];
            int len = 0;
            MessageDigest md = MessageDigest.getInstance("MD5");
            FileInputStream fis = new FileInputStream(file);
            while ((len = fis.read(buffer)) != -1) {
                md.update(buffer, 0, len);
            }
            fis.close();
            byte[] b = md.digest();
            bi = new BigInteger(1, b);
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (IOException e) {
            e.printStackTrace();
        }
        return bi.toString(16);
    }

    @Override
    public void onFileCreate(File file) {
        try {

            wait4Complete(file);
            String md5 = getMD5(file);
            LOGGER.debug("nmon file: {}, MD5: {}", file.getAbsoluteFile(), md5);

            NMONParser parser = new NMONParser();

            NMONDataSet dataSet = parser.parse(file, TimeZone.getDefault(), true);
            String hostname = dataSet.getHostname();
            Long startTime = dataSet.getStartTime();
            Long endTime = dataSet.getEndTime();
            String path = file.getPath();

            Map<String, Object> map = new HashMap<>();
            map.put("hostname", hostname);
            map.put("from", startTime);
            map.put("to", endTime);
            map.put("path", path);

            LOGGER.debug("parsed: {}, {}, {}, {}", hostname, startTime, endTime, path);
            HttpUtil.sendPostRequest("http://localhost:24162/api/nmon", new Gson().toJson(map));

        } catch (Exception e) {
            LOGGER.error(e.getMessage(), e);
        }
    }

}
