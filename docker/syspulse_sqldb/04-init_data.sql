USE `syspulse`;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (1,'admin','admin1!',1,0,0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;

create index idx_perf_cpu_linux_id_timestamp on perf_cpu(linux_id, timestamp);
create index idx_perf_disk_io_linux_id_timestamp on perf_disk_io(linux_id, timestamp);
create index idx_perf_fs_usage_linux_id_timestamp on perf_fs_usage(linux_id, timestamp);
create index idx_perf_if_io_linux_id_timestamp on perf_if_io(linux_id, timestamp);
create index idx_perf_load_linux_id_timestamp on perf_load(linux_id, timestamp);
create index idx_perf_mem_linux_id_timestamp on perf_mem(linux_id, timestamp);
create index idx_perf_swap_linux_id_timestamp on perf_swap(linux_id, timestamp);