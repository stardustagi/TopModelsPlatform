-- 创建存储过程：按用户ID和CallerKey统计模型调用情况
DELIMITER //

DROP PROCEDURE IF EXISTS GetUserCallStatusReport //

CREATE PROCEDURE GetUserCallStatusReport(IN p_user_id BIGINT)
BEGIN
SELECT
    user_id,
    caller_key,
    model,

    -- 最近5分钟统计
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) THEN 1 ELSE 0 END) AS call_count_5min,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) AS success_count_5min,
    ROUND(
            SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) * 100.0 /
            NULLIF(SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) THEN 1 ELSE 0 END), 0),
            2
    ) AS success_rate_5min,
    ROUND(AVG(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 5 MINUTE) THEN CAST(latency AS DECIMAL(10,4)) ELSE NULL END), 4) AS avg_latency_5min,

    -- 最近10分钟统计
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 10 MINUTE) THEN 1 ELSE 0 END) AS call_count_10min,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 10 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) AS success_count_10min,
    ROUND(
            SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 10 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) * 100.0 /
            NULLIF(SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 10 MINUTE) THEN 1 ELSE 0 END), 0),
            2
    ) AS success_rate_10min,
    ROUND(AVG(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 10 MINUTE) THEN CAST(latency AS DECIMAL(10,4)) ELSE NULL END), 4) AS avg_latency_10min,

    -- 最近30分钟统计
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) THEN 1 ELSE 0 END) AS call_count_30min,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) AS success_count_30min,
    ROUND(
            SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) * 100.0 /
            NULLIF(SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) THEN 1 ELSE 0 END), 0),
            2
    ) AS success_rate_30min,
    ROUND(AVG(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 30 MINUTE) THEN CAST(latency AS DECIMAL(10,4)) ELSE NULL END), 4) AS avg_latency_30min,

    -- 最近1小时统计
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR) THEN 1 ELSE 0 END) AS call_count_1hour,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) AS success_count_1hour,
    ROUND(
            SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) * 100.0 /
            NULLIF(SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR) THEN 1 ELSE 0 END), 0),
            2
    ) AS success_rate_1hour,
    ROUND(AVG(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 1 HOUR) THEN CAST(latency AS DECIMAL(10,4)) ELSE NULL END), 4) AS avg_latency_1hour,

    -- 最近24小时统计
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 1 ELSE 0 END) AS call_count_24hour,
    SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) AS success_count_24hour,
    ROUND(
            SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) AND (status_code = '' OR status_code = '0') THEN 1 ELSE 0 END) * 100.0 /
            NULLIF(SUM(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN 1 ELSE 0 END), 0),
            2
    ) AS success_rate_24hour,
    ROUND(AVG(CASE WHEN created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR) THEN CAST(latency AS DECIMAL(10,4)) ELSE NULL END), 4) AS avg_latency_24hour

FROM status_report
WHERE step = 'llm_agent_done'
  AND created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
  AND (p_user_id IS NULL OR user_id = p_user_id)
GROUP BY user_id, caller_key, model
ORDER BY user_id, caller_key, model;
END //

DELIMITER ;
