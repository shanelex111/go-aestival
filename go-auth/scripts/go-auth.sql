DROP TABLE
IF EXISTS account;
  CREATE TABLE `account` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '账户ID',
    `email` VARCHAR (255) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱',
    `phone_country_code` VARCHAR (8) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机号国家代码，如+86，+1',
    `phone_number` VARCHAR (20) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机号',
    `password` VARCHAR (255) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
    `status` TINYINT (1) NOT NULL DEFAULT '0' COMMENT '状态：1：enable，0：disable，-1：deleted',
    `signup_ip` VARCHAR (45) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '注册ip',
    `signup_ip_country_code` VARCHAR (8) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '注册ip国家码',
    `signup_ip_subdivision_code` VARCHAR (8) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '注册ip区域码',
    `signup_at` BIGINT NOT NULL COMMENT '注册时间，毫秒时间戳',
    `last_signin_ip` VARCHAR (45) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次登录IP',
    `last_signin_ip_country_code` VARCHAR (8) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次登录ip国家码',
    `last_signin_ip_subdivision_code` VARCHAR (8) CHARACTER
    SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次登录ip区域码',
    `last_signin_at` BIGINT NOT NULL COMMENT '最近一次登录时间',
    `signin_times` INT UNSIGNED NOT NULL COMMENT '总登录次数',
    `created_at` BIGINT NOT NULL COMMENT '创建时间，毫秒时间戳',
    `updated_at` BIGINT NOT NULL COMMENT '更新时间，毫秒时间戳',
    `deleted_at` BIGINT NOT NULL COMMENT '注销时间，毫秒时间戳',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_unique_email` (`email`),
    UNIQUE KEY `idx_unique_phone` (`phone_country_code`, `phone_number`)
  ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '账户表';
  DROP TABLE
  IF EXISTS device;
    CREATE TABLE device (
      id INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
      account_id INT UNSIGNED NOT NULL COMMENT '账户ID',
      device_id VARCHAR (255) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备id',
      device_type VARCHAR (10) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备类型，android/ipad/iphone/web/mweb',
      device_model VARCHAR (255) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备型号，如iPhone11，iPhone6，Letv X501\r\n\nWeb传递浏览器类型，如Chrome，Firefox',
      app_version INT UNSIGNED NOT NULL COMMENT '应用版本号，如：203000',
      signin_times INT UNSIGNED NOT NULL COMMENT '该设备总登录次数',
      created_ip VARCHAR (45) NOT NULL COMMENT '首次使用该设备登录时IP',
      created_ip_country_code VARCHAR (8) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备ip国家码',
      created_ip_subdivision_code VARCHAR (8) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备ip区域码',
      created_at BIGINT NOT NULL COMMENT '创建时间，毫秒时间戳',
      updated_ip VARCHAR (45) NOT NULL COMMENT '最近一次登录IP',
      updated_ip_country_code VARCHAR (8) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备ip国家码',
      updated_ip_subdivision_code VARCHAR (8) CHARACTER
      SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备ip区域码',
      updated_at BIGINT NOT NULL COMMENT '更新时间，毫秒时间戳',
      deleted_at BIGINT NOT NULL COMMENT '注销时间，毫秒时间戳',
      PRIMARY KEY (id),
      KEY idx_device_id (device_id),
      KEY idx_account_id (account_id)
    ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '设备表';
    DROP TABLE
    IF EXISTS verification_code;
      CREATE TABLE `verification_code` (
        `id` INT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键id',
        `account_id` INT UNSIGNED NOT NULL COMMENT '账户id',
        `type` VARCHAR (10) NOT NULL COMMENT '类型，email | phone',
        `target` VARCHAR (100) CHARACTER
        SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱/手机号',
        `country_code` VARCHAR (8) NOT NULL COMMENT '国家码',
        `code` VARCHAR (10) CHARACTER
        SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '验证码',
        `status` ENUM ('created', 'used', 'expired') CHARACTER
        SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '状态：created | used | expired',
        `expired_at` BIGINT UNSIGNED NOT NULL COMMENT '过期时间，毫秒时间戳',
        `created_at` BIGINT UNSIGNED NOT NULL COMMENT '创建时间，毫秒时间戳',
        `updated_at` BIGINT UNSIGNED NOT NULL COMMENT '更新时间，毫秒时间戳',
        `deleted_at` BIGINT UNSIGNED NOT NULL COMMENT '注销时间，毫秒时间戳',
        `platform` VARCHAR (20) CHARACTER
        SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '第三方：aliyun | tencent',
        `template_id` VARCHAR (100) CHARACTER
        SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模版id',
        `content` TEXT NOT NULL COMMENT '发送内容',
        PRIMARY KEY (`id`),
        KEY `idx_account_id` (`account_id`) USING BTREE,
        KEY `idx_target_type` (`type`, `target`) USING BTREE,
        KEY `idx_target_country` (`target`, `country_code`) USING BTREE
      ) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '验证码表';

DROP TABLE IF EXISTS `account_platform`;
CREATE TABLE `account_platform` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `account_id` int unsigned NOT NULL COMMENT '账户ID',
  `type` tinyint NOT NULL COMMENT '平台类型，0：未知，1：wechat，2：qq，3：weibo\r\n11：google，12：twitter，13：facebook',
  `platform_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '平台id',
  `platform_token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '平台access_token',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '头像',
  `created_at` bigint NOT NULL COMMENT '创建时间，毫秒时间戳',
  `updated_at` bigint NOT NULL COMMENT '更新时间，毫秒时间戳',
  `deleted_at` bigint NOT NULL COMMENT '注销时间，毫秒时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_account_id` (`account_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='第三方账户表';


DROP TABLE IF EXISTS `aestival_member`;
CREATE TABLE `aestival_member` (
  `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `account_id` int unsigned NOT NULL COMMENT '账户ID',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '昵称',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '头像相对地址',
  `role` tinyint(1) NOT NULL DEFAULT '0' COMMENT '角色：0：普通用户，1：vip',
  `created_at` bigint NOT NULL COMMENT '创建时间，毫秒时间戳',
  `updated_at` bigint NOT NULL COMMENT '更新时间，毫秒时间戳',
  `deleted_at` bigint NOT NULL COMMENT '注销时间，毫秒时间戳',
  PRIMARY KEY (`id`),
  KEY `idx_account_id` (`account_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户表';