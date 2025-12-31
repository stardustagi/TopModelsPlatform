-- 存储过程：GetUserInfoList
-- 功能：分页获取用户列表，包含用户信息和钱包余额
-- 参数：
--   p_skip: 跳过条数
--   p_limit: 每页条数
--   p_sort: 排序字段 (默认 'id')

DELIMITER //

DROP PROCEDURE IF EXISTS GetUserInfoList //

CREATE PROCEDURE GetUserInfoList(
    IN p_skip INT,
    IN p_limit INT,
    IN p_sort VARCHAR(50)
)
BEGIN
    -- 设置默认值
    IF p_skip IS NULL OR p_skip < 0 THEN
        SET p_skip = 0;
    END IF;

    IF p_limit IS NULL OR p_limit <= 0 THEN
        SET p_limit = 20;
    END IF;

    IF p_limit > 100 THEN
        SET p_limit = 100;
    END IF;

    IF p_sort IS NULL OR p_sort = '' THEN
        SET p_sort = 'id desc';
    END IF;

    -- 查询总数
    SELECT COUNT(*) AS total
    FROM users
    WHERE deleted = 0;

    -- 查询用户列表（不包含password和salt），并关联钱包信息
    SET @sql = CONCAT('
        SELECT
            u.id,
            u.user_name,
            u.email,
            u.phone,
            u.real_name,
            u.id_number,
            u.active,
            u.created_at,
            u.last_update,
            u.company_id,
            u.wallet_address_id,
            u.spread_id,
            u.is_realname_authentication,
            u.last_login_ip,
            u.is_ban,
            u.deleted,
            u.mail_code,
            u.phone_code,
            u.is_admin,
            u.is_private,
            IFNULL(w.balance, 0)/10000 AS balance,
            IFNULL(w.rebate_balance, 0)/10000 AS rebate_balance,
            w.wallet_type,
            w.wallet_address
        FROM users u
        LEFT JOIN user_wallet w ON u.id = w.user_id
        WHERE u.deleted = 0
        ORDER BY u.', p_sort, '
        LIMIT ', p_limit, ' OFFSET ', p_skip
    );

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

END //

DELIMITER ;

