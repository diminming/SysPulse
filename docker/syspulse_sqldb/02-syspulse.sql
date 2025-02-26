USE `syspulse`
-- MySQL dump 10.13  Distrib 8.4.3, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: syspulse
-- ------------------------------------------------------
-- Server version	8.4.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `alarm`
--

DROP TABLE IF EXISTS `alarm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `alarm` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `timestamp` bigint NOT NULL,
  `time_tag` char(10) NOT NULL,
  `linux_id` bigint NOT NULL,
  `biz_id` bigint NOT NULL,
  `trigger_id` varchar(100) NOT NULL,
  `trigger` varchar(500) NOT NULL,
  `level` varchar(45) NOT NULL,
  `ack` tinyint(1) NOT NULL,
  `source` varchar(45) NOT NULL,
  `msg` varchar(200) NOT NULL,
  `create_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=620 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `biz`
--

DROP TABLE IF EXISTS `biz`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `biz` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `biz_name` varchar(100) NOT NULL,
  `biz_id` varchar(200) NOT NULL,
  `biz_desc` varchar(5000) NOT NULL,
  `create_timestamp` bigint NOT NULL,
  `update_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1256 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cfg_cpu`
--

DROP TABLE IF EXISTS `cfg_cpu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cfg_cpu` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `cpu` int NOT NULL,
  `vendorid` varchar(500) NOT NULL,
  `family` varchar(500) NOT NULL,
  `model` varchar(500) NOT NULL,
  `stepping` int NOT NULL,
  `physicalid` varchar(500) NOT NULL,
  `coreid` varchar(500) NOT NULL,
  `cores` int NOT NULL,
  `modelname` varchar(500) NOT NULL,
  `mhz` double NOT NULL,
  `cachesize` int NOT NULL,
  `microcode` varchar(500) NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=65 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cfg_host`
--

DROP TABLE IF EXISTS `cfg_host`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cfg_host` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `hostname` varchar(200) NOT NULL,
  `uptime` bigint NOT NULL,
  `boottime` bigint NOT NULL,
  `procs` bigint NOT NULL,
  `os` varchar(200) NOT NULL,
  `platform` varchar(200) NOT NULL,
  `platformfamily` varchar(200) NOT NULL,
  `platformversion` varchar(200) NOT NULL,
  `kernelversion` varchar(200) NOT NULL,
  `kernelarch` varchar(200) NOT NULL,
  `virtualizationsystem` varchar(200) NOT NULL,
  `virtualizationrole` varchar(200) NOT NULL,
  `hostid` varchar(200) NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cfg_if`
--

DROP TABLE IF EXISTS `cfg_if`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cfg_if` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `index` int NOT NULL,
  `name` varchar(200) NOT NULL,
  `addr` varchar(200) NOT NULL,
  `hard_addr` varchar(200) NOT NULL,
  `mtu` int NOT NULL,
  `flags` varchar(200) NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `conn_lst`
--

DROP TABLE IF EXISTS `conn_lst`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `conn_lst` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `family` int NOT NULL,
  `type` int NOT NULL,
  `l_ip` varchar(15) NOT NULL,
  `l_port` int NOT NULL,
  `r_ip` varchar(15) NOT NULL,
  `r_port` int NOT NULL,
  `status` varchar(20) NOT NULL,
  `pid` int NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `db_record`
--

DROP TABLE IF EXISTS `db_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_record` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `name` varchar(40) NOT NULL COMMENT 'database name',
  `db_id` varchar(40) NOT NULL COMMENT 'identifier of DB',
  `type` varchar(40) NOT NULL COMMENT 'database type, MySQL, PGSql, Oracle, DB2...',
  `biz_id` int NOT NULL,
  `linux_id` int NOT NULL,
  `create_timestamp` bigint NOT NULL,
  `update_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `job`
--

DROP TABLE IF EXISTS `job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `job` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_name` varchar(200) NOT NULL,
  `category` varchar(200) NOT NULL,
  `type` varchar(200) DEFAULT NULL,
  `status` int NOT NULL,
  `startup_time` bigint DEFAULT NULL,
  `linux_id` bigint DEFAULT NULL,
  `pid` bigint DEFAULT NULL,
  `duration` int DEFAULT NULL,
  `immediately` tinyint(1) DEFAULT NULL,
  `direction` varchar(100) DEFAULT NULL,
  `count` int DEFAULT NULL,
  `ip_addr` varchar(255) DEFAULT NULL,
  `if_name` varchar(255) DEFAULT NULL,
  `port` int DEFAULT NULL,
  `extend` varchar(500) DEFAULT NULL,
  `create_timestamp` bigint NOT NULL,
  `update_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=55 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `linux`
