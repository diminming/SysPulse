package com.syspulse.nmon;

import org.apache.commons.io.monitor.FileAlterationListener;
import org.apache.commons.io.monitor.FileAlterationMonitor;
import org.apache.commons.io.monitor.FileAlterationObserver;

import java.io.File;

public class NMONFileHandler {

    private static final Long INTERVAL = 1000L;

    private final FileAlterationMonitor monitor = new FileAlterationMonitor(INTERVAL);

    public NMONFileHandler(String path) {
        this(path, new NMONFileListener());
    }

    public NMONFileHandler(String path, FileAlterationListener listener) {
        FileAlterationObserver observer = new FileAlterationObserver(new File(path));
        monitor.addObserver(observer);
        observer.addListener(listener);
    }

    public void start() throws Exception {
        monitor.start();
    }

    public void stop() throws Exception {
        monitor.stop();
    }

}
