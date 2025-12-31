-- 存储过程：GetUserDayReportList
-- 功能：分页获取用户日消费报表，关联用户、供应商、模型信息
-- 参数：
--   p_user_id: 用户ID（0表示查询所有用户）
--   p_skip: 跳过条数
--   p_limit: 每页条数
--   p_sort: 排序方式（默认 'id desc'）

DELIMITER //

DROP PROCEDURE IF EXISTS GetUserDayReportList //

CREATE PROCEDURE GetUserDayReportList(
    IN p_user_id BIGINT,
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

    -- 查询数据列表（第一个数据集）
    IF p_user_id IS NULL OR p_user_id = 0 THEN
        SELECT
            pds.id,
            pds.user_id,
            u.user_name,
            pds.actual_provider_id,
            mp.name AS provider_name,
            pds.model_id,
            mi.name AS model_name,
            pds.consume_type,
            pds.date,
            pds.total_consumed,
            pds.total_cost,
            pds.updated_at
        FROM provider_model_daily_summary pds
        LEFT JOIN users u ON pds.user_id = u.id
        LEFT JOIN models_provider mp ON pds.actual_provider_id = mp.id
        LEFT JOIN models_info mi ON pds.model_id = mi.id
        ORDER BY pds.id DESC
        LIMIT p_limit OFFSET p_skip;
    ELSE
        SELECT
            pds.id,
            pds.user_id,
            u.user_name,
            pds.actual_provider_id,
            mp.name AS provider_name,
            pds.model_id,
            mi.name AS model_name,
            pds.consume_type,
            pds.date,
            pds.total_consumed,
            pds.total_cost,
            pds.updated_at
        FROM provider_model_daily_summary pds
        LEFT JOIN users u ON pds.user_id = u.id
        LEFT JOIN models_provider mp ON pds.actual_provider_id = mp.id
        LEFT JOIN models_info mi ON pds.model_id = mi.id
        WHERE pds.user_id = p_user_id
        ORDER BY pds.id DESC
        LIMIT p_limit OFFSET p_skip;
    END IF;

    -- 查询总数（第二个数据集）
    IF p_user_id IS NULL OR p_user_id = 0 THEN
        SELECT COUNT(*) AS total
        FROM provider_model_daily_summary;
    ELSE
        SELECT COUNT(*) AS total
        FROM provider_model_daily_summary
        WHERE user_id = p_user_id;
    END IF;

END //

DELIMITER ;

