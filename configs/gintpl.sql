-- Copyright 2025 莫维龙 <kalandramo@gmail.com>. All rights reserved.
-- Use of this source code is governed by a MIT style
-- license that can be found in the LICENSE file. The original repo for
-- this file is https://github.com/srxstack/gintpl. The professional
-- version of this repository is https://github.com/srxstack/srxstack.

CREATE USER 'gintpl'@'%' IDENTIFIED BY 'gintpl(#)888';
GRANT ALL ON gintpl.* TO 'gintpl'@'%';
FLUSH PRIVILEGES;

DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `ptype` varchar(100) DEFAULT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=1;

LOCK TABLES `casbin_rule` WRITE;
INSERT INTO `casbin_rule` VALUES
(1,'g','user-000000','role::admin',NULL,NULL,'',''),
(2,'p','role::admin','*','*','allow','',''),
(3,'p','role::user','/v1.MiniBlog/DeleteUser','CALL','deny','',''),
(4,'p','role::user','/v1.MiniBlog/ListUser','CALL','deny','',''),
(5,'p','role::user','/v1/users','GET','deny','',''),
(6,'p','role::user','/v1/users/*','DELETE','deny','','');
UNLOCK TABLES;

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `userID` varchar(36) NOT NULL DEFAULT '' COMMENT '用户唯一 ID',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名（唯一）',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '用户密码（加密后）',
  `nickname` varchar(30) NOT NULL DEFAULT '' COMMENT '用户昵称',
  `email` varchar(256) NOT NULL DEFAULT '' COMMENT '用户电子邮箱地址',
  `phone` varchar(16) NOT NULL DEFAULT '' COMMENT '用户手机号',
  `createdAt` datetime NOT NULL DEFAULT current_timestamp() COMMENT '用户创建时间',
  `updatedAt` datetime NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp() COMMENT '用户最后修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user.userID` (`userID`),
  UNIQUE KEY `user.username` (`username`),
  UNIQUE KEY `user.phone` (`phone`)
) ENGINE=InnoDB AUTO_INCREMENT=1 COMMENT='用户表';


LOCK TABLES `user` WRITE;
INSERT INTO `user` VALUES
(96,'user-000000','root','$2a$10$ctsFXEUAMd7rXXpmccNlO.ZRiYGYz0eOfj8EicPGWqiz64YBBgR1y','moweilong','kalandramo@gmail.com','18110000000','2025-05-12 03:04:05','2025-05-12 03:04:05');
UNLOCK TABLES;