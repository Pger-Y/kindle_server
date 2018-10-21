create table `user_info`(`userid` char(28) primary key,
	`kindle_address` varchar(256),
	`update_time` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`mail_address` varchar(256),
	`mail_passwd` varchar(32),
	`smtp_server` varchar(24)
)engine=Innodb default charset='utf-8'
