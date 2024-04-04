DROP TABLE IF EXISTS source_campaign;
DROP TABLE IF EXISTS sources;
DROP TABLE IF EXISTS campaigns;

CREATE TABLE IF NOT EXISTS sources (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        name VARCHAR(50)
                                   );

CREATE TABLE IF NOT EXISTS campaigns (
                                        id INT AUTO_INCREMENT PRIMARY KEY,
                                        name VARCHAR(50)
                                     );

CREATE TABLE IF NOT EXISTS source_campaign (
                                                source_id INT,
                                                campaign_id INT,
                                                FOREIGN KEY (source_id) REFERENCES sources(id),
                                                FOREIGN KEY (campaign_id) REFERENCES campaigns(id)
                                           );
