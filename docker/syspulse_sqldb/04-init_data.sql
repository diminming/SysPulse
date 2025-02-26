USE `syspulse`;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;
INSERT INTO `user` VALUES (0,'admin','admin1!',1,0,0);
/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UPDATE `user` SET id=0 WHERE `username`='admin';
UNLOCK TABLES;