--

DROP TABLE IF EXISTS `linux`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `linux` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `hostname` varchar(200) NOT NULL,
  `linux_id` varchar(200) NOT NULL,
  `biz_id` int NOT NULL,
  `agent_conn` varchar(200) DEFAULT NULL,
  `ext_id` varchar(20) DEFAULT NULL,
  `create_timestamp` bigint NOT NULL,
  `update_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1738 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `menu`
--

DROP TABLE IF EXISTS `menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(200) NOT NULL,
  `identity` varchar(200) NOT NULL,
  `type` varchar(45) NOT NULL,
  `url` varchar(500) NOT NULL,
  `index` int NOT NULL,
  `parent` bigint NOT NULL,
  `createTimestamp` bigint NOT NULL,
  `updateTimestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_user_identity` (`identity`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `menu_permission`
--

DROP TABLE IF EXISTS `menu_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `menu_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `menu_id` bigint NOT NULL,
  `permission_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nmon`
--

DROP TABLE IF EXISTS `nmon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nmon` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `hostname` varchar(200) NOT NULL,
  `from` bigint NOT NULL,
  `to` bigint NOT NULL,
  `path` varchar(1000) DEFAULT NULL,
  `source` varchar(100) NOT NULL,
  `createTimestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_cpu`
--

DROP TABLE IF EXISTS `perf_cpu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_cpu` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `cpu` varchar(10) NOT NULL,
  `user` double NOT NULL,
  `system` double NOT NULL,
  `idle` double NOT NULL,
  `nice` double NOT NULL,
  `iowait` double NOT NULL,
  `irq` double NOT NULL,
  `softirq` double NOT NULL,
  `steal` double NOT NULL,
  `guest` double NOT NULL,
  `guestnice` double NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_cpu_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=1064610 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_disk_io`
--

DROP TABLE IF EXISTS `perf_disk_io`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_disk_io` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `readcount` bigint NOT NULL,
  `mergedreadcount` bigint NOT NULL,
  `writecount` bigint NOT NULL,
  `mergedwritecount` bigint NOT NULL,
  `readbytes` bigint NOT NULL,
  `writebytes` bigint NOT NULL,
  `readtime` bigint NOT NULL,
  `writetime` bigint NOT NULL,
  `iopsinprogress` bigint NOT NULL,
  `iotime` bigint NOT NULL,
  `weightedio` bigint NOT NULL,
  `name` varchar(200) NOT NULL,
  `serialnumber` varchar(200) NOT NULL,
  `label` varchar(200) NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_disk_io_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=2750922 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_fs_usage`
--

DROP TABLE IF EXISTS `perf_fs_usage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_fs_usage` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `path` varchar(500) NOT NULL,
  `fstype` varchar(50) NOT NULL,
  `total` bigint NOT NULL,
  `free` bigint NOT NULL,
  `used` bigint NOT NULL,
  `usedpercent` double NOT NULL,
  `inodestotal` bigint NOT NULL,
  `inodesused` bigint NOT NULL,
  `inodesfree` bigint NOT NULL,
  `inodesusedpercent` double NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_fs_usage_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=1401785 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_if_io`
--

DROP TABLE IF EXISTS `perf_if_io`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_if_io` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `name` varchar(200) NOT NULL,
  `bytessent` bigint NOT NULL,
  `bytesrecv` bigint NOT NULL,
  `packetssent` bigint NOT NULL,
  `packetsrecv` bigint NOT NULL,
  `errin` bigint NOT NULL,
  `errout` bigint NOT NULL,
  `dropin` bigint NOT NULL,
  `dropout` bigint NOT NULL,
  `fifoin` bigint NOT NULL,
  `fifoout` bigint NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_if_io_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=3789282 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_load`
--

DROP TABLE IF EXISTS `perf_load`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_load` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `load1` double NOT NULL,
  `load5` double NOT NULL,
  `load15` double NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_load_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=516382 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_mem`
--

DROP TABLE IF EXISTS `perf_mem`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_mem` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `total` bigint NOT NULL,
  `available` bigint NOT NULL,
  `used` bigint NOT NULL,
  `usedpercent` double NOT NULL,
  `free` bigint NOT NULL,
  `active` bigint NOT NULL,
  `inactive` bigint NOT NULL,
  `wired` bigint NOT NULL,
  `laundry` bigint NOT NULL,
  `buffers` bigint NOT NULL,
  `cached` bigint NOT NULL,
  `writeback` bigint NOT NULL,
  `dirty` bigint NOT NULL,
  `writebacktmp` bigint NOT NULL,
  `shared` bigint NOT NULL,
  `slab` bigint NOT NULL,
  `sreclaimable` bigint NOT NULL,
  `sunreclaim` bigint NOT NULL,
  `pagetables` bigint NOT NULL,
  `swapcached` bigint NOT NULL,
  `commitlimit` bigint NOT NULL,
  `committedas` bigint NOT NULL,
  `hightotal` bigint NOT NULL,
  `highfree` bigint NOT NULL,
  `lowtotal` bigint NOT NULL,
  `lowfree` bigint NOT NULL,
  `swaptotal` bigint NOT NULL,
  `swapfree` bigint NOT NULL,
  `mapped` bigint NOT NULL,
  `vmalloctotal` bigint NOT NULL,
  `vmallocused` bigint NOT NULL,
  `vmallocchunk` bigint NOT NULL,
  `hugepagestotal` bigint NOT NULL,
  `hugepagesfree` bigint NOT NULL,
  `hugepagesrsvd` bigint NOT NULL,
  `hugepagessurp` bigint NOT NULL,
  `hugepagesize` bigint NOT NULL,
  `anonhugepages` bigint NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_mem_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=769850 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `perf_swap`
--

DROP TABLE IF EXISTS `perf_swap`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `perf_swap` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `linux_id` bigint NOT NULL,
  `total` bigint NOT NULL,
  `used` bigint NOT NULL,
  `free` bigint NOT NULL,
  `usedpercent` double NOT NULL,
  `sin` bigint NOT NULL,
  `sout` bigint NOT NULL,
  `pgin` bigint NOT NULL,
  `pgout` bigint NOT NULL,
  `pgfault` bigint NOT NULL,
  `pgmajfault` bigint NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_perf_swap_linux_id_timestamp` (`linux_id`,`timestamp`)
) ENGINE=InnoDB AUTO_INCREMENT=497090 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `permission`
--

DROP TABLE IF EXISTS `permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `identity` varchar(200) NOT NULL,
  `name` varchar(200) NOT NULL,
  `method` varchar(45) NOT NULL,
  `url` varchar(200) NOT NULL,
  `createTimestamp` bigint NOT NULL,
  `updateTimestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `identity` varchar(20) NOT NULL,
  `name` varchar(200) NOT NULL,
  `createTimestamp` bigint NOT NULL,
  `updateTimestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_role_identity` (`identity`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `role_permission`
--

DROP TABLE IF EXISTS `role_permission`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `role_id` bigint NOT NULL,
  `permission_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL,
  `passwd` varchar(200) NOT NULL,
  `is_active` int NOT NULL,
  `createTimestamp` bigint NOT NULL,
  `updateTimestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_user_username` (`username`) /*!80000 INVISIBLE */,
  KEY `idx_user_username_passwd_is_active` (`username`,`passwd`,`is_active`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_role`
--

DROP TABLE IF EXISTS `user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL,
  `role_id` bigint NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-02-25 13:46:29
