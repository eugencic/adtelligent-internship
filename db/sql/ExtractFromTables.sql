SELECT s.name AS source_name, COUNT(sc.source_id) AS campaign_count
FROM sources s
LEFT JOIN source_campaign sc ON s.id = sc.source_id
GROUP BY s.id
ORDER BY COUNT(sc.source_id) DESC
LIMIT 5;

SELECT c.name AS campaign_name
FROM campaigns c
LEFT JOIN source_campaign sc ON c.id = sc.campaign_id
WHERE sc.source_id IS NULL;

SELECT name FROM sources
UNION
SELECT name FROM campaigns;
