create database beyond_user;
use beyond_user;

CREATE TABLE `user` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
  `avatar` varchar(256) NOT NULL DEFAULT '' COMMENT '头像',
  `mobile` varchar(128) NOT NULL DEFAULT '' COMMENT '手机号',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
  PRIMARY KEY (`id`),
  KEY `ix_mtime` (`mtime`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='用户表';

