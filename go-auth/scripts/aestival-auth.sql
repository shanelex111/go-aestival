DROP TABLE IF EXISTS account;

CREATE TABLE `account` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '账户ID',
    `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱',
    `phone_country_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机号国家代码，如+86，+1',
    `phone_number` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '手机号',
    `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
    `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：1：enable，0：disable，-1：deleted',
    `created_at` bigint NOT NULL COMMENT '创建时间，毫秒时间戳',
    `updated_at` bigint NOT NULL COMMENT '更新时间，毫秒时间戳',
    `deleted_at` bigint NOT NULL COMMENT '注销时间，毫秒时间戳',
    PRIMARY KEY (`id`),
    KEY `idx_unique_email` (`email`) USING BTREE,
    KEY `idx_unique_phone` (
        `phone_country_code`,
        `phone_number`
    ) USING BTREE,
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '账户表';

DROP TABLE IF EXISTS device;

CREATE TABLE `device` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `account_id` int unsigned NOT NULL COMMENT '账户ID',
    `device_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备id',
    `device_type` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备类型，android/ipad/iphone/web/mweb',
    `device_model` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '设备型号，如iPhone11，iPhone6，Letv X501\r\n\nWeb传递浏览器类型，如Chrome，Firefox',
    `app_version` int unsigned NOT NULL COMMENT '应用版本号，如：203000',
    `signin_times` int unsigned NOT NULL COMMENT '该设备总登录次数',
    `created_ip` varchar(45) NOT NULL COMMENT '首次使用该设备登录时IP',
    `created_ip_content_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备登录时IP大洲代码',
    `created_ip_country_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备IP国家码',
    `created_ip_subdivision_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备IP区域码',
    `created_ip_city_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '首次使用该设备登录时IP城市名称',
    `created_at` bigint NOT NULL COMMENT '创建时间，毫秒时间戳',
    `updated_ip` varchar(45) NOT NULL COMMENT '最近一次登录IP',
    `updated_ip_continent_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备登录时IP大洲代码',
    `updated_ip_country_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备IP国家码',
    `updated_ip_subdivision_code` varchar(8) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备IP区域码',
    `updated_ip_city_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '最近一次该设备IP城市名称',
    `updated_at` bigint NOT NULL COMMENT '更新时间，毫秒时间戳',
    `deleted_at` bigint NOT NULL COMMENT '注销时间，毫秒时间戳',
    PRIMARY KEY (`id`),
    KEY `idx_device_id` (`device_id`),
    KEY `idx_account_id` (`account_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '设备表';

DROP TABLE IF EXISTS verification_code;

CREATE TABLE `verification_code` (
    `id` int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `scene` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '场景: signin | reset_password',
    `type` varchar(10) NOT NULL COMMENT '类型，email | phone',
    `target` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱/手机号',
    `country_code` varchar(8) NOT NULL COMMENT '国家码',
    `code` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '验证码',
    `status` enum('pending', 'used', 'expired') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '状态：pending | used | expired',
    `expired_at` bigint unsigned NOT NULL COMMENT '过期时间，毫秒时间戳',
    `created_at` bigint unsigned NOT NULL COMMENT '创建时间，毫秒时间戳',
    `updated_at` bigint unsigned NOT NULL COMMENT '更新时间，毫秒时间戳',
    `deleted_at` bigint unsigned NOT NULL COMMENT '注销时间，毫秒时间戳',
    `platform` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '第三方：aliyun | tencent',
    `template_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '模版id',
    `content` text NOT NULL COMMENT '发送内容',
    PRIMARY KEY (`id`),
    KEY `idx_target_type` (`type`, `target`) USING BTREE,
    KEY `idx_target_country` (`target`, `country_code`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '验证码表';

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '第三方账户表';

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
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '用户表';